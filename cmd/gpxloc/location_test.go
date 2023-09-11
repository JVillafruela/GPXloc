package main

import (
	"reflect"
	"testing"

	"github.com/paulmach/orb"
)

func Test_bboxBound(t *testing.T) {
	//bbox='5.630665,45.031614,5.634817,45.034214'
	p1 := orb.Point{5.630665, 45.031614}
	p2 := orb.Point{5.634817, 45.034214}
	okBound := orb.MultiPoint{p1, p2}.Bound()
	errBound := orb.Bound{}

	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    orb.Bound
		wantErr bool
	}{
		{"OK min,max", args{"5.630665,45.031614,5.634817,45.034214"}, okBound, false},
		{"OK max,min", args{"5.634817,45.034214,5.630665,45.031614"}, okBound, false},
		{"OK with leading spaces", args{"5.630665, 45.031614, 5.634817, 45.034214"}, okBound, false},
		{"KO with trailing spaces", args{"5.630665 ,45.031614,5.634817,45.034214"}, errBound, true},
		{"KO non numeric", args{"5.630665,45.031614,5.634817,fortytwo"}, errBound, true},
		{"KO missing coordinate", args{"5.630665,45.031614,5.634817"}, errBound, true},
		{"KO invalid min lat", args{"999, 45.031614, 5.634817, 45.034214"}, errBound, true},
		{"KO invalid max lat", args{"5.630665, 45.031614, 999, 45.034214"}, errBound, true},
		{"KO invalid min lon", args{"5.630665, 999, 5.634817, 45.034214"}, errBound, true},
		{"KO invalid max lon", args{"5.630665, 45.031614, 5.634817, 999"}, errBound, true},
		// TODO KO p2 = p1
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BboxBound(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("bboxBound() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bboxBound() = %v, want %v", got, tt.want)
			}
		})
	}
}
