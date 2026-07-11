package rtk

import (
	"testing"
)

func TestFormatConstants(t *testing.T) {
	tests := []struct {
		got  Format
		want int
	}{
		{FormatRTCM2, 0},
		{FormatRTCM3, 1},
		{FormatOEM4, 2},
		{FormatOEM3, 3},
		{FormatUBX, 4},
		{FormatSS2, 5},
		{FormatCRES, 6},
		{FormatSTQ, 7},
		{FormatJAVAD, 8},
		{FormatNVS, 9},
		{FormatBINEX, 10},
		{FormatRT17, 11},
		{FormatSEPT, 12},
		{FormatRINEX, 13},
		{FormatSP3, 14},
		{FormatRNXCLK, 15},
		{FormatSBAS, 16},
		{FormatNMEA, 17},
	}
	for _, tt := range tests {
		if int(tt.got) != tt.want {
			t.FailNow()
		}
	}
}

func TestFormatCount(t *testing.T) {
	all := []Format{
		FormatRTCM2, FormatRTCM3, FormatOEM4, FormatOEM3,
		FormatUBX, FormatSS2, FormatCRES, FormatSTQ,
		FormatJAVAD, FormatNVS, FormatBINEX, FormatRT17,
		FormatSEPT, FormatRINEX, FormatSP3, FormatRNXCLK,
		FormatSBAS, FormatNMEA,
	}
	if len(all) != 18 {
		t.FailNow()
	}
}

func TestFormatRoundTrip(t *testing.T) {
	if int(FormatRTCM2) != 0 || int(FormatNMEA) != 17 {
		t.FailNow()
	}
}
