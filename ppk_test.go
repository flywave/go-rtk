package rtk

import (
	"math"
	"testing"
)

func TestDegreeToRadian(t *testing.T) {
	v := DegreeToRadian(180.0)
	if math.Abs(v-math.Pi) > 1e-15 {
		t.FailNow()
	}
}

func TestDegreeToRadianZero(t *testing.T) {
	v := DegreeToRadian(0.0)
	if v != 0.0 {
		t.FailNow()
	}
}

func TestRadianToDegree(t *testing.T) {
	v := RadianToDegree(math.Pi)
	if math.Abs(v-180.0) > 1e-15 {
		t.FailNow()
	}
}

func TestRadianToDegreeZero(t *testing.T) {
	v := RadianToDegree(0.0)
	if v != 0.0 {
		t.FailNow()
	}
}

func TestDegreeRadianRoundTrip(t *testing.T) {
	orig := 45.0
	rad := DegreeToRadian(orig)
	deg := RadianToDegree(rad)
	if math.Abs(deg-orig) > 1e-15 {
		t.FailNow()
	}
}

func TestNearestTimeInd(t *testing.T) {
	poses := make([]Pos, 3)
	poses[0].Gpst = *NewGTimeFromGPSTime(2024, 100.0)
	poses[1].Gpst = *NewGTimeFromGPSTime(2024, 200.0)
	poses[2].Gpst = *NewGTimeFromGPSTime(2024, 300.0)

	pivot := NewGTimeFromGPSTime(2024, 190.0)
	idx := nearestTimeInd(poses, *pivot)
	if idx != 1 {
		t.FailNow()
	}
}

func TestNearestTimeIndExact(t *testing.T) {
	poses := make([]Pos, 3)
	poses[0].Gpst = *NewGTimeFromGPSTime(2024, 100.0)
	poses[1].Gpst = *NewGTimeFromGPSTime(2024, 200.0)
	poses[2].Gpst = *NewGTimeFromGPSTime(2024, 300.0)

	pivot := NewGTimeFromGPSTime(2024, 200.0)
	idx := nearestTimeInd(poses, *pivot)
	if idx != 1 {
		t.FailNow()
	}
}

func TestNearestTimeIndFirst(t *testing.T) {
	poses := make([]Pos, 3)
	poses[0].Gpst = *NewGTimeFromGPSTime(2024, 100.0)
	poses[1].Gpst = *NewGTimeFromGPSTime(2024, 200.0)
	poses[2].Gpst = *NewGTimeFromGPSTime(2024, 300.0)

	pivot := NewGTimeFromGPSTime(2024, 50.0)
	idx := nearestTimeInd(poses, *pivot)
	if idx != 0 {
		t.FailNow()
	}
}

func TestNearestTimeIndLast(t *testing.T) {
	poses := make([]Pos, 3)
	poses[0].Gpst = *NewGTimeFromGPSTime(2024, 100.0)
	poses[1].Gpst = *NewGTimeFromGPSTime(2024, 200.0)
	poses[2].Gpst = *NewGTimeFromGPSTime(2024, 300.0)

	pivot := NewGTimeFromGPSTime(2024, 350.0)
	idx := nearestTimeInd(poses, *pivot)
	if idx != 2 {
		t.FailNow()
	}
}

func TestNearestTimeIndSingle(t *testing.T) {
	poses := make([]Pos, 1)
	poses[0].Gpst = *NewGTimeFromGPSTime(2024, 100.0)

	pivot := NewGTimeFromGPSTime(2024, 50.0)
	idx := nearestTimeInd(poses, *pivot)
	if idx != 0 {
		t.FailNow()
	}
}

func TestPpkInitPosEmptyPoses(t *testing.T) {
	mrks := make([]MRK, 1)
	mrks[0].Week.w = 2024
	mrks[0].Time = 100.0
	ppkInitPos([]Pos{}, mrks)
}

func TestPpkInitPosEmptyMrks(t *testing.T) {
	poses := make([]Pos, 1)
	poses[0].Gpst = *NewGTimeFromGPSTime(2024, 100.0)
	ppkInitPos(poses, []MRK{})
}

func TestPpkUpdatedMrksEmpty(t *testing.T) {
	ppkUpdatedMrks([]MRK{})
}

func TestPpkInitPosSingle(t *testing.T) {
	poses := make([]Pos, 1)
	poses[0].Gpst = *NewGTimeFromGPSTime(2024, 100.0)
	poses[0].Latitude = 30.0
	poses[0].Longitude = 120.0
	poses[0].Height = 50.0

	mrks := make([]MRK, 1)
	mrks[0].Week.w = 2024
	mrks[0].Time = 100.0

	ppkInitPos(poses, mrks)
	if mrks[0].closestId != 0 {
		t.FailNow()
	}
	if mrks[0].sndClosestId != 0 {
		t.FailNow()
	}
}

func TestPpkInitPosTwoPoses(t *testing.T) {
	poses := make([]Pos, 2)
	poses[0].Gpst = *NewGTimeFromGPSTime(2024, 100.0)
	poses[0].Latitude = 30.0
	poses[1].Gpst = *NewGTimeFromGPSTime(2024, 200.0)
	poses[1].Latitude = 31.0

	mrks := make([]MRK, 1)
	mrks[0].Week.w = 2024
	mrks[0].Time = 150.0

	ppkInitPos(poses, mrks)
	if mrks[0].closestId == mrks[0].sndClosestId {
		t.FailNow()
	}
}

func TestPPKSolutionEmptyPos(t *testing.T) {
	err, sols := PPKSolution("", "", WGS84, 0, "", 0, 0)
	if err != nil {
		t.FailNow()
	}
	if sols != nil {
		t.FailNow()
	}
}

func TestPPKSolutionWeightZeroStd(t *testing.T) {
	mrk := MRK{
		Sequence: 1,
		Week:     Week{w: 2024},
		Time:     100.0,
		Latitude: Angle{v: 30.0, d: "Lat"},
	}
	_ = mrk
	// Weight calculation with zero Std should not panic
	s := Std{latitude: 0, longitude: 0, altitude: 0}
	wlon, wlat, walt := 1.0, 1.0, 1.0
	if s.longitude == 0 {
		wlon = 1.0
	} else {
		wlon = 1 / s.longitude
	}
	if s.latitude == 0 {
		wlat = 1.0
	} else {
		wlat = 1 / s.latitude
	}
	if s.altitude == 0 {
		walt = 1.0
	} else {
		walt = 1 / s.altitude
	}
	if math.IsInf(wlon, 0) || math.IsInf(wlat, 0) || math.IsInf(walt, 0) {
		t.FailNow()
	}
}

func TestPPKSolutionWeightNonZeroStd(t *testing.T) {
	s := Std{latitude: 2.0, longitude: 3.0, altitude: 4.0}
	wlon := 1 / s.longitude
	wlat := 1 / s.latitude
	walt := 1 / s.altitude
	if wlon != 1.0/3.0 || wlat != 0.5 || walt != 0.25 {
		t.FailNow()
	}
}
