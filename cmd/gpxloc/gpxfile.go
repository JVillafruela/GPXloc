package main

import (
	"github.com/paulmach/orb"
	"github.com/tkrajina/gpxgo/gpx"
)

type GPXfile struct {
	gpx    *gpx.GPX
	bounds orb.Bound
}

// get the coordinates from the exif data
func ReadGPXfile(fname string) (gpxfile GPXfile, err error) {

	gpxfile.gpx, err = gpx.ParseFile(fname)
	if err != nil {
		return
	}

	bounds := gpxfile.gpx.Bounds()
	ptMin := orb.Point{bounds.MinLongitude, bounds.MinLatitude}
	ptMax := orb.Point{bounds.MaxLongitude, bounds.MaxLatitude}
	gpxfile.bounds = orb.Bound{Min: ptMin, Max: ptMax}

	return
}
