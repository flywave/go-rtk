package rtk

/*
#include <stdlib.h>
#include <time.h>
#include <string.h>
#include <ctype.h>
#include <stdio.h>

double parse_time_to_utc(const char *inbuff, struct tm *tt) {
  int getbuf, i;
  int year, mon, doy, hour, minute, second, tmsec;
  int iArg;

  char trsbuf[120];

  for (i = 0; i < strlen(inbuff); i++) {
    if (isupper(inbuff[i]))
      trsbuf[i] = tolower(inbuff[i]);
    else
      trsbuf[i] = inbuff[i];
  }

  trsbuf[i++] = '\0';

  year = 0;
  mon = 0;
  doy = 0;
  hour = 0;
  minute = 0;
  second = 0;
  tmsec = 0;

  iArg = sscanf(trsbuf, "%d/%d/%d %d:%d:%d.%d", &year, &mon, &doy, &hour,
                &minute, &second, &tmsec);

  tt->tm_year = year - 1900;
  tt->tm_mday = doy;
  tt->tm_mon = mon - 1;
  tt->tm_hour = hour;
  tt->tm_min = minute;
  tt->tm_sec = second;
  tt->tm_gmtoff = 0;

  return 1000 / tmsec;
}

*/
import "C"

import (
	"time"
	"unsafe"
)

type Tm struct {
	core C.struct_tm
}

func (tm Tm) native() *C.struct_tm {
	return &tm.core
}

func NewTm(t time.Time) Tm {
	var tm C.struct_tm
	tm.tm_sec = C.int(t.Second())
	tm.tm_min = C.int(t.Minute())
	tm.tm_hour = C.int(t.Hour())
	tm.tm_mday = C.int(t.Day())
	tm.tm_mon = C.int(t.Month() - 1)
	tm.tm_year = C.int(t.Year() - 1900)
	tm.tm_wday = C.int(t.Weekday())
	tm.tm_yday = C.int(t.YearDay() - 1)
	tm.tm_isdst = -1
	return Tm{
		core: tm,
	}
}

func ParseUtcTime(str string) (ep [6]float64, sec float64) {
	var tt C.struct_tm
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	sec = float64(C.parse_time_to_utc(cstr, &tt))

	ep[0] = float64(tt.tm_year) + 1900
	ep[1] = float64(tt.tm_mon) + 1
	ep[2] = float64(tt.tm_mday)
	ep[3] = float64(tt.tm_hour)
	ep[4] = float64(tt.tm_min)
	ep[5] = float64(tt.tm_sec)

	return
}

const initialBufSize = 256

func Strftime(format string, tm Tm) (s string) {
	if format == "" {
		return
	}

	fmt := C.CString(format)
	defer C.free(unsafe.Pointer(fmt))

	for size := initialBufSize; ; size *= 2 {
		buf := (*C.char)(C.malloc(C.size_t(size))) // can panic
		defer C.free(unsafe.Pointer(buf))
		n := C.strftime(buf, C.size_t(size), fmt, tm.native())
		if n == 0 {
			if size > 20*len(format) {
				return
			}
		} else if int(n) < size {
			s = C.GoStringN(buf, C.int(n))
			return
		}
	}
}
