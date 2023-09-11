package main

import (
	"reflect"
	"testing"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

func TestNewCircleLocationFinder(t *testing.T) {
	point := orb.Point{5.688, 45.087}
	radius := int64(20)
	bbox := geo.NewBoundAroundPoint(point, float64(radius))
	okLf := CircleLocationFinder{point, radius, bbox}
	koLf := CircleLocationFinder{}
	type args struct {
		lat    float64
		lon    float64
		radius int64
	}
	tests := []struct {
		name    string
		args    args
		want    CircleLocationFinder
		wantErr bool
	}{
		{"OK", args{45.087, 5.688, 20}, okLf, false},
		{"KO lat", args{999, 5.688, 20}, koLf, true},
		{"KO lon", args{45.087, 999, 20}, koLf, true},
		{"KO radius zero", args{45.087, 999, 0}, koLf, true},
		{"KO radius", args{45.087, 999, -42}, koLf, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCircleLocationFinder(tt.args.lat, tt.args.lon, tt.args.radius)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCircleLocationFinder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCircleLocationFinder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCircleLocationFinder_Match(t *testing.T) {

	lf, _ := NewCircleLocationFinder(45.087, 5.688, 20)
	type args struct {
		lat float64
		lon float64
	}
	tests := []struct {
		name string
		lf   CircleLocationFinder
		args args
		want bool
	}{
		{"OK", lf, args{45.087, 5.688}, true},
		{"KO", lf, args{45.090, 5.6815}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lf.Match(tt.args.lat, tt.args.lon); got != tt.want {
				t.Errorf("CircleLocationFinder.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
