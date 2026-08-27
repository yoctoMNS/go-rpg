// Package config provides the startup configuration shared across the game.
package config

// Config represents the game's startup configuration.
type Config struct {
	// Title is the window title.
	Title string

	// ScreenWidth and ScreenHeight are the logical resolution.
	ScreenWidth  int
	ScreenHeight int

	// WindowScale is the initial window scale relative to the logical resolution.
	WindowScale int

	// TileSize is the pixel size of one tile.
	TileSize int
}

// Default returns the development-time default values.
func Default() Config {
	return Config{
		Title:        "Go RPG",
		ScreenWidth:  320,
		ScreenHeight: 240,
		WindowScale:  3,
		TileSize:     16,
	}
}

// WindowSize returns the actual pixel size of the initial window.
func (c Config) WindowSize() (width, height int) {
	return c.ScreenWidth * c.WindowScale, c.ScreenHeight * c.WindowScale
}
