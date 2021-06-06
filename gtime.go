package rtk

// #include <stdlib.h>
// #include <string.h>
// #include <rtklib.h>
// #cgo CFLAGS: -I ./  -I ./lib
// #cgo CXXFLAGS: -I ./ -I ./lib
// #cgo linux LDFLAGS:  -L ./lib -Wl,--start-group  -lstdc++ -lm -pthread -ldl -lrtklib -Wl,--end-group
// #cgo windows LDFLAGS: -L ./lib -lrtklib
// #cgo darwin LDFLAGS: -L　./lib -lrtklib
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

type GTime struct {
	t C.gtime_t
}

func NewGTimeFromStr(str string) *GTime {
	var ct C.gtime_t
	cstr := C.CString(str)
	C.str2time(cstr, C.int(0), C.int(len(str)), &ct)
	return &GTime{t: ct}
}

func NewGTimeFromEpoch(ep [6]float64) *GTime {
	ct := C.epoch2time((*C.double)(unsafe.Pointer(&ep[0])))
	return &GTime{t: ct}
}

func NewGTimeFromGPSTime(week int32, sec float64) *GTime {
	ct := C.gpst2time(C.int(week), C.double(sec))
	return &GTime{t: ct}
}

func NewGTimeFromGalileoTime(week int32, sec float64) *GTime {
	ct := C.gst2time(C.int(week), C.double(sec))
	return &GTime{t: ct}
}

func NewGTimeFromBDTime(week int32, sec float64) *GTime {
	ct := C.bdt2time(C.int(week), C.double(sec))
	return &GTime{t: ct}
}

func NewGPSTFromUTC(utc *GTime) *GTime {
	ct := C.utc2gpst(utc.t)
	return &GTime{t: ct}
}

func NewUtcTimeFromTime(t time.Time) *GTime {
	var ep [6]float64

	ep[0] = float64(t.Year())
	ep[1] = float64(t.Month())
	ep[2] = float64(t.Day())
	ep[3] = float64(t.Hour())
	ep[4] = float64(t.Minute())
	ep[5] = float64(t.Second())

	g := NewGTimeFromEpoch(ep)
	return g
}

func NewGPSTTimeFromTime(t time.Time) *GTime {
	g := NewUtcTimeFromTime(t)
	return NewGPSTFromUTC(g)
}

func NewUtcTime(str string) *GTime {
	ep, sec := ParseUtcTime(str)
	g := NewGTimeFromEpoch(ep)
	g.Add(sec)
	return g
}

func NewGPSTTime(str string) *GTime {
	g := NewUtcTime(str)
	return NewGPSTFromUTC(g)
}

const (
	TIME_LAYOUT = "2006-01-02 15:04:05"
)

func parseWithLocation(name string, timeStr string) (time.Time, error) {
	locationName := name
	if l, err := time.LoadLocation(locationName); err != nil {
		println(err.Error())
		return time.Time{}, err
	} else {
		lt, _ := time.ParseInLocation(TIME_LAYOUT, timeStr, l)
		fmt.Println(locationName, lt)
		return lt, nil
	}
}

func NewUtcTimeFromLocal(timeStr string, name string) *GTime {
	t, err := parseWithLocation(name, timeStr)
	if err != nil {
		return nil
	}
	return NewUtcTimeFromTime(t)
}

func NewGPSTTimeFromLocal(timeStr string, name string) *GTime {
	utc := NewUtcTimeFromLocal(timeStr, name)
	return NewGPSTFromUTC(utc)
}

func NewGPSTTimeFromLocalTM(t Tm, sec float64) *GTime {
	utc := NewUTCFromLocalTM(t, sec)
	return NewGPSTFromUTC(utc)
}

func getLocalTimeZoneOffset() (name string, offset int) {
	return time.Now().Zone()
}

func NewUTCFromLocalTM(tt Tm, sec float64) *GTime {
	var ep [6]float64

	_, offset := getLocalTimeZoneOffset()

	ep[0] = float64(tt.core.tm_year) + 1900
	ep[1] = float64(tt.core.tm_mon) + 1
	ep[2] = float64(tt.core.tm_mday)
	ep[3] = float64(tt.core.tm_hour)
	ep[4] = float64(tt.core.tm_min)
	ep[5] = float64(tt.core.tm_sec) + float64(offset)

	g := NewGTimeFromEpoch(ep)
	g.Add(sec)
	return g
}

func NewUtcTimeFromCurrentLocal(timeStr string) *GTime {
	zone, _ := getLocalTimeZoneOffset()
	t, err := parseWithLocation(zone, timeStr)
	if err != nil {
		return nil
	}
	return NewUtcTimeFromTime(t)
}

func NewGPSTTimeFromCurrentLocal(timeStr string) *GTime {
	utc := NewUtcTimeFromCurrentLocal(timeStr)
	return NewGPSTFromUTC(utc)
}

func Current() *GTime {
	ct := C.timeget()
	return &GTime{t: ct}
}

func (g *GTime) Epoch() (ep [6]float64) {
	C.time2epoch(g.t, (*C.double)(unsafe.Pointer(&ep[0])))
	return
}

func (g *GTime) GpsTime() (week int32, sec float64) {
	var cweek C.int
	sec = float64(C.time2gpst(g.t, &cweek))
	week = int32(cweek)
	return
}

func (g *GTime) GalileoTime() (week int32, sec float64) {
	var cweek C.int
	sec = float64(C.time2gst(g.t, &cweek))
	week = int32(cweek)
	return
}

func (g *GTime) BDTime() (week int32, sec float64) {
	var cweek C.int
	sec = float64(C.time2bdt(g.t, &cweek))
	week = int32(cweek)
	return
}

func (g *GTime) UTC() *GTime {
	ct := C.gpst2utc(g.t)
	return &GTime{t: ct}
}

func (g *GTime) Time() int {
	return int(g.t.time)
}

func (g *GTime) Sec() float64 {
	return float64(g.t.sec)
}

func (g *GTime) Add(sec float64) {
	g.t = C.timeadd(g.t, C.double(sec))
}

func (g *GTime) Diff(o *GTime) float64 {
	return float64(C.timediff(g.t, o.t))
}

func (g *GTime) ToString(n int) string {
	return C.GoString(C.time_str(g.t, C.int(n)))
}

func (g *GTime) DayOfYear(n int) float64 {
	return float64(C.time2doy(g.t))
}
