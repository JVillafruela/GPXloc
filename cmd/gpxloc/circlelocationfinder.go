package main

import (
	"errors"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

type CircleLocationFinder struct {
	center orb.Point
	// radius in meters
	radius int64
	bbox   orb.Bound
}

// CircleLocationFinder constructor
func NewCircleLocationFinder(lat, lon float64, radius int64) (CircleLocationFinder, error) {
	err := validateLatitude(lat)
	if err != nil {
		return CircleLocationFinder{}, err
	}
	err = validateLongitude(lon)
	if err != nil {
		return CircleLocationFinder{}, err
	}
	if radius <= 0 {
		return CircleLocationFinder{}, errors.New("invalid radius")
	}
	point := orb.Point{lon, lat}
	bbox := geo.NewBoundAroundPoint(point, float64(radius))

	return CircleLocationFinder{point, radius, bbox}, nil
}

func (lf CircleLocationFinder) PossibleMatch(bbox orb.Bound) bool {
	return BoundsIntersect(lf.bbox, bbox)
}

func (lf CircleLocationFinder) Match(lat, lon float64) bool {
	point := orb.Point{lon, lat}
	return geo.Distance(lf.center, point) <= float64(lf.radius)
}

func (lf CircleLocationFinder) Bounds() orb.Bound {
	return lf.bbox
}
