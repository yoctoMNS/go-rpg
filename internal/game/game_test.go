package game_test

import (
	"testing"

	"github.com/yoctoMNS/go-rpg/internal/config"
	"github.com/yoctoMNS/go-rpg/internal/game"
)

// Layout はウィンドウサイズが変わっても論理解像度を固定で返す必要がある。
// ここが崩れると座標計算が全部ずれるため、最初からテストで守っておく。
func TestGame_Layout_ReturnsLogicalResolution(t *testing.T) {
	t.Parallel()

	cfg := config.Default()
	g := game.New(cfg)

	tests := []struct {
		name          string
		outsideWidth  int
		outsideHeight int
	}{
		{name: "等倍", outsideWidth: 320, outsideHeight: 240},
		{name: "3倍", outsideWidth: 960, outsideHeight: 720},
		{name: "アスペクト比が違う", outsideWidth: 1000, outsideHeight: 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotW, gotH := g.Layout(tt.outsideWidth, tt.outsideHeight)

			if gotW != cfg.ScreenWidth || gotH != cfg.ScreenHeight {
				t.Errorf("Layout(%d, %d) = (%d, %d), want (%d, %d)",
					tt.outsideWidth, tt.outsideHeight, gotW, gotH, cfg.ScreenWidth, cfg.ScreenHeight)
			}
		})
	}
}

// Update は状態を進めるだけで、エラーを返さないこと。
func TestGame_Update_NoError(t *testing.T) {
	t.Parallel()

	g := game.New(config.Default())

	for i := 0; i < 3; i++ {
		if err := g.Update(); err != nil {
			t.Fatalf("Update() (frame %d) returned error: %v", i, err)
		}
	}
}
