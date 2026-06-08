# go-rtk — agent instructions

## What this is

Flat Go package (`package rtk`) wrapping [RTKLIB](https://github.com/tomojitakasu/RTKLIB) via CGo for PPK/RTK GNSS post-processing. Works with drone image metadata (EXIF/XMP), RINEX observation files, and `.pos`/`.MRK` files.

## Build prerequisites

- **Sibling repos** — `go.mod` uses local `replace` directives for `go-geoid` and `go-proj`. They must exist at `../go-geoid` and `../go-proj` relative to the repo root.
- **Prebuilt RTKLIB** — `lib/librtklib.a` is checked in. If missing, build it from `external/rtklib/` with CMake.
- **CGo toolchain** — requires gcc/clang. Platform-specific LDFLAGS are in `gtime.go`, `pos.go`, `solution.go`, etc.

## Commands

```sh
go build ./...
go test ./...
go vet ./...
```

No Makefile or taskfile — plain Go toolchain.

## Architecture (12 `.go` files, single package)

| File | Responsibility |
|---|---|
| `datum.go` | Coordinate transforms (WGS84/CGCS2000), geoid height conversions (EGM84/96/2008) |
| `gtime.go` | GNSS time types (GPS, Galileo, BeiDou, UTC) — CGo wrappers around RTKLIB `gtime_t` |
| `pos.go` | Position type + reading `.pos` files via RTKLIB |
| `ppk.go` | PPK processing: nearest-time matching, coordinate transforms, CSV output |
| `ppkraw.go` | RTCM3 → RINEX conversion via CGo |
| `mark.go` | `.MRK` file reader (CSV-based), mark-to-pos interpolation |
| `solution.go` | RTKLIB solution init/processing (PPP, RINEX reading) |
| `exif.go` | EXIF/XMP metadata parser (DJI drone images, typically) |
| `tm.go` | Text timestamp parser via embedded C |

`init()` functions in `datum.go` and `exif.go` set up local data paths and EXIF parsers (no env vars to configure).

## Data directories

- `proj_data/` (9 MB) — PROJ datum grids and transform data
- `geoid_data/` (20 MB) — EGM84/96/2008 geoid height grids (PGM format)
- `lib/` — prebuilt `librtklib.a` + `rtklib.h`

## Tests

Run `go test -v ./...`. Tests use files in `testdata/`:
- `testdata/test.pos` — solution file for `TestPosRead`
- `testdata/100_0004_Timestamp.MRK` — mark file for `TestMarkRead`
- `testdata/100_0004_PPKRAW.bin`, `100_0004_Rinex.obs`, `dji_rtk.exif`, etc.

## Style

- Go 1.13 (no generics, no `errors.Is`/`errors.As`, no `os.ReadFile`)
- CGo with `#cgo` directives in each file that needs them
- Test failures use `t.FailNow()` (no error messages)
- No `go.sum` entries for replace-only deps — CI may need `-mod=mod`
