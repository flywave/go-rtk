package rtk

import (
	"testing"
	"time"
)

func TestParseUtcTimeStandard(t *testing.T) {
	ep, sec := ParseUtcTime("2024/06/15 10:30:45.123")
	if ep[0] != 2024 {
		t.FailNow()
	}
	if ep[1] != 6 {
		t.FailNow()
	}
	if ep[2] != 15 {
		t.FailNow()
	}
	if ep[3] != 10 {
		t.FailNow()
	}
	if ep[4] != 30 {
		t.FailNow()
	}
	if ep[5] != 45 {
		t.FailNow()
	}
	if sec != 0.123 {
		t.FailNow()
	}
}

func TestParseUtcTimeNoFraction(t *testing.T) {
	ep, sec := ParseUtcTime("2024/06/15 10:30:45")
	if ep[0] != 2024 {
		t.FailNow()
	}
	if sec != 0.0 {
		t.FailNow()
	}
}

func TestParseUtcTimeLowerCase(t *testing.T) {
	ep, _ := ParseUtcTime("2024/06/15 10:30:45.500")
	if ep[0] != 2024 {
		t.FailNow()
	}
}

func TestParseUtcTimeLeadingZeros(t *testing.T) {
	ep, sec := ParseUtcTime("2024/01/02 03:04:05.006")
	if ep[1] != 1 || ep[2] != 2 || ep[3] != 3 || ep[4] != 4 || ep[5] != 5 {
		t.FailNow()
	}
	if sec != 0.006 {
		t.FailNow()
	}
}

func TestNewTm(t *testing.T) {
	tm, err := parseWithLocation("UTC", "2024-06-15 10:30:45")
	if err != nil {
		t.FailNow()
	}
	_ = NewTm(tm)
}

func TestStrftime(t *testing.T) {
	tm, err := parseWithLocation("UTC", "2024-06-15 10:30:45")
	if err != nil {
		t.FailNow()
	}
	gotm := NewTm(tm)
	s := Strftime("%Y-%m-%d", gotm)
	if s != "2024-06-15" {
		t.FailNow()
	}
}

func TestStrftimeFullFormat(t *testing.T) {
	tm, err := parseWithLocation("UTC", "2024-06-15 10:30:45")
	if err != nil {
		t.FailNow()
	}
	gotm := NewTm(tm)
	s := Strftime("%Y-%m-%d %H:%M:%S", gotm)
	if s != "2024-06-15 10:30:45" {
		t.FailNow()
	}
}

func TestStrftimeEmptyFormat(t *testing.T) {
	tm, err := parseWithLocation("UTC", "2024-06-15 10:30:45")
	if err != nil {
		t.FailNow()
	}
	gotm := NewTm(tm)
	s := Strftime("", gotm)
	if s != "" {
		t.FailNow()
	}
}

func TestStrftimeLongFormat(t *testing.T) {
	tm, err := parseWithLocation("UTC", "2024-06-15 10:30:45")
	if err != nil {
		t.FailNow()
	}
	gotm := NewTm(tm)
	s := Strftime("%A, %B %d, %Y %H:%M:%S", gotm)
	if len(s) == 0 {
		t.FailNow()
	}
}

func TestNewTmNow(t *testing.T) {
	now := time.Now()
	tm := NewTm(now)
	_ = tm
}

func TestParseUtcTimeGpsEpoch(t *testing.T) {
	ep, _ := ParseUtcTime("1980/01/06 00:00:00")
	if ep[0] != 1980 || ep[1] != 1 || ep[2] != 6 {
		t.FailNow()
	}
}
