package rtk

import (
	"os"
	"strings"
	"testing"
)

func TestMarkRead(t *testing.T) {
	f, err := os.Open("./testdata/100_0004_Timestamp.MRK")
	if err != nil {
		t.FailNow()
	}
	mrks, err := ReadMRK(f)
	if err != nil {
		t.FailNow()
	}
	if len(mrks) == 0 {
		t.FailNow()
	}
}

func TestMRKGetGpst(t *testing.T) {
	f, err := os.Open("./testdata/100_0004_Timestamp.MRK")
	if err != nil {
		t.FailNow()
	}
	mrks, err := ReadMRK(f)
	if err != nil || len(mrks) == 0 {
		t.FailNow()
	}
	gt := mrks[0].GetGpst()
	ep := gt.Epoch()
	if ep[0] < 2020 || ep[0] > 2030 {
		t.FailNow()
	}
}

func TestMRKFieldValues(t *testing.T) {
	f, err := os.Open("./testdata/100_0004_Timestamp.MRK")
	if err != nil {
		t.FailNow()
	}
	mrks, err := ReadMRK(f)
	if err != nil || len(mrks) == 0 {
		t.FailNow()
	}
	m := mrks[0]
	if m.Sequence != 1 {
		t.FailNow()
	}
	if m.Week.w != 2024 {
		t.FailNow()
	}
	if m.State.e != 16 {
		t.FailNow()
	}
}

func TestPhaseCompUnmarshal(t *testing.T) {
	var pc PhaseComp
	err := pc.UnmarshalCSV("  -26,N")
	if err != nil {
		t.FailNow()
	}
	if pc.v != -26 {
		t.FailNow()
	}
	if pc.d != "N" {
		t.FailNow()
	}
	s := pc.String()
	if !strings.Contains(s, "-26") {
		t.FailNow()
	}
}

func TestAngleUnmarshal(t *testing.T) {
	var a Angle
	err := a.UnmarshalCSV("-36.41144502,Lat")
	if err != nil {
		t.FailNow()
	}
	if a.v != -36.41144502 {
		t.FailNow()
	}
	if a.d != "Lat" {
		t.FailNow()
	}
}

func TestWeekUnmarshal(t *testing.T) {
	var w Week
	err := w.UnmarshalCSV("[2024]")
	if err != nil {
		t.FailNow()
	}
	if w.w != 2024 {
		t.FailNow()
	}
}

func TestRtkStateUnmarshal(t *testing.T) {
	var rs RtkState
	err := rs.UnmarshalCSV("16,Q")
	if err != nil {
		t.FailNow()
	}
	if rs.e != 16 {
		t.FailNow()
	}
	if rs.f != "Q" {
		t.FailNow()
	}
}

func TestStdUnmarshal(t *testing.T) {
	var s Std
	err := s.UnmarshalCSV("1.492559, 1.368525, 3.180128")
	if err != nil {
		t.FailNow()
	}
}

func TestPosRead(t *testing.T) {
	pos, _ := ReadPos("./testdata/test.pos")
	if len(pos) == 0 {
		t.FailNow()
	}
}
