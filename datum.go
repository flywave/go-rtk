package rtk

import (
	"path/filepath"
	"runtime"

	"github.com/flywave/go-geoid"
	"github.com/flywave/go-proj"
)

func getCurrentDir() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}

func init() {
	dir := getCurrentDir()
	proj.SetFinder([]string{filepath.Join(dir, "../proj_data")})
}

type Datum uint32

const (
	WGS84    Datum = 0
	CGCS2000 Datum = 1
)

const (
	wgs84_ecef_srs    = "+proj=geocent +ellps=WGS84 +datum=WGS84 +units=m +no_defs"
	wgs84_lla_srs     = "+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs"
	cgcs2000_ecef_srs = "+proj=geocent +ellps=GRS80 +units=m +no_defs"
	cgcs2000_lla_srs  = "+proj=longlat +ellps=GRS80 +no_defs"
)

var (
	wgs84_ecef, _    = proj.NewProj(wgs84_ecef_srs)
	wgs84_lla, _     = proj.NewProj(wgs84_lla_srs)
	cgcs2000_ecef, _ = proj.NewProj(cgcs2000_ecef_srs)
	cgcs2000_lla, _  = proj.NewProj(cgcs2000_lla_srs)
)

func transform(src *proj.Proj, dst *proj.Proj, v [3]float64) [3]float64 {
	if src.IsLatLong() {
		v[0] = DegreeToRadian(v[0])
		v[1] = DegreeToRadian(v[1])
	}
	v[0], v[1], v[2], _ = proj.Transform3(src, dst, v[0], v[1], v[2])
	if dst.IsLatLong() {
		v[0] = RadianToDegree(v[0])
		v[1] = RadianToDegree(v[1])
	}
	return v
}

func transformFromLLA(dst *proj.Proj, v [3]float64, d Datum) [3]float64 {
	if d == WGS84 {
		return transform(wgs84_lla, dst, v)
	} else {
		return transform(cgcs2000_lla, dst, v)
	}
}

func is_datums_wgs84(dst *proj.Proj) bool {
	return wgs84_lla.CompareDatums(dst)
}

func is_datums_cgcs2000(dst *proj.Proj) bool {
	return cgcs2000_lla.CompareDatums(dst)
}

func isDatums(dst *proj.Proj, d Datum) bool {
	if d == WGS84 {
		return is_datums_wgs84(dst)
	} else {
		return is_datums_cgcs2000(dst)
	}
}

func init() {
	dir := getCurrentDir()
	geoid.SetGeoidPath(filepath.Join(dir, "../geoid_data"))
}

func geoid84_30() *geoid.Geoid {
	return geoid.NewGeoid(geoid.EGM84, false)
}

func geoid2008_5() *geoid.Geoid {
	return geoid.NewGeoid(geoid.EGM2008, false)
}

func geoid96_15() *geoid.Geoid {
	return geoid.NewGeoid(geoid.EGM96, false)
}

func getGeoid(g geoid.VerticalDatum) *geoid.Geoid {
	switch g {
	case geoid.EGM84:
		return geoid84_30()
	case geoid.EGM2008:
		return geoid2008_5()
	case geoid.EGM96:
		return geoid96_15()
	default:
		return geoid84_30()
	}
}

func MSLToWGS84(h, lon, lat float64, g geoid.VerticalDatum) float64 {
	return getGeoid(g).ConvertHeight(lat, lon, h, geoid.GEOIDTOELLIPSOID)
}

func WGS84ToMSL(lon, lat, altitude float64, g geoid.VerticalDatum) float64 {
	return getGeoid(g).ConvertHeight(lat, lon, altitude, geoid.ELLIPSOIDTOGEOID)
}

func HAEToMSL(lon, lat, altitude, ellipsoidOffset float64, g geoid.VerticalDatum) float64 {
	return getGeoid(g).ConvertHeight(lat, lon, altitude+ellipsoidOffset, geoid.ELLIPSOIDTOGEOID)
}

func MSLToHAE(h, lon, lat, ellipsoidOffset float64, g geoid.VerticalDatum) float64 {
	return getGeoid(g).ConvertHeight(lat, lon, h, geoid.GEOIDTOELLIPSOID) + ellipsoidOffset
}
