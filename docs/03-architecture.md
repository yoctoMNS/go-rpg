# アーキテクチャ方針

## 1. 最終的なパッケージ構成（12週後の姿）

いきなりこの形を作らない。**必要になった週に、そのパッケージを切り出す。**
「今どこまで作られているか」は各週の計画書に書いてある。

```
go-rpg/
├── cmd/
│   └── rpg/
│       └── main.go            # 起動だけ。ロジックを書かない (Composition Root)
├── internal/
│   ├── config/                # 起動時設定（W1）
│   ├── game/                  # ebiten.Game 実装 = ゲームループ（W1）
│   ├── geom/                  # Vec2, Rect など座標・矩形（W2）
│   ├── input/                 # 物理キー → 論理アクションへの変換（W2）
│   ├── assets/                # //go:embed によるアセット読み込み（W3）
│   ├── sprite/                # スプライトシート・アニメーション（W3）
│   ├── entity/                # プレイヤー/NPC/敵などの実体（W3〜）
│   ├── tilemap/               # タイルマップの読み込みと描画（W4）
│   ├── physics/               # 当たり判定・移動解決（W5）
│   ├── camera/                # 描画ビューポートとスクロール（W5）
│   ├── scene/                 # Scene インターフェースと SceneManager（W6）
│   │   ├── title/
│   │   ├── field/
│   │   └── battle/
│   ├── ui/                    # ウィンドウ・テキスト・メニュー（W7）
│   ├── dialogue/              # 会話スクリプトとコマンド列（W7）
│   ├── battle/                # 戦闘ロジック（FSM・行動・AI）（W8-W9）
│   ├── gamedata/              # JSON定義の読み込み（Repository）（W10）
│   ├── save/                  # セーブデータの入出力（W11）
│   ├── audio/                 # BGM/SE（W11）
│   └── event/                 # イベントバス（Observer）（W11）
├── assets/                    # 画像・音・JSON（//go:embed 対象）
├── docs/                      # 計画書とログ（このディレクトリ）
└── scripts/
```

---

## 2. 3つのレイヤ

パッケージは大きく3層に分ける。**依存は必ず上から下へ一方向。**

```
┌──────────────────────────────────────────┐
│ Presentation   scene / ui / camera        │  画面と入出力
├──────────────────────────────────────────┤
│ Domain         battle / entity / tilemap  │  ゲームのルール
├──────────────────────────────────────────┤
│ Infrastructure input / assets / audio /   │  外の世界との接続
│                save / config / geom       │
└──────────────────────────────────────────┘
```

### Domain 層に ebiten を import しない

これがこのプロジェクトで一番大事な設計判断。

```go
// ○ ダメージ計算は ebiten を知らない → テストがそのまま書ける
package battle

func Damage(atk, def int) int { ... }

// × Domain が描画ライブラリに依存すると、テストのたびに GPU が要る
package battle

import "github.com/hajimehoshi/ebiten/v2"
func (b *Battle) Update(screen *ebiten.Image) { ... }
```

「ルール」と「見た目」が分かれていれば、後でライブラリを変えても、CLI 版を作っても、
AI に自動対戦させても、ドメインのコードはそのまま使える。

---

## 3. デザインパターンの適用計画

### 3.1 Game Loop（W1）

Ebitengine が `ebiten.Game` インターフェースとして提供済み。自分で書く必要はないが、
**「Update は状態更新のみ、Draw は描画のみ」** という契約を守る責任は自分にある。

```go
type Game interface {
    Update() error                 // 論理更新（60回/秒固定）
    Draw(screen *ebiten.Image)     // 描画（可変。スキップされうる）
    Layout(w, h int) (int, int)    // 論理解像度
}
```

### 3.2 Input Abstraction / Adapter（W2）

「Zキー」ではなく「決定アクション」でロジックを書く。

```go
package input

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

// Provider は入力の抽象。ebiten 実装のほかにテスト用フェイクを作れる。
type Provider interface {
    IsPressed(Action) bool
    IsJustPressed(Action) bool
}
```

**得られるもの**: キーコンフィグ / ゲームパッド対応 / テスト自動化 / リプレイ機能が
すべてこの1枚のインターフェースの裏側で完結する。

### 3.3 State パターン（W3: アニメーション、W6: シーン、W8: バトル）

このプロジェクトで最も出番の多いパターン。**「状態ごとに構造体を作り、遷移を明示する」。**

```go
// 悪い例: switch が膨らみ続ける
switch p.state {
case "idle":  ...
case "walk":  ...
case "attack": ...
case "damage": ... // 追加のたびここを触る = 既存を壊すリスク
}

// 良い例: 状態が自分の振る舞いと次の遷移先を知っている
type State interface {
    Enter(*Player)
    Update(*Player) State // 次の状態を返す。同じなら自分を返す
    Exit(*Player)
}
```

シーン管理では **Scene Stack** にする（スタックにするとメニューを「重ねて」開ける）。

```go
type Manager struct {
    stack []Scene
}

func (m *Manager) Push(s Scene) // フィールドの上にメニューを重ねる
func (m *Manager) Pop()         // メニューを閉じてフィールドに戻る
func (m *Manager) Replace(s Scene) // タイトル → フィールド（戻らない遷移）
```

### 3.4 Flyweight（W4）

タイルは同じ絵が何千個も並ぶ。**「見た目の定義」は共有し、「配置」だけを持つ。**

