package geom

// Rect is an axis-aligned rectangle with its origin at the top-left corner.
type Rect struct {
	X, Y, W, H float64
}

// Intersects reports whether r and other overlap.
// Rectangles that only touch at an edge or corner are not considered overlapping.
func (r Rect) Intersects(other Rect) bool {
	return r.X < other.X+other.W &&
		other.X < r.X+r.W &&
		r.Y < other.Y+other.H &&
		other.Y < r.Y+r.H
}

// Contains reports whether p lines within r.
func (r Rect) Contains(p Vec2) bool {
	return p.X >= r.X &&
		p.X < r.X+r.W &&
		p.Y >= r.Y &&
		p.Y < r.Y+r.H
}

// Move returns r translated by d.
func (r Rect) Move(d Vec2) Rect {
	return Rect{
		X: r.X + d.X,
		Y: r.Y + d.Y,
		W: r.W,
		H: r.H,
	}
}
