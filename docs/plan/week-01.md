# Week 01 — 土台を作り、四角形を動かす

> **今週のゴール**: `make run` でウィンドウが開き、**十字キーで四角形が動く**。
> RPGの見た目はまだ無いが、「入力 → 更新 → 描画」というゲームの心臓部が完成する。

| | |
| --- | --- |
| 導入するパターン | Game Loop / Input Abstraction (Adapter) |
| 追加するパッケージ | `config` `game` `geom` `input` `entity` |

---

## Day 01 — プロジェクトの骨組みとウィンドウ表示

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `go.mod` `cmd/rpg/main.go` `internal/config/config.go` `internal/game/game.go` `Makefile` |
| 学ぶこと | `ebiten.Game` インターフェース / 論理解像度 / Composition Root |

### やること

1. Go モジュールを作る
   ```sh
   go mod init github.com/yoctoMNS/go-rpg
   go get github.com/hajimehoshi/ebiten/v2@latest
   ```
2. `internal/config/config.go` に `Config` 構造体と `Default()` を書く
   - 論理解像度 320x240、ウィンドウ倍率 3、タイルサイズ 16
   - **なぜ 320x240 か**: 16px タイルが 20x15 個ちょうど収まり、整数倍拡大で滲まないため
3. `internal/game/game.go` に `Game` 型を作り、`Update` / `Draw` / `Layout` を実装する
   - `Draw` では背景を塗り、`ebitenutil.DebugPrint` で TPS/FPS を出す
4. `cmd/rpg/main.go` は **起動だけ**。`config.Default()` を作って `game.Run(cfg)` を呼ぶ
5. `Makefile` を用意して `make run` / `make check` を作る

### 動作確認

```sh
make run
```

- [ ] 960x720 のウィンドウが開き、タイトルが `Go RPG`
- [ ] 左上に `ticks` が増え続け、`TPS: 60.0` 前後が出る
- [ ] ウィンドウをドラッグしてリサイズしても、文字の大きさが**画面に対する比率のまま**変わる
      （＝ `Layout` が論理解像度を固定できている証拠）

### ポイント

- `Update` は **状態更新のみ**、`Draw` は **描画のみ**。この規律を Day 01 から徹底する
- 時間は `time.Now()` ではなく **フレーム数（ticks）** で数える

---

## Day 02 — テストと CI を通す

| 項目 | 内容 |
| --- | --- |
| 目安 | 40分 |
| 触るファイル | `internal/config/config_test.go` `internal/game/game_test.go` `.github/workflows/ci.yml` |
| 学ぶこと | テーブル駆動テスト / 外部テストパッケージ / 「テストできる形」に切る感覚 |

### やること

1. `config_test.go` を書く（パッケージ名は `config_test` にする＝公開APIだけを使う）
   - `WindowSize()` が倍率どおりか
   - 画面サイズがタイルサイズで割り切れるか
2. `game_test.go` で `Layout` をテーブル駆動テストする
   - ウィンドウサイズが変わっても論理解像度が固定で返ること
3. `.github/workflows/ci.yml` を追加する
   - Linux で Ebitengine をビルドするには X11/OpenGL の開発ヘッダが必要（`libx11-dev` ほか）
4. `make check` が通ることを確認する

### 動作確認

```sh
make check   # fmt + vet + test がすべて通る
make run     # 壊れていないことを目視確認
```

- [ ] `ok github.com/yoctoMNS/go-rpg/internal/config` が出る
- [ ] プッシュ後、GitHub Actions が緑になる

### ポイント

- **`Draw` はテストしない。** 目で見た方が速い。代わりに `Layout` のような「数値を返す関数」を厚くテストする
- テストが書きづらいと感じたら、それは設計が悪いサイン（この先ずっと効いてくる判断基準）

---

## Day 03 — 座標の型を作る（geom パッケージ）

| 項目 | 内容 |
| --- | --- |
| 目安 | 40分 |
| 触るファイル | `internal/geom/vec2.go` `internal/geom/rect.go` + それぞれのテスト |
| 学ぶこと | 値型の設計 / ゼロ値で使える型 / 「小さく正しい部品」を持つ意味 |

### やること

