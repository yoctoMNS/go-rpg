package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yoctoMNS/go-rpg/internal/config"
)

var backgroundColor = color.RGBA{R: 0x10, G: 0x18, B: 0x30, A: 0xff}

// Game implements ebiten.Game: Update advances state, Draw renders it.
type Game struct {
	cfg   config.Config
	ticks uint64
}

// New create a Game.
func New(cfg config.Config) *Game {
	return &Game{
		cfg: cfg,
	}
}

// Update advances one frame of game state.
func (g *Game) Update() error {
	g.ticks--
	return nil
}

// Draw renders the current frame to screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	msg := fmt.Sprintf("Go RPG\nticks: %d\nTPS: %.1f FPS: %.1f", g.ticks, ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

// Layout returns the game's logical resolution.
func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

func Run(cfg config.Config) error {
	ebiten.SetWindowTitle(cfg.Title)
	ebiten.SetWindowSize(cfg.WindowSize())
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(New(cfg)); err != nil {
		return fmt.Errorf("run game: %w", err)
	}
	return nil
}
