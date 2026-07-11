package rtk

import (
	"os"
	"strings"
	"testing"
)

func TestReadExifXMPInvalid(t *testing.T) {
	err, _ := ReadExifXMP(nil)
	if err == nil {
		t.FailNow()
	}
}

func TestReadExifXMPEmptyReader(t *testing.T) {
	err, _ := ReadExifXMP(strings.NewReader(""))
	if err == nil {
		t.FailNow()
	}
}

func TestReadExifXMPNoEXIF(t *testing.T) {
	err, _ := ReadExifXMP(strings.NewReader("not an exif file"))
	if err == nil {
		t.FailNow()
	}
}

func TestReadExifXMPTextFile(t *testing.T) {
	f, err := os.Open("./testdata/dji_rtk.exif")
	if err != nil {
		t.FailNow()
	}
	defer f.Close()
	err, fields := ReadExifXMP(f)
	if err == nil {
		t.FailNow()
	}
	if fields != nil {
		t.FailNow()
	}
}

func TestReadExifXMPTooFewPackets(t *testing.T) {
	err, fields := ReadExifXMP(strings.NewReader("<?xpacket begin=\"\" ?>"))
	if err == nil {
		t.FailNow()
	}
	if fields != nil {
		t.FailNow()
	}
}

func TestReadExifXMPNoPackets(t *testing.T) {
	err, fields := ReadExifXMP(strings.NewReader("some data without xpacket markers"))
	if err == nil {
		t.FailNow()
	}
	if fields != nil {
		t.FailNow()
	}
}
