package rtk

// #cgo CFLAGS: -I ./  -I ./libs
// #cgo CXXFLAGS: -I ./ -I ./libs
/*
#include <stdlib.h>
#include <time.h>
#include <string.h>
#include <ctype.h>
#include <stdio.h>
#include <rtklib.h>

void ppk_raw_to_rindex(gtime_t gpst, int format, const char *bin,
                       const char *ofile, const char *nfile,
                       const char *gfile) {
  rnxopt_t rnxopt = {0};
  int i;
  char file[1024], *outfile[6], ofile_[6][1024] = {""};
  for (i = 0; i < 6; i++)
    outfile[i] = ofile_[i];
  strncpy(file, bin, sizeof(file) - 1);
  file[sizeof(file) - 1] = '\0';
  rnxopt.rnxver = RNX3VER;
  strncpy(outfile[0], ofile, 1023);
  outfile[0][1023] = '\0';
  strncpy(outfile[1], nfile, 1023);
  outfile[1][1023] = '\0';
  if (gfile[0] != '\0') {
    strncpy(outfile[2], gfile, 1023);
    outfile[2][1023] = '\0';
  }
  rnxopt.trtcm = gpst;
  rnxopt.navsys = 0x3;
  rnxopt.obstype = 0xF;
  rnxopt.freqtype = 0x3;
  convrnx(format, &rnxopt, file, outfile);
}
*/
import "C"
import "unsafe"

type Format int

const (
	FormatRTCM2   Format = 0  // RTCM 2
	FormatRTCM3   Format = 1  // RTCM 3
	FormatOEM4    Format = 2  // NovAtel OEMV/4
	FormatOEM3    Format = 3  // NovAtel OEM3
	FormatUBX     Format = 4  // u-blox LEA-*T
	FormatSS2     Format = 5  // NovAtel Superstar II
	FormatCRES    Format = 6  // Hemisphere
	FormatSTQ     Format = 7  // SkyTraq S1315F
	FormatJAVAD   Format = 8  // JAVAD GRIL/GREIS
	FormatNVS     Format = 9  // NVS NVC08C
	FormatBINEX   Format = 10 // BINEX
	FormatRT17    Format = 11 // Trimble RT17
	FormatSEPT    Format = 12 // Septentrio
	FormatRINEX   Format = 13 // RINEX
	FormatSP3     Format = 14 // SP3
	FormatRNXCLK  Format = 15 // RINEX CLK
	FormatSBAS    Format = 16 // SBAS messages
	FormatNMEA    Format = 17 // NMEA 0183
)

func RawToRIndex(gpst GTime, format Format, binfile, ofile, nfile, gfile string) {
	cbinfile := C.CString(binfile)
	cofile := C.CString(ofile)
	cnfile := C.CString(nfile)
	cgfile := C.CString(gfile)

	defer C.free(unsafe.Pointer(cbinfile))
	defer C.free(unsafe.Pointer(cofile))
	defer C.free(unsafe.Pointer(cnfile))
	defer C.free(unsafe.Pointer(cgfile))

	C.ppk_raw_to_rindex(gpst.t, C.int(format), cbinfile, cofile, cnfile, cgfile)
}
