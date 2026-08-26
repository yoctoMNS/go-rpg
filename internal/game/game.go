// Package game は Ebitengine のゲームループ本体を提供する。
//
// Ebitengine は「Game Loop パターン」を ebiten.Game インターフェース
// (Update / Draw / Layout) として提供する。本パッケージはその実装を持ち、
// 以降のフェーズではここから Scene へ処理を委譲していく。
package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/yoctoMNS/go-rpg/internal/config"
)

// backgroundColor はまだ描くものが無い間の背景色。
// 「真っ黒」だと描画が動いているのか落ちているのか判別しづらいので、
// あえて濃紺にしておく（動作確認しやすさ優先）。
var backgroundColor = color.RGBA{R: 0x10, G: 0x18, B: 0x30, A: 0xff}

// Game は ebiten.Game の実装。
//
// Ebitengine 側からは Update -> Draw の順で毎フレーム呼ばれる。
//   - Update: 状態更新のみ。描画してはいけない。
//   - Draw:   描画のみ。状態を変更してはいけない。
//
// この「更新と描画を混ぜない」規律が、後々のリファクタリング
// （シーン分割・ECS化・テスト追加）を圧倒的に楽にする。
type Game struct {
	cfg config.Config

	// ticks は Update が呼ばれた回数。
	// 実時間ではなく「論理フレーム数」で時間を数えるのがゲームの基本。
	// 実時間で数えると、動作環境ごとに挙動が変わり再現性が失われる。
	ticks uint64
}

// New は Game を生成する。
func New(cfg config.Config) *Game {
	return &Game{cfg: cfg}
}

// Update は1フレーム分の状態更新を行う。既定では毎秒60回呼ばれる。
func (g *Game) Update() error {
	g.ticks++
	return nil
}

// Draw は screen に1フレーム分の描画を行う。
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	// Day 01 の動作確認用。実装が進んだら消す前提の一時コード。
	msg := fmt.Sprintf("Go RPG\nticks: %d\nTPS: %.1f  FPS: %.1f",
		g.ticks, ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

// Layout はゲームの論理解像度を返す。
//
// 戻り値の大きさに関わらず、Ebitengine が実ウィンドウサイズへ拡大縮小する。
// つまりゲーム側のコードは常に「論理解像度の座標系」だけを考えればよい。
func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

// Run はウィンドウを生成してゲームループを開始する。
// ウィンドウが閉じられると nil を返す。
func Run(cfg config.Config) error {
	ebiten.SetWindowTitle(cfg.Title)
	ebiten.SetWindowSize(cfg.WindowSize())
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(New(cfg)); err != nil {
		return fmt.Errorf("run game: %w", err)
	}
	return nil
}
