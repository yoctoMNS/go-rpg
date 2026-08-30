package geom_test

import (
	"math"
	"testing"

	"github.com/yoctoMNS/go-rpg/internal/geom"
)

func TestVec2_Add(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		a    geom.Vec2
		b    geom.Vec2
		want geom.Vec2
	}{
		{
			name: "positive",
			a:    geom.Vec2{X: 1, Y: 2},
			b:    geom.Vec2{X: 3, Y: 4},
			want: geom.Vec2{X: 4, Y: 6},
		},
		{
			name: "negative",
			a:    geom.Vec2{X: -1, Y: 5},
			b:    geom.Vec2{X: 2, Y: -3},
			want: geom.Vec2{X: 1, Y: 2},
		},
		{
			name: "zero",
			a:    geom.Vec2{},
			b:    geom.Vec2{X: 1, Y: 1},
			want: geom.Vec2{X: 1, Y: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.a.Add(tt.b); got != tt.want {
				t.Errorf("Add(%v) = %v, want %v", tt.b, got, tt.want)
			}
		})
	}
}

func TestVec2_Sub(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		a    geom.Vec2
		b    geom.Vec2
		want geom.Vec2
	}{
		{
			name: "positive",
			a:    geom.Vec2{X: 5, Y: 5},
			b:    geom.Vec2{X: 2, Y: 3},
			want: geom.Vec2{X: 3, Y: 2},
		},
		{
			name: "result negative",
			a:    geom.Vec2{X: 1, Y: 1},
			b:    geom.Vec2{X: 3, Y: 3},
			want: geom.Vec2{X: -2, Y: -2},
		},
		{
			name: "same vector",
			a:    geom.Vec2{X: 4, Y: 4},
			b:    geom.Vec2{X: 4, Y: 4},
			want: geom.Vec2{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.a.Sub(tt.b); got != tt.want {
				t.Errorf("Sub(%v) = %v, want %v", tt.b, got, tt.want)
			}
		})
	}
}

func TestVec2_Scale(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		v    geom.Vec2
		s    float64
		want geom.Vec2
	}{
		{
			name: "double",
			v:    geom.Vec2{X: 1, Y: 2},
			s:    2,
			want: geom.Vec2{X: 2, Y: 4},
		},
		{
			name: "zero scale",
			v:    geom.Vec2{X: 5, Y: 5},
			s:    0,
			want: geom.Vec2{},
		},
		{
			name: "negative scale",
			v:    geom.Vec2{X: 1, Y: -1},
			s:    -1,
			want: geom.Vec2{X: -1, Y: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.v.Scale(tt.s); got != tt.want {
				t.Errorf("Scale(%v) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestVec2_Len(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		v    geom.Vec2
		want float64
	}{
		{
			name: "3-4-5 triangle",
			v:    geom.Vec2{X: 3, Y: 4},
			want: 5,
		},
		{
			name: "zero vector",
			v:    geom.Vec2{},
			want: 0,
		},
		{
			name: "negative components",
			v:    geom.Vec2{X: -3, Y: -4},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.v.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVec2_Normalize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		v    geom.Vec2
		want geom.Vec2
	}{
		{
			name: "unit x",
			v:    geom.Vec2{X: 5, Y: 0},
			want: geom.Vec2{X: 1, Y: 0},
		},
		{
			name: "3-4-5 triangle",
			v:    geom.Vec2{X: 3, Y: 4},
			want: geom.Vec2{X: 0.6, Y: 0.8},
		},
		{
			name: "zero vector stays zero",
			v:    geom.Vec2{},
			want: geom.Vec2{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.v.Normalize()
			if math.Abs(got.X-tt.want.X) > 1e-9 || math.Abs(got.Y-tt.want.Y) > 1e-9 {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}
