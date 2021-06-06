package rtk

/*
#include <stdlib.h>
#include <time.h>
#include <string.h>
#include <ctype.h>
#include <stdio.h>
#include <rtklib.h>

void ppk_raw_to_rindex(gtime_t gpst, const char *bin,
                       const char *ofile, const char *nfile,
                       const char *gfile) {
  rnxopt_t rnxopt = {0};
  int i;
  int format = STRFMT_RTCM3;
  char file[1024], *outfile[6], ofile_[6][1024] = {""}, *p;
  char buff[256], tstr[32];
  for (i = 0; i < 6; i++)
    outfile[i] = ofile_[i];
  strcpy(file, bin);
  rnxopt.rnxver = RNX3VER;
  strcpy(outfile[0], ofile);
  strcpy(outfile[1], nfile);
  if (gfile != "") {
    strcpy(outfile[2], gfile);
  }
  rnxopt.trtcm = gpst;
  rnxopt.navsys = 0x3;
  rnxopt.obstype = 0xF;
  rnxopt.freqtype = 0x3;
  convrnx(format, &rnxopt, file, outfile);
}
*/
import "C"
