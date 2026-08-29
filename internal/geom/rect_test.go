package geom_test

import (
	"testing"

	"github.com/yoctoMNS/go-rpg/internal/geom"
)

func TestRect_Intersects(t *testing.T) {
	t.Parallel()

	base := geom.Rect{X: 0, Y: 0, W: 10, H: 10}

	tests := []struct {
		name  string
		other geom.Rect
		want  bool
	}{
		{name: "fully overlapping", other: geom.Rect{X: 2, Y: 2, W: 4, H: 4}, want: true},
		{name: "overlaps by 1px", other: geom.Rect{X: 9, Y: 0, W: 10, H: 10}, want: true},
		{name: "edges touch exactly", other: geom.Rect{X: 10, Y: 0, W: 10, H: 10}, want: false},
		{name: "corners touch exactly", other: geom.Rect{X: 10, Y: 10, W: 10, H: 10}, want: false},
		{name: "completely separate", other: geom.Rect{X: 100, Y: 100, W: 10, H: 10}, want: false},
		{name: "other contains base", other: geom.Rect{X: -5, Y: -5, W: 20, H: 20}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := base.Intersects(tt.other); got != tt.want {
				t.Errorf("Intersects(%v) = %v, want %v", tt.other, got, tt.want)
			}
		})
	}
}

func TestRect_Contains(t *testing.T) {
	t.Parallel()

	r := geom.Rect{X: 0, Y: 0, W: 10, H: 10}

	tests := []struct {
		name string
		p    geom.Vec2
		want bool
	}{
		{name: "inside", p: geom.Vec2{X: 5, Y: 5}, want: true},
		{name: "top-left corner is inclusive", p: geom.Vec2{X: 0, Y: 0}, want: true},
		{name: "bottom-right corner is exclusive", p: geom.Vec2{X: 10, Y: 10}, want: false},
		{name: "outside", p: geom.Vec2{X: 20, Y: 20}, want: false},
		{name: "negative coordinates", p: geom.Vec2{X: -1, Y: 5}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := r.Contains(tt.p); got != tt.want {
				t.Errorf("Contains(%v) = %v, want %v", tt.p, got, tt.want)
			}
		})
	}
}

func TestRect_Move(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		r    geom.Rect
		d    geom.Vec2
		want geom.Rect
	}{
		{name: "positive delta", r: geom.Rect{X: 0, Y: 0, W: 10, H: 5}, d: geom.Vec2{X: 3, Y: 4}, want: geom.Rect{X: 3, Y: 4, W: 10, H: 5}},
		{name: "negative delta", r: geom.Rect{X: 10, Y: 10, W: 4, H: 4}, d: geom.Vec2{X: -5, Y: -5}, want: geom.Rect{X: 5, Y: 5, W: 4, H: 4}},
		{name: "zero delta", r: geom.Rect{X: 1, Y: 2, W: 3, H: 4}, d: geom.Vec2{}, want: geom.Rect{X: 1, Y: 2, W: 3, H: 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.r.Move(tt.d); got != tt.want {
				t.Errorf("Move(%v) = %v, want %v", tt.d, got, tt.want)
			}
		})
	}
}
