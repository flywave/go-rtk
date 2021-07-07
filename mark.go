package rtk

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
)

type PhaseComp struct {
	v float64
	d string
}

func (date *PhaseComp) MarshalCSV() (string, error) {
	return fmt.Sprintf("%f,%s", date.v, date.d), nil
}

func (date *PhaseComp) String() string {
	return fmt.Sprintf("%f,%s", date.v, date.d)
}

func (date *PhaseComp) UnmarshalCSV(csv string) (err error) {
	strs := strings.Split(csv, ",")
	strs[0] = strings.Trim(strs[0], " ")
	date.v, err = strconv.ParseFloat(strs[0], 64)
	date.d = strings.Trim(strs[1], " ")
	if err != nil {
		return
	}
	return
}

type Week struct {
	w int64
}

func (w *Week) MarshalCSV() (string, error) {
	return fmt.Sprintf("[%d]", w.w), nil
}

func (w *Week) String() string {
	return fmt.Sprintf("[%d]", w.w)
}

func (w *Week) UnmarshalCSV(csv string) (err error) {
	csv = strings.Trim(csv, " ")
	if strings.HasPrefix(csv, "[") && strings.HasSuffix(csv, "]") {
		s := csv[1 : len(csv)-1]
		w.w, err = strconv.ParseInt(s, 10, 32)
	}
	return
}

type Angle struct {
	v float64
	d string
}

func (date *Angle) MarshalCSV() (string, error) {
	return fmt.Sprintf("%f,%s", date.v, date.d), nil
}

func (date *Angle) String() string {
	return fmt.Sprintf("%f,%s", date.v, date.d)
}

func (date *Angle) UnmarshalCSV(csv string) (err error) {
	strs := strings.Split(csv, ",")
	strs[0] = strings.Trim(strs[0], " ")
	date.v, err = strconv.ParseFloat(strs[0], 64)
	date.d = strings.Trim(strs[1], " ")
	if err != nil {
		return
	}
	return
}

type RtkState struct {
	e int64
	f string
}

func (date *RtkState) MarshalCSV() (string, error) {
	return fmt.Sprintf("%d,%s", date.e, date.f), nil
}

func (date *RtkState) String() string {
	return fmt.Sprintf("%d,%s", date.e, date.f)
}

func (date *RtkState) UnmarshalCSV(csv string) (err error) {
	strs := strings.Split(csv, ",")
	date.e, err = strconv.ParseInt(strs[0], 10, 32)
	date.f = strs[1]
	if err != nil {
		return
	}
	return
}

type Std struct {
	latitude  float64
	longitude float64
	altitude  float64
}

func (date *Std) MarshalCSV() (string, error) {
	return fmt.Sprintf("%f, %f, %f", date.latitude, date.longitude, date.altitude), nil
}

func (date *Std) String() string {
	return fmt.Sprintf("%f, %f, %f", date.latitude, date.longitude, date.altitude)
}

func (date *Std) UnmarshalCSV(csv string) (err error) {
	strs := strings.Split(csv, ",")

	strs[0] = strings.Trim(strs[0], " ")
	strs[1] = strings.Trim(strs[1], " ")
	strs[2] = strings.Trim(strs[2], " ")

	date.latitude, err = strconv.ParseFloat(strs[0], 64)
	if err != nil {
		return
	}

	date.longitude, err = strconv.ParseFloat(strs[1], 64)
	if err != nil {
		return
	}

	date.altitude, err = strconv.ParseFloat(strs[2], 64)
	if err != nil {
		return
	}
	return
}

type MRK struct {
	Sequence            int32     `csv:"Sequence"`
	Time                float64   `csv:"GPSSecondOfWeek"`
	Week                Week      `csv:"GPSWeek"`
	PhaseCompNs         PhaseComp `csv:"NorthOff"`
	PhaseCompEw         PhaseComp `csv:"EastOff"`
	PhaseCompV          PhaseComp `csv:"VelOff"`
	Latitude            Angle     `csv:"Latitude"`
	Longitude           Angle     `csv:"Longitude"`
	Altitude            Angle     `csv:"EllipsoideHight"`
	Std                 Std       `csv:"Std"`
	State               RtkState  `csv:"RtkState"`
	closestId           int32
	closestTime         GTime
	closestLatitude     float64
	closestLongitude    float64
	closestAltitude     float64
	sndClosestId        int32
	sndClosestTime      GTime
	sndClosestLatitude  float64
	sndClosestLongitude float64
	sndClosestAltitude  float64
	gpstDiff            float64
	interpLatitude      float64
	interpLongitude     float64
	interpAltitude      float64
	phaseCompNsDeg      float64
	phaseCompEwDeg      float64
	phaseCompVM         float64
	updatedLatitude     float64
	updatedLongitude    float64
	updatedAltitude     float64
	diffLatitude        float64
	diffLongitude       float64
	diffAltitude        float64
}

func ReadMRK(reader io.Reader) ([]MRK, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '\t'
	var mrks []MRK
	if err := gocsv.UnmarshalCSVWithoutHeaders(csvReader, &mrks); err != nil {
		return nil, err
	}
	return mrks, nil
}
