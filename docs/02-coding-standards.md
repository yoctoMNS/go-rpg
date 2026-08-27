# コーディング規約

前提として [Effective Go](https://go.dev/doc/effective_go) と
[Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) に従う。
ここには **それに加えて、このプロジェクト（ゲーム開発）で特に守ること** だけを書く。

---

## 1. 命名

| 対象 | ルール | 例 |
| --- | --- | --- |
| パッケージ | 小文字1語。複数形にしない。`util` / `common` / `helper` は禁止 | `player`, `tilemap`, `battle` |
| 型 | 名詞。パッケージ名と重複させない | `player.State`（`player.PlayerState` は×） |
| インターフェース | 「〜するもの」。1メソッドなら `-er` | `Drawer`, `input.Provider` |
| 定数 | キャメルケース。関連するものは型付き定数にまとめる | `DirUp`, `DirDown` |
| 受け手(レシーバ) | 型名の1〜2文字。`this` / `self` は禁止 | `func (g *Game) Update()` |

### ゲーム特有の略語は使ってよいが、統一する

`HP` / `MP` / `EXP` / `NPC` / `AI` / `SE` / `BGM` は略語のまま使う。
Go の慣習に従い **全部大文字**（`ID`, `URL` と同じ扱い）。

```go
type Stats struct {
    HP    int // ○
    Hp    int // ×
    MaxHP int // ○
}
```

---

## 2. パッケージ設計

- **`internal/` に入れる。** 外部から import されない前提のコードを `internal/` 以外に置かない
- **`util` パッケージを作らない。** 置き場所に困るコードは、たいてい所属すべき型がある
- **循環 import を作らない。** 依存の向きは常に「具体 → 抽象」「上位 → 下位」
- 1パッケージのファイル数が10を超えたら分割を検討する

依存の向き（上が上位）:

```
cmd/rpg
   ↓
internal/game      … ゲームループ
   ↓
internal/scene     … 画面ごとの状態
   ↓
internal/entity, internal/tilemap, internal/ui, internal/battle
   ↓
internal/input, internal/assets, internal/config, internal/geom
```

**下位パッケージが上位を import したら設計ミス。** 必要になったらインターフェースを
「使う側（上位）」に定義して逆転させる（依存性逆転）。

---

## 3. Update と Draw を混ぜない

Ebitengine で最も重要な規律。

```go
// ○ 正しい
func (g *Game) Update() error { g.player.Update(); return nil }
func (g *Game) Draw(screen *ebiten.Image) { g.player.Draw(screen) }

// × Draw で状態を変える（フレームスキップ時に挙動が壊れる）
func (g *Game) Draw(screen *ebiten.Image) {
    g.animFrame++ // これを Update に移す
    ...
}
```

理由: Ebitengine は処理落ちしたとき `Draw` をスキップすることがある。
`Draw` に状態更新を書くと、重い環境でだけゲームが遅くなるバグが生まれる。

---

## 4. 時間はフレーム数で数える

`time.Now()` や実時間の差分でゲームロジックを動かさない。
Ebitengine の `Update` は既定で 60TPS 固定なので、**フレーム数を数える**。

```go
// ○ 再現性がある
const invincibleFrames = 60 // 1秒

// × 環境で挙動が変わる／リプレイやテストが再現しない
if time.Since(g.hitAt) > time.Second { ... }
```

---

## 5. 毎フレーム走るコードでメモリを確保しない

`Update` / `Draw` の中で `make` / `new` / 文字列結合をすると、GC が走って
カクつき（フレーム落ち）の原因になる。

```go
// × 毎フレーム slice を確保
func (g *Game) Update() error {
    visible := []*Entity{}          // 毎フレーム alloc
    for _, e := range g.entities { ... }
}

// ○ フィールドに持って使い回す
type Game struct {
    visibleBuf []*Entity // 使い回すバッファ
}

func (g *Game) Update() error {
    g.visibleBuf = g.visibleBuf[:0] // 長さだけリセット。容量は再利用される
    for _, e := range g.entities { ... }
}
```

同様に `*ebiten.Image` と `*ebiten.DrawImageOptions` も使い回す。

---

## 6. エラーハンドリング

- **起動時（アセット読み込み・データ読み込み）のエラーは即座に返して落とす。** 握り潰さない
- **ゲームループ中のエラーは基本的に起こさない設計にする。** 毎フレーム失敗しうる処理をループに置かない
- エラーは `fmt.Errorf("...: %w", err)` でラップし、文脈を足す
- エラーメッセージは小文字始まり・句点なし（Go の慣習）

```go
func LoadTileset(path string) (*Tileset, error) {
    f, err := assets.FS.Open(path)
    if err != nil {
        return nil, fmt.Errorf("open tileset %q: %w", path, err)
    }
    defer f.Close()
    ...
}
```

`panic` を使ってよいのは「プログラマのミスであり、直さない限り絶対に動かない」場合のみ
（例: `//go:embed` したアセットのデコード失敗）。その場合は `MustXxx` という名前にする。

---

## 7. コメント

- **コメントは英語で書く。** Go の標準ライブラリ・コミュニティの慣習に合わせる
  （このリポジトリの計画書・作業ログは学習記録として日本語のままでよいが、
  ソースコード中のコメントは英語に統一する）
- **公開（大文字始まり）の型・関数には必ず doc コメント。** `// Xxx does ...` の形式
- **「何をしているか」ではなく「なぜそうしたか」を書く。** コードを読めばわかることは書かない
- マジックナンバーには必ず理由を書くか、名前付き定数にする

```go
// Good: explains why.
// Movement speed is 1.5px/frame. Using a value that doesn't evenly divide
// TileSize (16px) avoids getting stuck exactly on tile boundaries.
const walkSpeed = 1.5

// Bad: just restates the code.
// Assign 1.5 to speed.
const walkSpeed = 1.5
```

---

## 8. 構造体の作り方

- **ゼロ値で使えるなら使えるようにする**（`var s Stats` がそのまま有効）
- 初期化が必要なら `New` コンストラクタを用意し、**必須引数は引数で、任意設定は Option で** 受ける
- **公開フィールドは慎重に。** 外から書き換えられて困る値はメソッド経由にする

```go
// 必須は引数、任意は可変長オプション（Functional Options パターン）
func NewPlayer(pos geom.Vec2, opts ...PlayerOption) *Player

// 使う側
p := NewPlayer(start)
p := NewPlayer(start, WithSpeed(2.0))
```

Option を使うのは **引数が4つを超えたとき、または任意設定が増えたとき**。
最初から全部 Option にはしない（過剰設計）。

---

## 9. インターフェースの作り方

- **「使う側」のパッケージに定義する。** 実装側に置かない（Go の重要な慣習）
- **小さく保つ。** メソッド3つを超えたら分割を疑う
- **実装が1つしかないうちはインターフェースを作らない。** テスト用フェイクを作る時点で初めて切る

```go
// internal/scene/scene.go … 使う側に定義
package scene

// InputProvider はシーンが必要とする入力だけを表す。
// ebiten への依存をここで断ち切り、テストではフェイクを差し込む。
type InputProvider interface {
    IsPressed(a input.Action) bool
    IsJustPressed(a input.Action) bool
}
```

---

## 10. テストの書き方

- **テーブル駆動テスト**を基本にする
- テストパッケージは `xxx_test`（外部テストパッケージ）を優先。公開APIだけを使うことで、設計の使いにくさに気づける
- `t.Parallel()` を付ける。`make test` は `-race -shuffle=on` で走るので、状態の共有ミスがすぐ露見する
- 失敗メッセージは `got = X, want Y` の形式で、**何を入れて何が返ったか** がわかるように書く

```go
func TestDamage(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name string
        atk  int
        def  int
        want int
    }{
        {name: "通常", atk: 20, def: 10, want: 10},
        {name: "防御が上回る場合は最低1", atk: 5, def: 100, want: 1},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            if got := Damage(tt.atk, tt.def); got != tt.want {
                t.Errorf("Damage(%d, %d) = %d, want %d", tt.atk, tt.def, got, tt.want)
            }
        })
    }
}
```

---

## 11. アセットの扱い

- 画像・音・データはすべて `assets/` に置き、`//go:embed` で実行ファイルに埋め込む
  → 配布が1ファイルで済み、実行時のパス問題が消える
- 読み込みは **起動時に一括**。ゲームループ中にファイルを開かない
- 素材のライセンスは `assets/CREDITS.md` に必ず記録する（後から調べ直すのは不可能に近い）

---

## 12. やらないこと（このプロジェクトの禁止事項）

- グローバル変数に可変状態を置く（`var currentPlayer *Player` など）
- `init()` で重い処理をする
- 型アサーションの連鎖（`if p, ok := e.(*Player); ok { ... }` が増えたら設計を疑う）
- 1関数100行超え（画面描画でも分割できる）
- 「後で直す」コメントを消えないまま放置する（`TODO:` を書いたらログにも書く）
