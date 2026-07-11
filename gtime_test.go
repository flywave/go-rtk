package rtk

import (
	"testing"
	"time"
)

func TestGTimeFromGPSTimeRoundTrip(t *testing.T) {
	week := int32(2024)
	sec := 275295.301059
	gt := NewGTimeFromGPSTime(week, sec)
	if gt == nil {
		t.FailNow()
	}
	w, s := gt.GpsTime()
	if w != week {
		t.FailNow()
	}
	if s != sec {
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

func TestGTimeDiffNegative(t *testing.T) {
	gt1 := NewGTimeFromGPSTime(2024, 200.0)
	gt2 := NewGTimeFromGPSTime(2024, 100.0)
	diff := gt2.Diff(gt1)
	if diff != -100.0 {
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

func TestGTimeAddNegative(t *testing.T) {
	gt := NewGTimeFromGPSTime(2024, 100.0)
	gt.Add(-30.0)
	_, s := gt.GpsTime()
	if s != 70.0 {
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
	doy := gt.DayOfYear()
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

func TestNewGPSTFromUTCNil(t *testing.T) {
	gpst := NewGPSTFromUTC(nil)
	if gpst != nil {
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

func TestGTimeTimeAndSec(t *testing.T) {
	gt := NewGTimeFromGPSTime(2024, 275295.301059)
	if gt == nil {
		t.FailNow()
	}
	if gt.Sec() == 0 && gt.Time() == 0 {
		t.FailNow()
	}
}

func TestNewGTimeFromStr(t *testing.T) {
	gt := NewGTimeFromStr("2024/06/15 10:30:45")
	if gt == nil {
		t.FailNow()
	}
	ep := gt.Epoch()
	if ep[0] < 2024 || ep[0] > 2025 {
		t.FailNow()
	}
}

func TestNewGTimeFromStrWithFraction(t *testing.T) {
	gt := NewGTimeFromStr("2024/06/15 10:30:45.500")
	if gt == nil {
		t.FailNow()
	}
}

func TestNewUtcTime(t *testing.T) {
	gt := NewUtcTime("2024/06/15 10:30:45")
	if gt == nil {
		t.FailNow()
	}
}

func TestNewUtcTimeWithFraction(t *testing.T) {
	gt := NewUtcTime("2024/06/15 10:30:45.500")
	if gt == nil {
		t.FailNow()
	}
}

func TestNewGPSTime(t *testing.T) {
	gt := NewGPSTTime("2024/06/15 10:30:45")
	if gt == nil {
		t.FailNow()
	}
}

func TestNewUtcTimeFromTime(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 30, 45, 0, time.UTC)
	gt := NewUtcTimeFromTime(now)
	if gt == nil {
		t.FailNow()
	}
	ep := gt.Epoch()
	if ep[0] != 2024 || ep[1] != 6 || ep[2] != 15 {
		t.FailNow()
	}
}

func TestGPSTimeFromTime(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 30, 45, 0, time.UTC)
	gt := NewGPSTTimeFromTime(now)
	if gt == nil {
		t.FailNow()
	}
}

func TestNewUtcTimeFromLocalInvalid(t *testing.T) {
	gt := NewUtcTimeFromLocal("invalid", "UTC")
	if gt != nil {
		t.FailNow()
	}
}

func TestNewGPSTTimeFromLocalInvalid(t *testing.T) {
	gt := NewGPSTTimeFromLocal("invalid", "UTC")
	if gt != nil {
		t.FailNow()
	}
}

func TestNewUtcTimeFromCurrentLocalInvalid(t *testing.T) {
	gt := NewUtcTimeFromCurrentLocal("invalid")
	if gt != nil {
		t.FailNow()
	}
}

func TestEpochRoundTrip(t *testing.T) {
	ep := [6]float64{2024, 6, 15, 10, 30, 45.5}
	gt := NewGTimeFromEpoch(ep)
	if gt == nil {
		t.FailNow()
	}
	ep2 := gt.Epoch()
	if ep[0] != ep2[0] || ep[1] != ep2[1] || ep[2] != ep2[2] {
		t.FailNow()
	}
}
