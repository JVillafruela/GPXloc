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
	okBound := BoundingBox(p1, p2)
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

func TestBoundsIntersect(t *testing.T) {
	p0 := orb.Point{0.0, 0.0}
	p1 := orb.Point{1.0, 1.0}
	p2 := orb.Point{2.0, 2.0}
	p3 := orb.Point{3.0, 3.0}
	bbox_0_1 := BoundingBox(p0, p1)
	bbox_0_2 := BoundingBox(p0, p2)
	bbox_1_3 := BoundingBox(p1, p3)
	bbox_2_3 := BoundingBox(p2, p3)
	type args struct {
		bbox1 orb.Bound
		bbox2 orb.Bound
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"OK bbox_0_1,bbox_0_1", args{bbox_0_1, bbox_0_1}, true},

		{"OK bbox_0_1,bbox_0_2", args{bbox_0_1, bbox_0_2}, true},
		{"OK bbox_0_2,bbox_0_1", args{bbox_0_2, bbox_0_1}, true},

		{"OK bbox_0_2,bbox_1_3", args{bbox_0_2, bbox_1_3}, true},
		{"OK bbox_1_3,bbox_0_2", args{bbox_1_3, bbox_0_2}, true},

		{"KO bbox_0_1,bbox_2_3", args{bbox_0_1, bbox_2_3}, false},
		{"KO bbox_2_3,bbox_0_1", args{bbox_2_3, bbox_0_1}, false},

		{"OK bbox_0_1,bbox_1_3", args{bbox_0_1, bbox_1_3}, true}, // one point intersection
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoundsIntersect(tt.args.bbox1, tt.args.bbox2); got != tt.want {
				t.Errorf("BoundsIntersect() = %v, want %v", got, tt.want)
			}
		})
	}
}
