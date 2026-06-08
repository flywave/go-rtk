module github.com/flywave/go-rtk

go 1.13

require (
	github.com/flywave/go-geoid v0.0.0-00010101000000-000000000000
	github.com/flywave/go-proj v0.0.0-00010101000000-000000000000
	github.com/gocarina/gocsv v0.0.0-20210516172204-ca9e8a8ddea8
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
	trimmer.io/go-xmp v0.0.0-20200923092433-f9b6ca6c4a87
)

replace github.com/flywave/go-geoid => ../go-geoid

replace github.com/flywave/go-proj => ../go-proj
