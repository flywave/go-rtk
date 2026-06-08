package rtk

import (
	"testing"
)

func TestGTimeFromGPSTimeRoundTrip(t *testing.T) {
	week := int32(2024)
	sec := 275295.301059
	gt := NewGTimeFromGPSTime(week, sec)
	if gt == nil {
		t.FailNow()
	}
	w, _ := gt.GpsTime()
	if w != week {
		t.FailNow()
	}
}

func TestGTimeFromEpoch(t *testing.T) {
	ep := [6]float64{2024, 6, 15, 10, 30, 0}
	gt := NewGTimeFromEpoch(ep)
	if gt == nil {
		t.FailNow()
	}
}

func TestGTimeDiff(t *testing.T) {
	gt1 := NewGTimeFromGPSTime(2024, 100.0)
	gt2 := NewGTimeFromGPSTime(2024, 200.0)
	if gt1 == nil || gt2 == nil {
		t.FailNow()
	}
	diff := gt2.Diff(gt1)
	if diff != 100.0 {
		t.FailNow()
	}
}

func TestGTimeAdd(t *testing.T) {
	gt := NewGTimeFromGPSTime(2024, 100.0)
	if gt == nil {
		t.FailNow()
	}
	gt.Add(50.0)
	_, s := gt.GpsTime()
	if s != 150.0 {
		t.FailNow()
	}
}

func TestGTimeEpoch(t *testing.T) {
	gt := NewGTimeFromGPSTime(2024, 100.0)
	if gt == nil {
		t.FailNow()
	}
	ep := gt.Epoch()
	if ep[0] < 2000 || ep[0] > 2100 {
		t.FailNow()
	}
}

func TestGTimeToString(t *testing.T) {
	gt := NewGTimeFromGPSTime(2024, 275295.0)
	if gt == nil {
		t.FailNow()
	}
	s := gt.ToString(0)
	if len(s) == 0 {
		t.FailNow()
	}
}

func TestGTimeDayOfYear(t *testing.T) {
	gt := NewGTimeFromGPSTime(2024, 275295.0)
	if gt == nil {
		t.FailNow()
	}
	doy := gt.DayOfYear(0)
	if doy < 1 || doy > 366 {
		t.FailNow()
	}
}

func TestNewGPSTFromUTC(t *testing.T) {
	utc := NewGTimeFromGPSTime(2024, 275295.0)
	if utc == nil {
		t.FailNow()
	}
	gpst := NewGPSTFromUTC(utc)
	if gpst == nil {
		t.FailNow()
	}
}

func TestGTimeUTCConversion(t *testing.T) {
	gpst := NewGTimeFromGPSTime(2024, 275295.0)
	if gpst == nil {
		t.FailNow()
	}
	utc := gpst.UTC()
	if utc == nil {
		t.FailNow()
	}
}

func TestGalileoTime(t *testing.T) {
	gt := NewGTimeFromGalileoTime(2024, 100.0)
	if gt == nil {
		t.FailNow()
	}
	w, s := gt.GalileoTime()
	if w != 2024 {
		t.FailNow()
	}
	if s != 100.0 {
		t.FailNow()
	}
}

func TestBDTime(t *testing.T) {
	gt := NewGTimeFromBDTime(2024, 100.0)
	if gt == nil {
		t.FailNow()
	}
	w, s := gt.BDTime()
	if w != 2024 {
		t.FailNow()
	}
	if s != 100.0 {
		t.FailNow()
	}
}

func TestCurrent(t *testing.T) {
	ct := Current()
	if ct == nil {
		t.FailNow()
	}
}
