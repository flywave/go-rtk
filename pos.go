package rtk

import (
	_ "github.com/jszwec/csvutil"
)

type Pos struct {
	Gpst      GTime
	Latitude  float64
	Longitude float64
	Height    float64
	Q         int8
	Ns        int16
	Sdn       float64
	Sde       float64
	Sdu       float64
	Sdne      float64
	Sdeu      float64
	Sdun      float64
	Age       float64
	Ratio     float64
}
