// Package geom provides basic 2D geometry types shared across the game.
package geom

import "math"

// Vec2 is a 2D vector or point. All methods return a new value; Vec2 itself is immutable.
type Vec2 struct {
	X, Y float64
}

// Add returns v + other.
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub returns v - other.
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{X: v.X - other.X, Y: v.Y - other.Y}
}

// Scale returns v scaled by s.
func (v Vec2) Scale(s float64) Vec2 {
	return Vec2{X: v.X * s, Y: v.Y * s}
}

// Len returns the length of v.
func (v Vec2) Len() float64 {
	return math.Hypot(v.X, v.Y)
}

// Normalize returns v scaled to length 1, or the zero vector if v is the zero vector.
func (v Vec2) Normalize() Vec2 {
	l := v.Len()
	if l == 0 {
		return Vec2{}
	}
	return v.Scale(1 / l)
}
