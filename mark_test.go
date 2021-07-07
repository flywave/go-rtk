package rtk

import (
	"os"
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

func TestPosRead(t *testing.T) {
	f, err := os.Open("./testdata/test.pos")
	if err != nil {
		t.FailNow()
	}
	pos, err := ReadPos(f)
	if err != nil {
		t.FailNow()
	}
	if len(pos) == 0 {
		t.FailNow()
	}
}
