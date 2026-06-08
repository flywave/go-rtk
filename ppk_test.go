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

func TestRadianToDegree(t *testing.T) {
	v := RadianToDegree(math.Pi)
	if math.Abs(v-180.0) > 1e-15 {
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
