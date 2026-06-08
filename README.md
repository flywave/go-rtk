# go-rtk

Go package wrapping [RTKLIB](https://github.com/tomojitakasu/RTKLIB) via CGo for PPK (Post-Processed Kinematic) GNSS post-processing. Works with drone image metadata, RINEX observation files, `.pos` solutions, and `.MRK` timestamp files.

## Package

```
import "github.com/flywave/go-rtk"
```

Single flat package `rtk` — no sub-packages.

## Dependencies

| Dependency | Source |
|---|---|
| RTKLIB (C) | `external/rtklib/` — prebuilt as `lib/librtklib.a` |
| [go-proj](https://github.com/flywave/go-proj) | coordinate transforms via PROJ |
| [go-geoid](https://github.com/flywave/go-geoid) | EGM84/96/2008 geoid height conversions |
| goexif | EXIF tag parsing |
| go-xmp | XMP metadata parsing |
| gocsv | CSV marshaling for `.MRK` files |

`go-proj` and `go-geoid` are local replaces in `go.mod` — they must exist as siblings at `../go-proj` and `../go-geoid`.

## Build

```sh
go build ./...
go test ./...
go vet ./...
```

`lib/librtklib.a` is checked in. To rebuild from source: `cmake external/rtklib && make`.

## Usage

### PPK processing

```go
poses, trange := rtk.ReadPos("solution.pos")
mrks, err := rtk.ReadMRK(mrkFile)

rtk.PPKInitPos(poses, mrks)
rtk.PPKUpdatedMrks(mrks)

sols, err := rtk.PPKSolution("solution.pos", "marks.MRK",
    rtk.WGS84, rtk.HAE, "+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs",
    rtk.HAE, 0.0)
```

### GNSS time handling

```go
gt := rtk.NewGTimeFromGPSTime(2024, 275295.301059)
week, sec := gt.GpsTime()
ep := gt.Epoch()      // [year, month, day, hour, min, sec]
utc := gt.UTC()
diff := gt.Diff(otherGT)
gt.Add(30.0)           // add seconds
```

### EXIF / XMP metadata

```go
f, _ := os.Open("image.jpg")
err, fields := rtk.ReadExifXMP(f)
// fields["drone-dji: RtkFlag"], fields["drone-dji: RtkStdLon"], etc.
```

### Coordinate transforms

```go
h := rtk.MSLToWGS84(100.0, 120.0, 30.0, geoid.EGM2008)
h = rtk.WGS84ToMSL(120.0, 30.0, 100.0, geoid.EGM96)
h = rtk.HAEToMSL(120.0, 30.0, 100.0, 0.0, geoid.EGM84)
```

### Solution processing (RINEX → .pos)

```go
err := rtk.Solution("rover.obs", "base.obs", "brdc.nav", "output.pos")
```

## File formats

| Format | Reader | Description |
|---|---|---|
| `.pos` | `ReadPos` | RTKLIB solution file (tab-separated, header lines start with `%`) |
| `.MRK` | `ReadMRK` | Tab-separated timestamp marks with GPS week/SOW, phase center offsets, and RTK state |
| `.exif` / XMP | `ReadExifXMP` | DJI drone image metadata (RtkFlag, RtkStd, absolute altitude) |
| `.obs` / `.nav` | `Solution` | RINEX observation and navigation files (via RTKLIB) |
| `.bin` | `RawToRIndex` | RTCM3 binary → RINEX conversion |

### `.MRK` columns

```
Sequence  GPSSecondOfWeek  GPSWeek  NorthOff  EastOff  VelOff  Latitude  Longitude  EllipsoidHeight  Std  RtkState
```

Tab-delimited, no header row.

## Data directories

- `proj_data/` — PROJ datum grids and transform files
- `geoid_data/` — EGM84/96/2008 geoid height grids (PGM format)
- `testdata/` — sample `.pos`, `.MRK`, `.exif`, and RINEX files for testing

Paths are set automatically via `init()` — no environment variables needed.

## Notes

- Go 1.13 compatible (no generics, no `errors.Is`/`errors.As`, no `os.ReadFile`)
- CGo with per-file `#cgo` directives in `gtime.go`, `pos.go`, `solution.go`, `ppkraw.go`
- Tests use `t.FailNow()` on failure (no error messages)
- `go.sum` has no entries for replace-only deps — CI may need `-mod=mod`