```go
// 共有される内在状態（Intrinsic）: タイルの種類ごとに1つだけ存在
type TileDef struct {
    ID       int
    SrcRect  image.Rectangle // タイルセット画像内の位置
    Passable bool
}

// マップは TileDef のIDだけを持つ（外在状態 = Extrinsic）
type Layer struct {
    Width, Height int
    Tiles         []int // len == Width*Height。ここは int のみ
}
```

### 3.5 Command パターン（W7: 会話、W9: 戦闘行動）

**「やること」をオブジェクトにする。** 会話イベントも戦闘行動も、これで統一的に扱える。

```go
// 会話スクリプトを命令列として表現する
type Command interface {
    // Execute は1フレーム分進め、完了したら true を返す。
    Execute(ctx *Context) (done bool)
}

type ShowText struct{ Text string }
type Wait struct{ Frames int }
type GiveItem struct{ ItemID string; Count int }
type Choice struct{ Options []string; Branches [][]Command }
```

**得られるもの**: 会話イベントを JSON で書けるようになり、シナリオ追加でコードを触らなくなる。
戦闘では「行動をキューに積んで素早さ順にソートして実行」がそのまま書ける。
さらに「実行した Command を記録する」だけで、リプレイと Undo が手に入る。

### 3.6 Strategy（W9）

敵AIの思考ルーチンを差し替え可能にする。

```go
type AI interface {
    // Decide はこのターンの行動を決める。
    Decide(self *Enemy, field *Field) Action
}

type AggressiveAI struct{} // 常に最大ダメージを狙う
type HealerAI struct{}     // 味方のHPが半分以下なら回復
type RandomAI struct{ rng *rand.Rand }
```

敵データの JSON に `"ai": "healer"` と書くだけで挙動が変わる状態を目指す。

### 3.7 Repository + Factory + データ駆動（W10）

**敵・アイテム・スキルの追加でコードを書き換えない。**

```go
// gamedata パッケージ: JSON を読んで「定義」を保持する
type EnemyRepository interface {
    Find(id string) (EnemyDef, error)
    All() []EnemyDef
}

// Factory: 「定義」から「実体」を作る
func NewEnemy(def EnemyDef) *entity.Enemy
```

定義（`EnemyDef`: 変わらないマスタデータ）と
実体（`Enemy`: HPが減る実行時の状態）を **必ず別の型にする**。
ここを混ぜると「マスタデータのHPが減っていた」という致命的バグが生まれる。

### 3.8 Observer / Event Bus（W11）

「HPが0になった」を、UI・音・実績・セーブが同時に知りたくなる。
呼び出し元が全部を知る設計にすると、依存が爆発する。

```go
package event

type Event interface{ eventMarker() }

type PlayerDamaged struct{ Amount int }
type EnemyDefeated struct{ EnemyID string }
type LevelUp       struct{ NewLevel int }

type Bus struct { ... }
func (b *Bus) Subscribe(handler func(Event))
func (b *Bus) Publish(e Event)
```

**注意**: イベントバスは便利すぎて乱用しがち。「誰がいつ反応するか」が追えなくなる。
**直接呼べる関係は直接呼ぶ。** 1対多で、かつ呼ぶ側が相手を知りたくない場合だけ使う。

### 3.9 Component 指向（W12）

`Player` / `NPC` / `Enemy` が似たフィールドを持ち始めたら、**継承ではなく合成** で解く。

```go
// 「性質」を独立した型にして持たせる
type Transform struct{ Pos, Vel geom.Vec2 }
type Sprite    struct{ Sheet *sprite.Sheet; Anim *sprite.Animator }
type Collider  struct{ Box geom.Rect; Solid bool }
type Health    struct{ HP, MaxHP int }

type Entity struct {
    Transform Transform
    Sprite    *Sprite   // nil = 描画しない
    Collider  *Collider // nil = 当たり判定なし
    Health    *Health   // nil = 無敵/破壊不能
}
```

W12 ではここまで（Component 指向）で止める。**本格的な ECS（Entity Component System）は、
エンティティ数が数千を超えて実際に処理が重くなってから。** それ以前は複雑さの負債にしかならない。

---

## 4. アンチパターン（避けること）

| アンチパターン | なぜダメか | 代わりに |
| --- | --- | --- |
| God Object（`Game` が全部持つ） | 変更のたび全体に影響。テスト不能 | Scene に分割（W6） |
| Service Locator / グローバル変数 | 依存が隠れる。テストで差し替え不能 | コンストラクタで注入 |
| 深い継承もどき（埋め込みの多段） | Go では特に追いづらい | Component で合成 |
| 全部イベントバス経由 | 実行順が追えず、デバッグ不能に | 直接呼べるものは直接呼ぶ |
| 早すぎるECS化 | 必要ない複雑さ。速度も出ない | まず構造体、重くなってから |
| `interface{}` / `any` の多用 | 型安全性を捨てている | ジェネリクスか具体型 |

---

## 5. 判断に迷ったときの原則

1. **動くものを先に作る。** 設計はあとで直せるが、動かないコードからは何も学べない
2. **痛みが出てからパターンを入れる。** 「将来のため」の抽象化は、たいてい間違った抽象化になる
3. **消しやすさ > 再利用性。** 再利用しやすいコードより、捨てやすいコードの方が価値が高い
4. **テストが書きづらいなら設計が悪い。** テストのしづらさは設計の欠陥を教えてくれる唯一のシグナル
