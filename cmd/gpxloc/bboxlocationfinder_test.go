package main

import (
	"reflect"
	"testing"

	"github.com/paulmach/orb"
)

func TestNewBboxLocationFinder(t *testing.T) {
	//bbox='5.630665,45.031614,5.634817,45.034214'
	p1 := orb.Point{5.630665, 45.031614}
	p2 := orb.Point{5.634817, 45.034214}
	okBound := orb.MultiPoint{p1, p2}.Bound()
	koBound := orb.Bound{}
	okLf := BboxLocationFinder{bbox: okBound}
	koLf := BboxLocationFinder{bbox: koBound}

	type args struct {
		bbox string
	}
	tests := []struct {
		name    string
		args    args
		want    BboxLocationFinder
		wantErr bool
	}{
		{"OK", args{"5.630665,45.031614,5.634817,45.034214"}, okLf, false},
		{"KO", args{"999,999,999,999"}, koLf, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBboxLocationFinder(tt.args.bbox)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBboxLocationFinder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBboxLocationFinder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBboxLocationFinder_Match(t *testing.T) {
	lf, _ := NewBboxLocationFinder("5.686787,45.085960,5.689791,45.087786")

	type args struct {
		lat float64
		lon float64
	}
	tests := []struct {
		name string
		lf   BboxLocationFinder
		args args
		want bool
	}{
		{"OK within", lf, args{45.0869, 5.6880}, true},
		{"KO outside", lf, args{45.0412, 5.6959}, false},
		{"KO empty bound", BboxLocationFinder{}, args{45.0412, 5.6959}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lf.Match(tt.args.lat, tt.args.lon); got != tt.want {
				t.Errorf("BboxLocationFinder.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
