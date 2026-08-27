// Package game provides the core Ebitengine game loop.
//
// Ebitengine exposes the Game Loop pattern via the ebiten.Game interface
// (Update / Draw / Layout). This package holds that implementation; later
// phases will delegate from here down to individual Scenes.
package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/yoctoMNS/go-rpg/internal/config"
)

// backgroundColor is the background color used while there is nothing else
// to draw yet. Plain black makes it hard to tell whether rendering is
// running or has crashed, so a dark navy is used instead (favors easy
// visual verification).
var backgroundColor = color.RGBA{R: 0x10, G: 0x18, B: 0x30, A: 0xff}

// Game implements ebiten.Game.
//
// Ebitengine calls Update then Draw, in that order, every frame:
//   - Update: state changes only. Must not draw.
//   - Draw:   drawing only. Must not change state.
//
// Keeping update and draw from mixing is a discipline that makes later
// refactoring (splitting into scenes, moving to ECS, adding tests) far easier.
type Game struct {
	cfg config.Config

	// ticks counts how many times Update has been called.
	// Counting logical frames instead of wall-clock time is the standard
	// approach in games: real-time counting makes behavior vary by
	// environment and destroys reproducibility.
	ticks uint64
}

// New creates a Game.
func New(cfg config.Config) *Game {
	return &Game{cfg: cfg}
}

// Update advances one frame's worth of state. Called 60 times per second by default.
func (g *Game) Update() error {
	g.ticks++
	return nil
}

// Draw renders one frame's worth of content to screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	// Verification aid for Day 01; treated as temporary and expected to be
	// removed as the implementation progresses.
	msg := fmt.Sprintf("Go RPG\nticks: %d\nTPS: %.1f  FPS: %.1f",
		g.ticks, ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

// Layout returns the game's logical resolution.
//
// Regardless of the return value, Ebitengine scales it up to the actual
// window size — so game code only ever has to reason about the logical
// resolution's coordinate system.
func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

// Run creates the window and starts the game loop.
// Returns nil once the window is closed.
func Run(cfg config.Config) error {
	ebiten.SetWindowTitle(cfg.Title)
	ebiten.SetWindowSize(cfg.WindowSize())
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(New(cfg)); err != nil {
		return fmt.Errorf("run game: %w", err)
	}
	return nil
}
