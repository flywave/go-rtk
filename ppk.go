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
	if len(poses) == 0 || len(mrks) == 0 {
		return
	}
	for i := range mrks {
		first := nearestTimeInd(poses, mrks[i].GetGpst())
		mrks[i].closestId = int32(first)
		mrks[i].closestTime = poses[first].Gpst
		mrks[i].closestLatitude = poses[first].Latitude
		mrks[i].closestLongitude = poses[first].Longitude
		mrks[i].closestAltitude = poses[first].Height

		pivot := mrks[i].GetGpst()

		var second int
		if len(poses) == 1 {
			second = 0
		} else if first == 0 {
			second = 1
		} else if first == len(poses)-1 {
			second = len(poses) - 2
		} else {
			dta := poses[first-1].Gpst.Diff(&pivot)
			dtb := poses[first+1].Gpst.Diff(&pivot)
			if math.Abs(dta) < math.Abs(dtb) {
				second = first - 1
			} else {
				second = first + 1
			}
		}
		mrks[i].sndClosestId = int32(second)
		mrks[i].sndClosestTime = poses[second].Gpst
		mrks[i].sndClosestLatitude = poses[second].Latitude
		mrks[i].sndClosestLongitude = poses[second].Longitude
		mrks[i].sndClosestAltitude = poses[second].Height
	}
}

func ppkUpdatedMrks(mrks []MRK) {
	if len(mrks) == 0 {
		return
	}

	for i := range mrks {
		mtime := mrks[i].GetGpst()
		denom := mrks[i].sndClosestTime.Diff(&mrks[i].closestTime)
		if denom == 0 {
			mrks[i].gpstDiff = 0
		} else {
			mrks[i].gpstDiff = mtime.Diff(&mrks[i].closestTime) / denom
		}

		mrks[i].interpLatitude =
			(mrks[i].closestLatitude * (1 - mrks[i].gpstDiff)) +
				(mrks[i].sndClosestLatitude * mrks[i].gpstDiff)
		mrks[i].interpLongitude =
			(mrks[i].closestLongitude * (1 - mrks[i].gpstDiff)) +
				(mrks[i].sndClosestLongitude * mrks[i].gpstDiff)
		mrks[i].interpAltitude = (mrks[i].closestAltitude * (1 - mrks[i].gpstDiff)) +
			(mrks[i].sndClosestAltitude * mrks[i].gpstDiff)

		degLon := math.Cos(DegreeToRadian(mrks[i].Latitude.v)) * 111.321
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
	srspj, err := proj.NewProj(srs)
	if err != nil {
		return [3]float64{}
	}
	var val [3]float64
	latitude := mrk.updatedLatitude
	longitude := mrk.updatedLongitude
	altitude := mrk.updatedAltitude
	if pose_vertica_datum != geoid.HAE &&
		pose_vertica_datum != geoid.UNKNOWN {
		altitude = MSLToWGS84(altitude, longitude, latitude,
			pose_vertica_datum)
	}
	altitude += ellipsoid_offset
	if isDatums(srspj, pose_datum) {
		val = transformFromLLA(
			srspj, [3]float64{longitude, latitude, altitude}, pose_datum)
	} else {
		val = [3]float64{longitude, latitude, altitude}
	}
	if vertica_datum != geoid.HAE &&
		vertica_datum != geoid.UNKNOWN {
		val[2] = WGS84ToMSL(val[0], val[1], val[2], vertica_datum)
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
	if len(poses) == 0 {
		return nil, nil
	}
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
			&mrks[i], pose_datum, pose_vertica_datum,
			srs,
			vertica_datum, ellipsoid_offset)

		wlon, wlat, walt := 1.0, 1.0, 1.0
		if mrks[i].Std.longitude != 0 {
			wlon = 1 / mrks[i].Std.longitude
		}
		if mrks[i].Std.latitude != 0 {
			wlat = 1 / mrks[i].Std.latitude
		}
		if mrks[i].Std.altitude != 0 {
			walt = 1 / mrks[i].Std.altitude
		}
		sols[i].Weight = [3]float64{wlon, wlat, walt}
	}
	return nil, sols
}
