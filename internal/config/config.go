// Package config provides the startup configuration shared across the game.
//
// The goal is to keep magic numbers out of the rest of the codebase.
// Collecting them here means that switching to a config file (JSON/TOML)
// later only requires changing how Config is constructed.
package config

// Config represents the game's startup configuration.
//
// It must not hold values that change at runtime (e.g. player HP).
// Only values that are fixed at startup and never change afterward belong here.
type Config struct {
	// Title is the window title.
	Title string

	// ScreenWidth / ScreenHeight is the logical resolution (the game's
	// internal coordinate system). It is independent of the window size;
	// Ebitengine scales it up for display. Pixel-art RPGs conventionally
	// fix a low logical resolution and scale it by an integer factor.
	ScreenWidth  int
	ScreenHeight int

	// WindowScale is the initial window scale relative to the logical resolution.
	WindowScale int

	// TileSize is the pixel size of one tile, the base unit for maps and collision.
	TileSize int
}

// Default returns the development-time default values.
//
// 320x240 (QVGA) resembles the look of NES/SNES-era RPGs, and with 16px
// tiles it fits exactly 20x15 tiles — a convenient size to work with.
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
