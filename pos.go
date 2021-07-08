package rtk

// #include <stdlib.h>
// #include <string.h>
// #include <rtklib.h>
// #cgo CFLAGS: -I ./  -I ./lib
// #cgo CXXFLAGS: -I ./ -I ./lib
import "C"
import (
	"reflect"
	"unsafe"
)

type Pos struct {
	Gpst      GTime
	Latitude  float64
	Longitude float64
	Height    float64
	Q         uint8
	Ns        uint8
	CovEE     float32
	CovNN     float32
	CovUU     float32
	CovEN     float32
	CovNU     float32
	CovUE     float32
	Age       float32
	Ratio     float32
}

type unsafeSol struct {
	time  C.gtime_t
	rr_0  float64
	rr_1  float64
	rr_2  float64
	rr_3  float64
	rr_4  float64
	rr_5  float64
	qr_0  float32
	qr_1  float32
	qr_2  float32
	qr_3  float32
	qr_4  float32
	qr_5  float32
	dtr_0 float64
	dtr_1 float64
	dtr_2 float64
	dtr_3 float64
	dtr_4 float64
	dtr_5 float64
	_type uint8
	stat  uint8
	ns    uint8
	age   float32
	ratio float32
}

type TimeRange [2]GTime

func ReadPos(pospath string) ([]Pos, TimeRange) {
	var buf C.solbuf_t
	cpospath := C.CString(pospath)
	defer C.free(unsafe.Pointer(cpospath))
	paths := []*C.char{cpospath}
	ret := int(C.readsol(&paths[0], 1, &buf))
	defer C.freesolbuf(&buf)
	if ret == 0 {
		return nil, TimeRange{}
	}
	count := int(buf.n)

	var cpossSlice []unsafeSol
	cposHeader := (*reflect.SliceHeader)((unsafe.Pointer(&cpossSlice)))
	cposHeader.Cap = int(count)
	cposHeader.Len = int(count)
	cposHeader.Data = uintptr(unsafe.Pointer(buf.data))

	poss := make([]Pos, count)
	for i := range poss {
		sol := cpossSlice[i]
		poss[i].Gpst.t = sol.time

		poss[i].Latitude = sol.rr_0
		poss[i].Longitude = sol.rr_1
		poss[i].Height = sol.rr_2

		poss[i].Q = sol.stat
		poss[i].Ns = sol.ns

		poss[i].CovEE = sol.qr_0
		poss[i].CovNN = sol.qr_1
		poss[i].CovUU = sol.qr_2
		poss[i].CovEN = sol.qr_3
		poss[i].CovNU = sol.qr_4
		poss[i].CovUE = sol.qr_5

		poss[i].Age = sol.age
		poss[i].Ratio = sol.ratio
	}
	start, end := int(buf.start), int(buf.end)
	return poss, TimeRange{poss[start].Gpst, poss[end].Gpst}
}
