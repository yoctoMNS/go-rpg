# Week 02 — 絵を出す（アセットとスプライト）

> **今週のゴール**: 白い四角がキャラクターの絵に変わり、**歩くアニメーションが動く**。

| | |
| --- | --- |
| 導入するパターン | 埋め込みアセット / スプライトシート（Flyweight の予行演習） |
| 追加するパッケージ | `assets` `sprite` |

---

## Day 06 — アセットを実行ファイルに埋め込む

| 項目 | 内容 |
| --- | --- |
| 目安 | 40分 |
| 触るファイル | `assets/assets.go` `assets/images/` `assets/CREDITS.md` |
| 学ぶこと | `//go:embed` / `embed.FS` / 起動時一括ロード |

### やること

1. 素材を用意する（自作でも配布素材でもよい）
   - 推奨: [ぴぽや倉庫](https://pipoya.net/sozai/) などの 32x32 キャラチップ
   - **ライセンスを `assets/CREDITS.md` に必ず記録する**（後から調べ直すのはほぼ不可能）
2. `assets/assets.go` に埋め込み FS を作る
   ```go
   package assets

   import "embed"

   //go:embed images/*.png
   var FS embed.FS
   ```
3. `assets.LoadImage(path string) (*ebiten.Image, error)` を書く
   - `ebitenutil.NewImageFromFileSystem` か `image/png` + `ebiten.NewImageFromImage`
4. `MustLoadImage` も用意する（**埋め込み済み画像のデコード失敗はプログラマのミス** なので panic 可）

### 動作確認

```sh
make run
```

- [ ] 画面の左上にキャラクター画像がそのまま1枚表示される（まだ動かなくてよい）
- [ ] `go build` した実行ファイルを別ディレクトリにコピーしても動く（＝埋め込めている証拠）

### ポイント

- 埋め込みにすると **配布が実行ファイル1つ** で済み、実行時のパス問題が消える
- ゲームループ中にファイルを開かない。読み込みは起動時に全部やる

---

## Day 07 — スプライトシートを切り出す

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/sprite/sheet.go` + テスト |
| 学ぶこと | `SubImage` / 座標計算の切り出し方 |

### やること

1. `sprite.Sheet` を作る
   ```go
   type Sheet struct {
       img        *ebiten.Image
       frameW     int
       frameH     int
       cols, rows int
   }
   func NewSheet(img *ebiten.Image, frameW, frameH int) *Sheet
   func (s *Sheet) Frame(col, row int) *ebiten.Image  // SubImage を返す
   func (s *Sheet) FrameByIndex(i int) *ebiten.Image
   ```
2. **`SubImage` は元画像を共有する**（コピーしない）。だから毎フレーム呼んでも安い……が、
   返り値の `*ebiten.Image` は毎回 alloc されるので、**切り出し結果はキャッシュする**
3. インデックス → (col, row) の変換にテストを書く（範囲外は panic か error か決める）

### 動作確認

```sh
make run
```

- [ ] シートの中の**1コマだけ**（例: 下向き待機）が画面に表示される
- [ ] 表示位置・サイズが 32x32 でぴったり切り出されている（隣のコマが混ざらない）

### ポイント

- 切り出しがズレるときは、素材の余白（padding）やコマ間の隙間（spacing）を疑う
- ここでの `Frame(col, row)` が、後のアニメーションとタイルマップ両方で使い回される

---

## Day 08 — アニメーションを再生する

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/sprite/animation.go` `internal/sprite/animator.go` + テスト |
| 学ぶこと | フレーム管理 / データと再生状態の分離 |

### やること

1. **定義と状態を分ける**（W10 のデータ駆動につながる重要な感覚）
   ```go
   // Animation は「どう再生するか」の定義。共有される。書き換えない。
   type Animation struct {
       Frames   []int // シート上のインデックス列
       Duration int   // 1コマあたりのフレーム数
       Loop     bool
   }

   // Animator は「今どこを再生中か」の状態。インスタンスごとに持つ。
   type Animator struct {
       anim    *Animation
       elapsed int
       index   int
   }
   ```
2. `Animator.Update()` と `Animator.CurrentFrame() int` を実装する
3. `Play(a *Animation)` で切り替え（**同じアニメを渡されたらリセットしない**。これを忘れると
   歩き続けている間ずっと1コマ目に戻り続けるバグになる）
4. テスト: 3コマ・Duration=8 のアニメを 24 回 Update して 1周すること、ループ / 非ループの差

### 動作確認

```sh
make run
```

- [ ] キャラクターがその場で足踏みアニメーションをする
- [ ] `Duration` を 4 と 16 に変えると、明らかに速さが変わる

### ポイント

- アニメーションの進行は必ず `Update` で。`Draw` に書くと処理落ち時に速度が変わる
- テストが書けるのは、`Animator` が `ebiten` に依存していないから。**この分離を意識する**

---

## Day 09 — 移動とアニメーションをつなぐ

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/entity/player.go` `internal/entity/direction.go` |
| 学ぶこと | 型付き列挙 / 状態と表示の対応づけ |

### やること

1. `entity.Direction` を型付き定数で定義する（`DirDown` `DirLeft` `DirRight` `DirUp`）
   - **キャラチップの行順に合わせて iota を振る**と、`int(dir)` がそのまま行番号になる
   - `String()` メソッドを付ける（デバッグ表示が一気に楽になる）
2. `Player` に `dir Direction` と `animator *sprite.Animator` を持たせる
3. 移動入力から向きを決める（**斜めのときは左右を優先**するのが2Dアクションの定石）
4. 速度がゼロなら待機アニメ、動いていれば歩きアニメを再生する
5. `Draw` で `DrawImageOptions` を使ってスプライトを描く
   - **`DrawImageOptions` は使い回す**（`opts.GeoM.Reset()` してから再利用）

### 動作確認

```sh
make run
```

- [ ] 十字キーで押した方向にキャラが向く
- [ ] 動いている間だけ歩行アニメが再生され、離すと待機に戻る
- [ ] 斜め入力でも向きがガタつかない

### ポイント

- ここで `if` の分岐が増え始める。**まだ我慢する**（Day 11 で State パターンに置き換える）
- 「痛みを感じてから直す」ことで、パターンの必要性が体で分かる

---

## Day 10 — 今週のリファクタリング

| 項目 | 内容 |
| --- | --- |
| 目安 | 40分 |
| 触るファイル | 今週触った全ファイル |
| 学ぶこと | リファクタリングの進め方（機能を足さない日） |

### やること

1. `assets` の読み込みを1か所にまとめる（`assets.Load()` で全部読んで構造体で返す）
   - 呼ぶ側が個別のパスを知らなくて済む状態にする
2. `Player` の `Draw` が長くなっていたら、描画オプション作成を切り出す
3. マジックナンバーを名前付き定数にする（キャラのサイズ、アニメ速度、移動速度）
4. 公開している型・関数に doc コメントが付いているか確認する
5. `make check` を通し、**画面の挙動が Day 09 と1ミリも変わっていない**ことを確認する

### 動作確認

```sh
make run
```

- [ ] Day 09 とまったく同じ動作をする（リファクタなので当然。ここが確認事項）

### 今週の振り返り（ログに書く）

- `Player` の `Update` は何行になったか。50行を超えていないか
- 「向き」と「アニメ」の対応づけで気持ち悪かった箇所はどこか（→ 来週の State パターンの動機）
