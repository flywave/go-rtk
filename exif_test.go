package rtk

import (
	"os"
	"testing"
)

func TestReadExifXMP(t *testing.T) {
	f, err := os.Open("./testdata/dji_rtk.exif")
	if err != nil {
		t.FailNow()
	}
	defer f.Close()
	err, fields := ReadExifXMP(f)
	if err != nil {
		t.FailNow()
	}
	if fields == nil {
		t.FailNow()
	}
	if len(fields) == 0 {
		t.FailNow()
	}
}

func TestReadExifXMPHasTags(t *testing.T) {
	f, err := os.Open("./testdata/dji_rtk.exif")
	if err != nil {
		t.FailNow()
	}
	defer f.Close()
	err, fields := ReadExifXMP(f)
	if err != nil {
		t.FailNow()
	}
	v, ok := fields["drone-dji: RtkFlag"]
	if !ok {
		t.FailNow()
	}
	if v != "50" {
		t.FailNow()
	}
}

func TestReadExifXMPRtkStd(t *testing.T) {
	f, err := os.Open("./testdata/dji_rtk.exif")
	if err != nil {
		t.FailNow()
	}
	defer f.Close()
	err, fields := ReadExifXMP(f)
	if err != nil {
		t.FailNow()
	}
	_, ok := fields["drone-dji: RtkStdLon"]
	if !ok {
		t.FailNow()
	}
}

func TestReadExifXMPInvalid(t *testing.T) {
	err, _ := ReadExifXMP(nil)
	if err == nil {
		t.FailNow()
	}
}
