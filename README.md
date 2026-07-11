# go-rtk

Go package (`package rtk`) wrapping [RTKLIB](https://github.com/tomojitakasu/RTKLIB) via CGo for PPK (Post-Processed Kinematic) and RTK GNSS post-processing. Handles drone image metadata (EXIF/XMP), RINEX observation files, `.pos` solution files, and `.MRK` timestamp files.

```
import "github.com/flywave/go-rtk"
```

## Prerequisites

- **Go 1.13+** — no generics, no `errors.Is`/`errors.As`, no `os.ReadFile`
- **Sibling repos** — `go.mod` uses `replace` directives for `go-geoid` and `go-proj`. They must exist at `../go-geoid` and `../go-proj` relative to this repo root.
- **Prebuilt RTKLIB** — `libs/<platform_arch>/librtklib.a` is checked in per platform/arch. If missing, rebuild from `external/rtklib/` with CMake.
- **C toolchain** — gcc or clang. LDFLAGS are per-file in `gtime.go`, `pos.go`, etc.

### Supported platforms

| Platform | Arch | Lib directory |
|---|---|---|
| macOS | amd64 | `libs/darwin/` |
| macOS | arm64 | `libs/darwin_arm/` |
| Linux | amd64 | `libs/linux/` |
| Linux | arm64 | `libs/linux_arm/` |
| Windows | amd64 | `libs/windows/` |

## Build

```sh
go build ./...
go test ./...
go vet ./...
```

## Dependencies

| Dependency | Source | Purpose |
|---|---|---|
| RTKLIB (C) | `external/rtklib/` | GNSS processing engine |
| [go-proj](https://github.com/flywave/go-proj) | `../go-proj` | Coordinate transforms via PROJ |
| [go-geoid](https://github.com/flywave/go-geoid) | `../go-geoid` | EGM84/96/2008 geoid height conversions |
| goexif | — | EXIF tag parsing |
| go-xmp | — | XMP metadata parsing |
| gocsv | — | CSV marshaling for `.MRK` files |

## API reference

### GNSS time (`GTime`)

```go
// Constructors
gt := rtk.NewGTimeFromGPSTime(2024, 275295.301059)  // GPS week + seconds
gt := rtk.NewGTimeFromGalileoTime(2024, 100.0)       // Galileo week + seconds
gt := rtk.NewGTimeFromBDTime(2024, 100.0)            // BeiDou week + seconds
gt := rtk.NewGTimeFromEpoch([6]float64{2024, 6, 15, 10, 30, 45.5})  // calendar
gt := rtk.NewGTimeFromStr("2024/06/15 10:30:45.123") // parse string
gt := rtk.NewGPSTFromUTC(utc)        // GPS ↔ UTC
gt := rtk.NewUtcTime(str)            // UTC from string
gt := rtk.NewGPSTTime(str)           // GPS time from string
gt := rtk.NewUtcTimeFromTime(t)      // Go time.Time → UTC
gt := rtk.NewGPSTTimeFromTime(t)     // Go time.Time → GPS
gt := rtk.Current()                  // current system time

// Methods
week, sec := gt.GpsTime()            // → GPS week, seconds
week, sec := gt.GalileoTime()        // → Galileo week, seconds
week, sec := gt.BDTime()             // → BeiDou week, seconds
ep := gt.Epoch()                     // → [year, month, day, hour, min, sec]
utc := gt.UTC()                      // GPS → UTC conversion
diff := gt.Diff(other)               // seconds between two times
gt.Add(30.0)                         // add seconds (mutates)
s := gt.ToString(0)                  // formatted string
doy := gt.DayOfYear()                // day of year (1-366)
t := gt.Time()                       // raw time_t
sec := gt.Sec()                      // fractional seconds
```

### Position files (`.pos`)

```go
poses, trange := rtk.ReadPos("solution.pos")
// poses → []rtk.Pos with Gpst, Latitude, Longitude, Height, Q, Ns,
//         CovEE/CovNN/CovUU/CovEN/CovNU/CovUE, Age, Ratio
// trange → rtk.TimeRange — earliest and latest time in file
```

### Timestamp marks (`.MRK`)

```go
f, _ := os.Open("marks.MRK")
mrks, err := rtk.ReadMRK(f)
// mrks → []rtk.MRK with Sequence, Time, Week, PhaseCompNs/Ew/V,
//         Latitude, Longitude, Altitude, Std, State
gt := mrks[i].GetGpst()  // GPS time from Week + Time fields
```

#### Parse individual fields

```go
var pc rtk.PhaseComp
pc.UnmarshalCSV("  -26,N")            // → {v: -26, d: "N"}

var a rtk.Angle
a.UnmarshalCSV("-36.41144502,Lat")    // → {v: -36.41144502, d: "Lat"}

var w rtk.Week
w.UnmarshalCSV("[2024]")              // → {w: 2024}

var rs rtk.RtkState
rs.UnmarshalCSV("16,Q")               // → {e: 16, f: "Q"}

var s rtk.Std
s.UnmarshalCSV("1.492559, 1.368525, 3.180128")
```

### PPK processing

```go
err, sols := rtk.PPKSolution(
    "solution.pos",                    // RTKLIB .pos file
    "marks.MRK",                       // timestamp marks
    rtk.WGS84,                         // pose datum (WGS84 or CGCS2000)
    geoid.HAE,                         // pose vertical datum
    "+proj=longlat +ellps=WGS84 +datum=WGS84 +no_defs",  // output SRS
    geoid.HAE,                         // output vertical datum
    0.0,                               // ellipsoid offset
)
// sols → []rtk.PPKSol with Pos [3]float64 and Weight [3]float64
```

The PPK pipeline:
1. `ReadPos` loads the RTKLIB solution
2. `ReadMRK` loads the timestamp marks
3. `ppkInitPos` matches each mark to the nearest two solution positions
4. `ppkUpdatedMrks` interpolates, applies phase center offsets
5. `convertGPS` transforms coordinates to the target SRS and vertical datum

#### Standalone helpers

```go
deg := rtk.RadianToDegree(math.Pi)    // 180.0
rad := rtk.DegreeToRadian(180.0)      // π
```

### Coordinate transforms

```go
// Geoid height conversions
h := rtk.MSLToWGS84(100.0, 120.0, 30.0, geoid.EGM2008)
h := rtk.WGS84ToMSL(120.0, 30.0, 100.0, geoid.EGM96)
h := rtk.HAEToMSL(120.0, 30.0, 100.0, 0.0, geoid.EGM84)
h := rtk.MSLToHAE(100.0, 120.0, 30.0, 0.0, geoid.EGM2008)

// Datum constants
_ = rtk.WGS84    // 0
_ = rtk.CGCS2000 // 1
```

### Solution processing (RINEX → `.pos`)

```go
err := rtk.Solution("rover.obs", "base.obs", "brdc.nav", "output.pos")
// Uses PPP-Static mode with GPS+GLONASS+Galileo+BeiDou,
// precise ephemeris, ionosphere-free LC, tropo estimation
```

### Raw binary → RINEX conversion

```go
gt := rtk.NewGTimeFromGPSTime(2024, 275295.0)
rtk.RawToRIndex(*gt, rtk.FormatRTCM3,
    "raw.bin",    // input binary
    "out.obs",    // output observation file
    "out.nav",    // output navigation file
    "out.geo",    // optional geostationary file ("" to skip)
)
```

Supported input formats via the `Format` type:

| Constant | Value | Input format |
|---|---|---|
| `FormatRTCM2` | 0 | RTCM 2 |
| `FormatRTCM3` | 1 | RTCM 3 |
| `FormatOEM4` | 2 | NovAtel OEMV/4 |
| `FormatOEM3` | 3 | NovAtel OEM3 |
| `FormatUBX` | 4 | u-blox LEA-*T |
| `FormatSS2` | 5 | NovAtel Superstar II |
| `FormatCRES` | 6 | Hemisphere |
| `FormatSTQ` | 7 | SkyTraq S1315F |
| `FormatJAVAD` | 8 | JAVAD GRIL/GREIS |
| `FormatNVS` | 9 | NVS NVC08C |
| `FormatBINEX` | 10 | BINEX |
| `FormatRT17` | 11 | Trimble RT17 |
| `FormatSEPT` | 12 | Septentrio |
| `FormatRINEX` | 13 | RINEX |
| `FormatSP3` | 14 | SP3 |
| `FormatRNXCLK` | 15 | RINEX CLK |
| `FormatSBAS` | 16 | SBAS messages |
| `FormatNMEA` | 17 | NMEA 0183 |

### EXIF / XMP metadata

```go
f, _ := os.Open("DJI_0064.JPG")
err, fields := rtk.ReadExifXMP(f)
// fields["drone-dji: RtkFlag"]     = "50"  (fixed=50, float=60)
// fields["drone-dji: RtkStdLon"]   = "0.005"
// fields["drone-dji: RtkStdLat"]   = "0.004"
// fields["drone-dji: RtkStdAlt"]   = "0.008"
// fields["exif: GPSLatitude"]      = raw EXIF GPS tag
// fields["exif: GPSLongitude"]     = raw EXIF GPS tag
```

### Time utilities

```go
// Parse timestamp string → epoch + fractional seconds
ep, sec := rtk.ParseUtcTime("2024/06/15 10:30:45.123")
// ep = [2024, 6, 15, 10, 30, 45], sec = 0.123

// Go time.Time → Tm (struct tm wrapper)
tm := rtk.NewTm(time.Now())

// strftime formatting
s := rtk.Strftime("%Y-%m-%d %H:%M:%S", tm)
// s = "2024-06-15 10:30:45"

// Local time helpers
gt := rtk.NewUtcTimeFromLocal("2024-06-15 10:30:45", "Asia/Shanghai")
gt := rtk.NewGPSTTimeFromLocal("2024-06-15 10:30:45", "America/New_York")
gt := rtk.NewUTCFromLocalTM(tm, 0.0)
gt := rtk.NewGPSTTimeFromLocalTM(tm, 0.0)
```

## File formats

| Format | Reader | Description |
|---|---|---|
| `.pos` | `ReadPos` | RTKLIB solution file (header lines start with `%`; tab/space-separated) |
| `.MRK` | `ReadMRK` | Tab-separated timestamp marks with GPS week/SOW, phase center offsets, and RTK state |
| `.exif` / XMP | `ReadExifXMP` | DJI drone image metadata (RtkFlag, RtkStd, absolute altitude) |
| `.obs` / `.nav` | `Solution` | RINEX observation and navigation files |
| `.bin` | `RawToRIndex` | Raw GNSS binary → RINEX conversion (RTCM3, UBX, OEM4, etc.) |

### `.MRK` column layout

Tab-delimited, no header row:

```
Sequence  GPSSecondOfWeek  GPSWeek  NorthOff  EastOff  VelOff  Latitude  Longitude  EllipsoidHeight  Std  RtkState
```

### `.pos` column layout

Standard RTKLIB `.pos` format — header lines starting with `%`, data rows with time, position, statistics.

## Data directories

- `proj_data/` (9 MB) — PROJ datum grids and transform data
- `geoid_data/` (20 MB) — EGM84/96/2008 geoid height grids (PGM format)
- `testdata/` — sample `.pos`, `.MRK`, `.exif`, `.bin`, and RINEX files
- `libs/` — prebuilt `librtklib.a` per platform/arch + `rtklib.h`

Paths are set automatically via `init()` in `datum.go` and `exif.go` — no environment variables needed.

## Architecture

12 `.go` files, single flat package:

| File | Responsibility |
|---|---|
| `datum.go` | Coordinate transforms, datum types, geoid conversions |
| `gtime.go` | GNSS time types (GPS, Galileo, BeiDou, UTC) — CGo wrappers |
| `pos.go` | `.pos` file reader via RTKLIB `readsol` |
| `ppk.go` | PPK pipeline: time matching, interpolation, coordinate transforms |
| `ppkraw.go` | Raw GNSS binary → RINEX conversion (18 format types) |
| `mark.go` | `.MRK` file reader, CSV field types |
| `solution.go` | RTKLIB solution init/processing (PPP-Static) |
| `exif.go` | EXIF/XMP metadata parser for DJI drone images |
| `tm.go` | Text timestamp parser via embedded C |

## Notes

- Go 1.13 compatible
- CGo with per-file `#cgo` directives — no global CFLAGS/LDFLAGS
- Tests use `t.FailNow()` on failure (no error messages)
- `go.sum` has no entries for `replace`-only deps — CI may need `-mod=mod`
