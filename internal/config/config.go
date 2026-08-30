package config

// Config represents the game's startup configuration.
type Config struct {
	Title        string // window title
	ScreenWidth  int    // logical resolution width
	ScreenHeight int    // logical resolution height
	WindowScale  int    // initial window sdcale relative to the logical resolution
	TileSize     int    // pixel size of one tile
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

// WindowSize returns the actual pixel size of the inital window.
func (c Config) WindowSize() (width, height int) {
	return c.ScreenWidth * c.WindowScale, c.ScreenHeight * c.WindowScale
}
