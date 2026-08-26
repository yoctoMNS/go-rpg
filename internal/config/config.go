// Package config はゲーム全体で共有する起動時設定を提供する。
//
// 「マジックナンバーをコードに散らさない」ことが目的。
// 値をここに集約しておくと、後から設定ファイル(JSON/TOML)読み込みに
// 差し替えるときも Config 構造体の生成方法を変えるだけで済む。
package config

// Config はゲームの起動時設定を表す。
//
// 実行中に変化する値（プレイヤーHPなど）は入れない。
// あくまで「起動時に決まり、以後変わらない」ものだけを持たせる。
type Config struct {
	// Title はウィンドウタイトル。
	Title string

	// ScreenWidth / ScreenHeight は論理解像度（ゲーム内部の座標系）。
	// ウィンドウサイズとは独立しており、Ebitengine が拡大縮小して描画する。
	// ドット絵RPGでは低めの論理解像度を固定し、整数倍で拡大するのが定石。
	ScreenWidth  int
	ScreenHeight int

	// WindowScale は論理解像度に対する初期ウィンドウ倍率。
	WindowScale int

	// TileSize は1タイルのピクセルサイズ。マップ・当たり判定の基準単位。
	TileSize int
}

// Default は開発時の既定値を返す。
//
// 320x240 (QVGA) はファミコン〜スーファミ系RPGの見た目に近く、
// 16pxタイルで 20x15 タイルがちょうど収まる扱いやすいサイズ。
func Default() Config {
	return Config{
		Title:        "Go RPG",
		ScreenWidth:  320,
		ScreenHeight: 240,
		WindowScale:  3,
		TileSize:     16,
	}
}

// WindowSize は初期ウィンドウの実ピクセルサイズを返す。
func (c Config) WindowSize() (width, height int) {
	return c.ScreenWidth * c.WindowScale, c.ScreenHeight * c.WindowScale
}
