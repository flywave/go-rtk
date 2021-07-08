package rtk

import (
	"math"
	"os"

	"github.com/flywave/go-geoid"
	"github.com/flywave/go-proj"
)

type PPKOption struct {
	BaseAerialHeight  float32
	RoverAerialHeight float32
	RoverCameraOffset [3]float32
}

func DegreeToRadian(degree float64) float64 {
	return (degree * math.Pi) / 180.0
}
func RadianToDegree(radian float64) float64 {
	return (radian * 180.0) / math.Pi
}

func nearestTimeInd(poses []Pos, pivot GTime) int {
	minindex := 0
	mindiff := math.MaxFloat64
	for i := range poses {
		if math.Abs(poses[i].Gpst.Diff(&pivot)) < mindiff {
			mindiff = math.Abs(poses[i].Gpst.Diff(&pivot))
			minindex = i
		}
	}
	return minindex
}

func ppkInitPos(poses []Pos, mrks []MRK) {
	for i := range mrks {
		first := nearestTimeInd(poses, mrks[i].GetGpst())
		mrks[i].closestId = int32(first)
		mrks[i].closestTime = poses[i].Gpst
		mrks[i].closestLatitude = poses[i].Latitude
		mrks[i].closestLongitude = poses[i].Longitude
		mrks[i].closestAltitude = poses[i].Height

		pivot := mrks[i].GetGpst()

		dta := poses[first-1].Gpst.Diff(&pivot)
		dtb := poses[first+1].Gpst.Diff(&pivot)
		var second int
		if math.Abs(dta) < math.Abs(dtb) {
			second = first - 1
		} else {
			second = first + 1
		}
		mrks[i].sndClosestId = int32(second)
		mrks[i].sndClosestTime = poses[second].Gpst
		mrks[i].sndClosestLatitude = poses[second].Latitude
		mrks[i].sndClosestLongitude = poses[second].Longitude
		mrks[i].sndClosestAltitude = poses[second].Height
	}
}

func ppkUpdatedMrks(mrks []MRK) {
	degLon := math.Cos(RadianToDegree(mrks[0].Latitude.v)) * 111.321

	for i := range mrks {
		mtime := mrks[i].GetGpst()
		mrks[i].gpstDiff = mtime.Diff(&mrks[i].closestTime) /
			mrks[i].sndClosestTime.Diff(&mrks[i].closestTime)

		mrks[i].interpLatitude =
			(mrks[i].closestLatitude * (1 - mrks[i].gpstDiff)) +
				(mrks[i].sndClosestLatitude * mrks[i].gpstDiff)
		mrks[i].interpLongitude =
			(mrks[i].closestLongitude * (1 - mrks[i].gpstDiff)) +
				(mrks[i].sndClosestLongitude * mrks[i].gpstDiff)
		mrks[i].interpAltitude = (mrks[i].closestAltitude * (1 - mrks[i].gpstDiff)) +
			(mrks[i].sndClosestAltitude * mrks[i].gpstDiff)

		mrks[i].phaseCompNsDeg = mrks[i].PhaseCompNs.v / 1000000 / 111.111
		mrks[i].phaseCompEwDeg = mrks[i].PhaseCompEw.v / 1000000 / degLon
		mrks[i].phaseCompVM = mrks[i].PhaseCompV.v / 1000

		mrks[i].updatedLatitude =
			mrks[i].interpLatitude + mrks[i].phaseCompNsDeg
		mrks[i].updatedLongitude =
			mrks[i].interpLongitude + mrks[i].phaseCompEwDeg
		mrks[i].updatedAltitude = mrks[i].interpAltitude - mrks[i].phaseCompVM

		mrks[i].diffLatitude = mrks[i].Latitude.v - mrks[i].updatedLatitude
		mrks[i].diffLongitude = mrks[i].Longitude.v - mrks[i].updatedLongitude
		mrks[i].diffAltitude = mrks[i].Altitude.v - mrks[i].updatedAltitude
	}
}

func convertGPS(mrk *MRK, pose_datum Datum, pose_vertica_datum geoid.VerticalDatum,
	srs string, vertica_datum geoid.VerticalDatum, ellipsoid_offset float64) [3]float64 {
	srspj, _ := proj.NewProj(srs)
	var val [3]float64
	latitude := mrk.updatedLatitude
	longitude := mrk.updatedLongitude
	altitude := mrk.updatedAltitude
	if pose_vertica_datum != geoid.HAE &&
		pose_vertica_datum != geoid.UNKNOWN {
		altitude = MSLToWGS84(altitude, longitude, latitude,
			pose_vertica_datum)
	}
	if isDatums(srspj, pose_datum) {
		val = transformFromLLA(
			srspj, [3]float64{longitude, latitude, altitude}, pose_datum)
	}
	if vertica_datum != geoid.HAE &&
		vertica_datum != geoid.UNKNOWN {
		val[2] = WGS84ToMSL(longitude, latitude, altitude, vertica_datum)
	}
	return val
}

type PPKSol struct {
	Pos    [3]float64
	Weight [3]float64
}

func PPKSolution(posfile string, markfile string, pose_datum Datum, pose_vertica_datum geoid.VerticalDatum,
	srs string, vertica_datum geoid.VerticalDatum, ellipsoid_offset float64) (error, []PPKSol) {
	poses, _ := ReadPos(posfile)
	f, err := os.Open(markfile)

	if err != nil {
		return err, nil
	}
	defer f.Close()
	mrks, err := ReadMRK(f)

	if err != nil {
		return err, nil
	}

	ppkInitPos(poses, mrks)
	ppkUpdatedMrks(mrks)

	sols := make([]PPKSol, len(mrks))
	for i := range sols {
		sols[i].Pos = convertGPS(
			&mrks[i], pose_datum, vertica_datum,
			srs,
			vertica_datum, ellipsoid_offset)

		sols[i].Weight = [3]float64{1 / mrks[i].Std.longitude, 1 / mrks[i].Std.latitude,
			1 / mrks[i].Std.altitude}
	}
	return nil, sols
}
