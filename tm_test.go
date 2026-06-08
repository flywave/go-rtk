package rtk

import (
	"testing"
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
