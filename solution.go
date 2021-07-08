package rtk

/*
#include <stdlib.h>
#include <time.h>
#include <string.h>
#include <ctype.h>
#include <stdio.h>
#include <rtklib.h>

#define INFILEMAX 5
#define BUFSIZE 1024

int showmsg(char *format, ...) { return 0; }

void settspan(gtime_t ts, gtime_t te) {}

void settime(gtime_t time) {}

void init_popt(prcopt_t *opt) {
  opt->mode = PMODE_PPP_STATIC;
  opt->nf = 3;
  opt->navsys = SYS_GPS | SYS_GLO | SYS_GAL | SYS_CMP;
  opt->elmin = 10.0 * D2R;
  opt->sateph = EPHOPT_PREC;
  opt->modear = 0;
  opt->glomodear = 0;
  opt->bdsmodear = 0;
  opt->ionoopt = IONOOPT_IFLC;
  opt->tropopt = TROPOPT_EST;
}

void init_sopt(solopt_t *opt) {
  opt->posf = SOLF_LLH;
  opt->times = TIMES_GPST;
  opt->timef = 1;
  opt->outhead = 1;
}

void init_in_file(char *infile[], int *n, const char *rofile,
                  const char *bofile, const char *navfile) {
  for (int i = 0; i < INFILEMAX; i++) {
    infile[i] = (char *)malloc(sizeof(char) * BUFSIZE);
    *infile[i] = '\0';
  }
  strcpy(infile[0], rofile);
  strcpy(infile[1], bofile);
  strcpy(infile[2], navfile);

  *n = 0;
  while (strlen(infile[*n]) && *n <= INFILEMAX) {
    *n++;
  }
}

void free_in_file(char *infile[]) {
  for (int i = 0; i < INFILEMAX; i++) {
    free(infile[i]);
  }
}

int ppk_solution(const char *rofile, const char *bofile,
                 const char *navfile, const char *posfile) {
  filopt_t filopt = {""};
  prcopt_t prcopt = prcopt_default;
  solopt_t solopt = solopt_default;

  int n;
  char *infile[5];
  char outfile[1024] = "";
  strcpy(outfile, posfile);

  resetsysopts();
  getsysopts(&prcopt, &solopt, &filopt);

  init_popt(&prcopt);
  init_sopt(&solopt);
  init_in_file(infile, &n, rofile, bofile, navfile);

  gtime_t ts = {0}, te = {0};
  double ti = 0.0, tu = 0.0;

  if (postpos(ts, te, ti, tu, &prcopt, &solopt, &filopt, infile, n, outfile, "",
              "") == 1) {
    return -1;
  }

  free_in_file(infile);
  return 0;
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

func Solution(rofile string, bofile string, navfile string, posfile string) error {
	crofile := C.CString(rofile)
	cbofile := C.CString(bofile)
	cnavfile := C.CString(navfile)
	cposfile := C.CString(posfile)

	defer C.free(unsafe.Pointer(crofile))
	defer C.free(unsafe.Pointer(cbofile))
	defer C.free(unsafe.Pointer(cnavfile))
	defer C.free(unsafe.Pointer(cposfile))

	ret := int(C.ppk_solution(crofile, cbofile, cnavfile, cposfile))
	if ret < 0 {
		return errors.New("solutin error")
	}
	return nil
}
