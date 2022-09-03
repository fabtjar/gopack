package pack

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCanFit(t *testing.T) {
	var tests = []struct {
		r1, r2 Rectangle
		want   bool
	}{
		{Rectangle{0, 0, 200, 100}, Rectangle{0, 0, 200, 100}, true},
		{Rectangle{0, 0, 0, 0}, Rectangle{0, 0, 0, 0}, true},
		{Rectangle{0, 0, 200, 100}, Rectangle{0, 0, 400, 200}, false},
		{Rectangle{0, 0, 200, 100}, Rectangle{0, 0, 150, 50}, true},
		{Rectangle{0, 0, 200, 100}, Rectangle{0, 0, 9999, 1}, false},
		{Rectangle{0, 0, 200, 100}, Rectangle{0, 0, 1, 9999}, false},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%v,%v", tt.r1, tt.r2)
		t.Run(testName, func(t *testing.T) {
			ans := tt.r1.canFit(tt.r2)
			if ans != tt.want {
				t.Errorf("r1.canFit(r2) = %v; want %v", ans, tt.want)
			}
		})
	}
}

func TestOverlaps(t *testing.T) {
	var tests = []struct {
		r1, r2 Rectangle
		want   bool
	}{
		{Rectangle{0, 0, 200, 150}, Rectangle{0, 0, 200, 150}, true},
		{Rectangle{0, 0, 200, 150}, Rectangle{50, 50, 300, 200}, true},
		{Rectangle{0, 0, 500, 400}, Rectangle{450, 350, 50, 50}, true},
		{Rectangle{0, 0, 200, 150}, Rectangle{300, 200, 200, 150}, false},
		{Rectangle{0, 0, 200, 150}, Rectangle{-300, -200, 200, 150}, false},
		{Rectangle{0, 0, 200, 150}, Rectangle{200, 0, 200, 150}, false},
		{Rectangle{0, 0, 200, 150}, Rectangle{-200, 0, 200, 150}, false},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("%v,%v", tt.r1, tt.r2)
		t.Run(testName, func(t *testing.T) {
			ans := tt.r1.overlaps(tt.r2)
			if ans != tt.want {
				t.Errorf("r1.overlaps(r2) = %v; want %v", ans, tt.want)
			}
		})
	}
}

func TestCutInside(t *testing.T) {
	r1 := Rectangle{0, 0, 300, 200}
	r2 := Rectangle{10, 20, 30, 50}

	rects := r1.cut(r2)

	want := []Rectangle{
		{0, 0, 300, 20},
		{0, 70, 300, 130},
		{0, 0, 10, 200},
		{40, 0, 260, 200},
	}

	if !reflect.DeepEqual(rects, want) {
		t.Errorf("r1.cut(r2) = %v; want %v", rects, want)
	}
}

func TestCutTouchingLeftEdge(t *testing.T) {
	r1 := Rectangle{0, 0, 300, 200}
	r2 := Rectangle{0, 20, 30, 50}

	rects := r1.cut(r2)

	want := []Rectangle{
		{0, 0, 300, 20},
		{0, 70, 300, 130},
		{30, 0, 270, 200},
	}

	if !reflect.DeepEqual(rects, want) {
		t.Errorf("r1.cut(r2) = %v; want %v", rects, want)
	}
}

func TestCutOverlappingRightEdge(t *testing.T) {
	r1 := Rectangle{0, 0, 300, 200}
	r2 := Rectangle{290, 20, 30, 50}

	rects := r1.cut(r2)

	want := []Rectangle{
		{0, 0, 300, 20},
		{0, 70, 300, 130},
		{0, 0, 290, 200},
	}

	if !reflect.DeepEqual(rects, want) {
		t.Errorf("r1.cut(r2) = %v; want %v", rects, want)
	}
}

func TestCutOutsideBottomRight(t *testing.T) {
	r1 := Rectangle{0, 0, 300, 200}
	r2 := Rectangle{500, 400, 10, 20}

	rects := r1.cut(r2)

	if len(rects) > 0 {
		t.Errorf("r1.cut(r2) = %v; want []", rects)
	}
}

func TestCutOutsideTopLeft(t *testing.T) {
	r1 := Rectangle{0, 0, 300, 200}
	r2 := Rectangle{-50, -40, 10, 20}

	rects := r1.cut(r2)

	if len(rects) > 0 {
		t.Errorf("r1.cut(r2) = %v; want []", rects)
	}
}
