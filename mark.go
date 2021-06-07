package rtk

type PhaseComp struct {
	v float64
	d byte
}

type MRK struct {
	Photo               int32
	Time                float64
	Week                int32
	PhaseCompNs         PhaseComp
	PhaseCompEw         PhaseComp
	PhaseCompV          PhaseComp
	Latitude            float64
	Longitude           float64
	Altitude            float64
	StdLatitude         float64
	StdLongitude        float64
	StdAltitude         float64
	ClosestId           int32
	ClosestTime         GTime
	ClosestLatitude     float64
	ClosestLongitude    float64
	ClosestAltitude     float64
	SndClosestId        int32
	SndClosestTime      GTime
	SndClosestLatitude  float64
	SndClosestLongitude float64
	SndClosestAltitude  float64
	GpstDiff            float64
	InterpLatitude      float64
	InterpLongitude     float64
	InterpAltitude      float64
	PhaseCompNsDeg      float64
	PhaseCompEwDeg      float64
	PhaseCompVM         float64
	UpdatedLatitude     float64
	UpdatedLongitude    float64
	UpdatedAltitude     float64
	DiffLatitude        float64
	DiffLongitude       float64
	DiffAltitude        float64
}
