package main

import (
	"errors"
	"fmt"

	"github.com/paulmach/orb"
)

type LocationFinder interface {
	Match(lat, lon float64) bool
	PossibleMatch(bbox orb.Bound) bool
	Bounds() orb.Bound
}

// get a Bound object from bounding box coordinates
// bbox format : lon1,lat1,lon2,lat2 eg. 5.63066,45.03161,5.63481,45.03421
func BboxBound(str string) (orb.Bound, error) {

	var minLon, minLat, maxLat, maxLon float64

	n, err := fmt.Sscanf(str, "%f,%f,%f,%f", &minLon, &minLat, &maxLon, &maxLat)
	if n != 4 {
		return orb.Bound{}, err
	}
	err = validateLatitude(minLat)
	if err != nil {
		return orb.Bound{}, err
	}
	err = validateLatitude(maxLat)
	if err != nil {
		return orb.Bound{}, err
	}
	err = validateLongitude(minLon)
	if err != nil {
		return orb.Bound{}, err
	}
	err = validateLongitude(maxLon)
	if err != nil {
		return orb.Bound{}, err
	}

	p1 := orb.Point{minLon, minLat}
	p2 := orb.Point{maxLon, maxLat}

	return BoundingBox(p1, p2), nil
}

func BoundingBox(p0, p1 orb.Point) orb.Bound {
	return orb.MultiPoint{p0, p1}.Bound()
}

func BoundsIntersect(bbox1 orb.Bound, bbox2 orb.Bound) bool {
	return bbox1.Intersects(bbox2) // || bbox2.Intersects(bbox1)
}

// validate a latitude in WGS84 system
func validateLatitude(lat float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("invalid latitude (WGS84 [-90,+90])")
	}
	return nil
}

// validate a longitude in WGS84 system
func validateLongitude(lon float64) error {
	if lon < -180 || lon > 180 {
		return errors.New("invalid longitude (WGS84 [-180,+180])")
	}
	return nil
}
