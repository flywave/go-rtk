package rtk

import (
	_ "github.com/flywave/go-geoid"
	_ "github.com/flywave/go-proj"
)

type PPKOption struct {
	BaseAerialHeight  float32
	RoverAerialHeight float32
	RoverCameraOffset [3]float32
}
