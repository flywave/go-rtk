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
	f.Close()
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
	f.Close()
	if err != nil || len(mrks) == 0 {
		t.FailNow()
	}
	gt := mrks[0].GetGpst()
	ep := gt.Epoch()
	if ep[0] < 2018 || ep[0] > 2030 {
		t.FailNow()
	}
}

func TestMRKFieldValues(t *testing.T) {
	f, err := os.Open("./testdata/100_0004_Timestamp.MRK")
	if err != nil {
		t.FailNow()
	}
	mrks, err := ReadMRK(f)
	f.Close()
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

func TestPhaseCompUnmarshalLeadingSpace(t *testing.T) {
	var pc PhaseComp
	err := pc.UnmarshalCSV("  1.5,E")
	if err != nil {
		t.FailNow()
	}
	if pc.v != 1.5 || pc.d != "E" {
		t.FailNow()
	}
}

func TestPhaseCompMarshalRoundTrip(t *testing.T) {
	pc := PhaseComp{v: -26, d: "N"}
	s, err := pc.MarshalCSV()
	if err != nil {
		t.FailNow()
	}
	var pc2 PhaseComp
	err = pc2.UnmarshalCSV(s)
	if err != nil {
		t.FailNow()
	}
	if pc.v != pc2.v || pc.d != pc2.d {
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

func TestAngleUnmarshalLeadingSpace(t *testing.T) {
	var a Angle
	err := a.UnmarshalCSV("  30.0,Lon")
	if err != nil {
		t.FailNow()
	}
	if a.v != 30.0 || a.d != "Lon" {
		t.FailNow()
	}
}

func TestAngleMarshalRoundTrip(t *testing.T) {
	a := Angle{v: -36.41144502, d: "Lat"}
	s, err := a.MarshalCSV()
	if err != nil {
		t.FailNow()
	}
	var a2 Angle
	err = a2.UnmarshalCSV(s)
	if err != nil {
		t.FailNow()
	}
	if a.v != a2.v || a.d != a2.d {
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

func TestWeekUnmarshalWithSpace(t *testing.T) {
	var w Week
	err := w.UnmarshalCSV("  [2024]  ")
	if err != nil {
		t.FailNow()
	}
	if w.w != 2024 {
		t.FailNow()
	}
}

func TestWeekMarshalRoundTrip(t *testing.T) {
	w := Week{w: 2024}
	s, err := w.MarshalCSV()
	if err != nil {
		t.FailNow()
	}
	var w2 Week
	err = w2.UnmarshalCSV(s)
	if err != nil {
		t.FailNow()
	}
	if w.w != w2.w {
		t.FailNow()
	}
}

func TestWeekUnmarshalInvalid(t *testing.T) {
	var w Week
	err := w.UnmarshalCSV("2024")
	if err != nil {
		t.FailNow()
	}
	if w.w != 0 {
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

func TestRtkStateMarshalRoundTrip(t *testing.T) {
	rs := RtkState{e: 16, f: "Q"}
	s, err := rs.MarshalCSV()
	if err != nil {
		t.FailNow()
	}
	var rs2 RtkState
	err = rs2.UnmarshalCSV(s)
	if err != nil {
		t.FailNow()
	}
	if rs.e != rs2.e || rs.f != rs2.f {
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

func TestStdUnmarshalRoundTrip(t *testing.T) {
	s := Std{latitude: 1.492559, longitude: 1.368525, altitude: 3.180128}
	str, err := s.MarshalCSV()
	if err != nil {
		t.FailNow()
	}
	var s2 Std
	err = s2.UnmarshalCSV(str)
	if err != nil {
		t.FailNow()
	}
}

func TestStdUnmarshalError(t *testing.T) {
	var s Std
	err := s.UnmarshalCSV("invalid, 1.0, 2.0")
	if err == nil {
		t.FailNow()
	}
}

func TestPosRead(t *testing.T) {
	pos, _ := ReadPos("./testdata/test.pos")
	if len(pos) == 0 {
		t.FailNow()
	}
}

func TestReadMRKInvalidFile(t *testing.T) {
	_, err := ReadMRK(strings.NewReader("invalid\tdata"))
	if err == nil {
		t.FailNow()
	}
}

func TestMRKGetGpstMultiple(t *testing.T) {
	f, err := os.Open("./testdata/100_0004_Timestamp.MRK")
	if err != nil {
		t.FailNow()
	}
	mrks, err := ReadMRK(f)
	f.Close()
	if err != nil || len(mrks) < 2 {
		t.FailNow()
	}
	for i := range mrks {
		_ = mrks[i].GetGpst()
	}
}
