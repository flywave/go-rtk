package rtk

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
)

type YearMonthDay struct {
	year  int64
	month int64
	day   int64
}

func (date *YearMonthDay) MarshalCSV() (string, error) {
	return fmt.Sprintf("%d/%d/%d", date.year, date.month, date.day), nil
}

func (date *YearMonthDay) String() string {
	return fmt.Sprintf("%d/%d/%d", date.year, date.month, date.day)
}

func (date *YearMonthDay) UnmarshalCSV(csv string) (err error) {
	strs := strings.Split(csv, "/")

	strs[0] = strings.Trim(strs[0], " ")
	strs[1] = strings.Trim(strs[1], " ")
	strs[2] = strings.Trim(strs[2], " ")

	date.year, err = strconv.ParseInt(strs[0], 10, 32)
	if err != nil {
		return
	}
	date.month, err = strconv.ParseInt(strs[1], 10, 32)
	if err != nil {
		return
	}
	date.day, err = strconv.ParseInt(strs[2], 10, 32)
	if err != nil {
		return
	}
	return
}

type Time struct {
	hour   int64
	minute int64
	second float64
}

func (date *Time) MarshalCSV() (string, error) {
	return fmt.Sprintf("%d:%d:%f", date.hour, date.minute, date.second), nil
}

func (date *Time) String() string {
	return fmt.Sprintf("%d:%d:%f", date.hour, date.minute, date.second)
}

func (date *Time) UnmarshalCSV(csv string) (err error) {
	strs := strings.Split(csv, ":")

	strs[0] = strings.Trim(strs[0], " ")
	strs[2] = strings.Trim(strs[2], " ")

	date.hour, err = strconv.ParseInt(strs[0], 10, 32)
	if err != nil {
		return
	}
	date.minute, err = strconv.ParseInt(strs[1], 10, 32)
	if err != nil {
		return
	}
	date.second, err = strconv.ParseFloat(strs[2], 64)
	if err != nil {
		return
	}
	return
}

type Pos struct {
	Day             YearMonthDay `csv:"GPSTDay"`
	Time            Time         `csv:"GPSTTime"`
	LatitudeDegree  float64      `csv:"LatitudeDegree"`
	LatitudeMinute  float64      `csv:"LatitudeMinute"`
	LatitudeSecond  float64      `csv:"LatitudeSecond"`
	LongitudeDegree float64      `csv:"LongitudeDegree"`
	LongitudeMinute float64      `csv:"LongitudeMinute"`
	LongitudeSecond float64      `csv:"LongitudeSecond"`
	Height          float64      `csv:"Height"`
	Q               int8         `csv:"Q"`
	Ns              int16        `csv:"ns"`
	Sdn             float64      `csv:"sdn(m)"`
	Sde             float64      `csv:"sde(m)"`
	Sdu             float64      `csv:"sdu(m)"`
	Sdne            float64      `csv:"sdne(m)"`
	Sdeu            float64      `csv:"sdeu(m)"`
	Sdun            float64      `csv:"sdue(m)"`
	Age             float64      `csv:"age(s)"`
	Ratio           float64      `csv:"ratio"`
}

func ReadPos(reader io.Reader) ([]Pos, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comment = '%'
	csvReader.Comma = ' '

	var poss []Pos
	if err := gocsv.UnmarshalCSVWithoutHeaders(csvReader, &poss); err != nil {
		return nil, err
	}
	return poss, nil
}