1. `geom.Vec2` を作る（`X, Y float64`）
   - `Add` `Sub` `Scale` `Len` `Normalize` をメソッドで
   - **すべて値レシーバ**にして、新しい値を返す（イミュータブル）。ポインタにすると
     `a.Add(b)` が `a` を書き換えるのか分からず、バグの温床になる
2. `geom.Rect` を作る（`X, Y, W, H float64`）
   - `Intersects(other Rect) bool` `Contains(p Vec2) bool` `Move(d Vec2) Rect`
3. 両方にテーブル駆動テストを書く（境界がぴったり接する場合を必ず入れる）

### 動作確認

```sh
make test
```

- [ ] 「1pxだけ重なる」「辺がぴったり接する（重なっていない扱い）」のテストが通る

### ポイント

- なぜ `image.Point` を使わないか → **整数だと滑らかな移動ができない**。ゲームは float64 で計算し、描画のときだけ丸める
- ここで作る `Intersects` が、W5 の当たり判定でそのまま使われる。**土台の正しさは後で効く**

---

## Day 04 — 入力を抽象化する（input パッケージ）

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/input/action.go` `internal/input/keyboard.go` `internal/input/fake.go` |
| 学ぶパターン | **Adapter / 依存性逆転** |

### やること

1. `input.Action` を型付き定数で定義する
   ```go
   type Action int
   const (
       ActionUp Action = iota
       ActionDown
       ActionLeft
       ActionRight
       ActionConfirm
       ActionCancel
       ActionMenu
   )
   ```
2. `input.Provider` インターフェースを定義する（`IsPressed` / `IsJustPressed` の2つだけ）
3. `Keyboard` 実装を書く。キーマップは `map[Action][]ebiten.Key` で持つ
   - 十字キーと WASD の両方を割り当てる
   - `IsJustPressed` は `inpututil.IsKeyJustPressed` を使う
4. **`Fake` 実装をテスト用に書く**（`fake.go`。`Press(a Action)` で押した状態を作れる）
5. `Game` に `input.Provider` をコンストラクタで渡す（グローバル変数にしない）
6. `Draw` で「今押されているアクション名」を画面に出す

### 動作確認

```sh
make run
```

- [ ] 十字キーを押すと画面に `UP` `DOWN` などが表示される
- [ ] WASD でも同じ表示になる
- [ ] Z を押すと `CONFIRM`、X で `CANCEL` が出る

### ポイント

- ゲームロジックに `ebiten.KeyArrowUp` を**二度と書かない**。これがキーコンフィグ・パッド対応・
  テスト自動化・リプレイ機能のすべての土台になる
- `Fake` を先に作ることで、「インターフェースが本当に使いやすいか」が即座に分かる

---

## Day 05 — 四角形を動かす + 今週のリファクタリング

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/entity/player.go` `internal/game/game.go` |
| 学ぶこと | 責務の分離 / 速度と位置の分け方 |

### やること

1. `entity.Player` を作る
   ```go
   type Player struct {
       Pos   geom.Vec2
       Speed float64 // px/frame
   }
   func (p *Player) Update(in input.Provider)
   func (p *Player) Draw(screen *ebiten.Image)
   ```
2. `Update` で入力を速度に変換し、位置に足す
   - **斜め移動は正規化する**（そうしないと斜めが約1.41倍速くなる）
3. `Draw` で 16x16 の白い四角を描く（`vector.DrawFilledRect` が手軽）
4. 画面外に出ないようクランプする
5. **リファクタ**: `Game.Update` / `Game.Draw` から中身を `Player` へ移し、`Game` は
   「持っているものに委譲するだけ」の状態にする

### 動作確認

```sh
make run
```

- [ ] 十字キーで四角形が上下左右に動く
- [ ] 斜めに動いても速度が速くならない
- [ ] 画面の端で止まり、外に出ない

### ポイント

- **`Game` が薄くなったか確認する。** `Game.Update` が3行程度なら成功
- ここまでで「入力 → 更新 → 描画」の1周が完成。以降はこの型に載せていくだけ

### 今週の振り返り（ログに書く）

- `Game` に書きたくなった処理はあったか？ それはどこに置くべきだったか
- `input.Provider` のメソッドは足りたか、多すぎたか
