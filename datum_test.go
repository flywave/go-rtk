package rtk

import (
	"testing"
)

func TestDatumConstants(t *testing.T) {
	if WGS84 != 0 {
		t.FailNow()
	}
	if CGCS2000 != 1 {
		t.FailNow()
	}
}

func TestMSLToWGS84(t *testing.T) {
	h := MSLToWGS84(100.0, 120.0, 30.0, 0)
	if h == 0 {
		t.FailNow()
	}
}

func TestWGS84ToMSL(t *testing.T) {
	h := WGS84ToMSL(120.0, 30.0, 100.0, 0)
	if h == 0 {
		t.FailNow()
	}
}

func TestHAEToMSL(t *testing.T) {
	h := HAEToMSL(120.0, 30.0, 100.0, 0.0, 0)
	if h == 0 {
		t.FailNow()
	}
}

func TestMSLToHAE(t *testing.T) {
	h := MSLToHAE(100.0, 120.0, 30.0, 0.0, 0)
	if h == 0 {
		t.FailNow()
	}
}
