package config_test

import (
	"testing"

	"github.com/yoctoMNS/go-rpg/internal/config"
)

func TestDefault_WindowSize(t *testing.T) {
	t.Parallel()

	cfg := config.Default()

	gotW, gotH := cfg.WindowSize()

	wantW := cfg.ScreenWidth * cfg.WindowScale
	wantH := cfg.ScreenHeight * cfg.WindowScale

	if gotW != wantW || gotH != wantH {
		t.Errorf("WindowSize() = (%d, %d), want (%d, %d)", gotW, gotH, wantW, wantH)
	}
}

func TestDefault_ScreenFitsTileGrid(t *testing.T) {
	t.Parallel()

	cfg := config.Default()

	if cfg.ScreenWidth%cfg.TileSize != 0 {
		t.Errorf("ScreenWidth %d is not a multiple of TileSize %d", cfg.ScreenWidth, cfg.TileSize)
	}
	if cfg.ScreenHeight%cfg.TileSize != 0 {
		t.Errorf("ScreenHeight %d is not a multiple of TileSize %d", cfg.ScreenHeight, cfg.TileSize)
	}
}
