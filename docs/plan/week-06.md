# Week 06 — シーン管理（State パターンの本命）

> **今週のゴール**: タイトル画面 → フィールド → メニューを行き来できる。
> `Game` が痩せて、画面ごとにコードが分かれる。

| | |
| --- | --- |
| 導入するパターン | **State（Scene）** / **Scene Stack** / 依存性注入 |
| 追加するパッケージ | `scene` `scene/title` `scene/field` |

---

## Day 26 — Scene インターフェースと Manager

| 項目 | 内容 |
| --- | --- |
| 目安 | 55分 |
| 触るファイル | `internal/scene/scene.go` `internal/scene/manager.go` + テスト |
| 学ぶパターン | **State + Stack** |

### やること

1. `Scene` インターフェースを定義する
   ```go
   // Scene は1つの画面（タイトル、フィールド、バトル、メニュー）を表す。
   type Scene interface {
       // OnEnter はシーンがスタックに積まれたときに1度だけ呼ばれる。
       OnEnter()
       // Update は1フレーム分の更新を行う。
       Update() error
       // Draw はこのシーンを描画する。
       Draw(screen *ebiten.Image)
       // OnExit はシーンがスタックから降ろされるときに1度だけ呼ばれる。
       OnExit()
   }
   ```
2. `Manager` を **スタック**で作る
   ```go
   func (m *Manager) Push(s Scene)     // 上に重ねる（メニューを開く）
   func (m *Manager) Pop()             // 降ろす（メニューを閉じる）
   func (m *Manager) Replace(s Scene)  // 差し替える（タイトル→フィールド）
   ```
3. **更新は最上段だけ、描画は下から全部**（メニューの後ろにフィールドが透けて見える）
   - ただしシーンが「下も更新してよい」と宣言できる余地を残す（`Transparent()` などは後で）
4. **遷移をフレームの途中でやらない。** 遷移要求をキューに積み、`Update` の最後にまとめて適用する
   - スライスを走査中に書き換えると壊れる。これは実務でも頻出のバグ
5. `Manager` のテストを書く（フェイク Scene で Push/Pop/Replace の呼ばれ順を検証）

### 動作確認

```sh
make test
```

- [ ] `OnEnter` → `Update` → `OnExit` が正しい順で呼ばれるテストが通る

---

## Day 27 — FieldScene に既存機能を移す

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/scene/field/field.go` `internal/game/game.go` |
| 学ぶこと | 大きなリファクタリングの安全な進め方 |

### やること

1. `Game` が持っていたもの（マップ、プレイヤー、NPC、カメラ）を `FieldScene` に移す
2. `Game` は `*scene.Manager` を1つ持つだけにする
   ```go
   func (g *Game) Update() error            { return g.scenes.Update() }
   func (g *Game) Draw(screen *ebiten.Image) { g.scenes.Draw(screen) }
   ```
3. **依存はコンストラクタで渡す**（グローバル変数を作らない）
   ```go
   func NewFieldScene(deps Deps) *FieldScene
   type Deps struct {
       Input  input.Provider
       Assets *assets.Store
       Config config.Config
   }
   ```
4. 一気に動かそうとせず、**コンパイルを通しながら少しずつ移す**

### 動作確認

```sh
make run
```

- [ ] Week 05 とまったく同じ動作をする
- [ ] `Game` 構造体のフィールドが1〜2個になった

### ポイント

- **`Deps` 構造体で依存をまとめる**と、引数が増えても呼び出し側が壊れない
- ここで Service Locator（グローバルから取ってくる方式）を選ばないこと。楽だがテストが死ぬ

---

## Day 28 — タイトル画面を作る

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/scene/title/title.go` |
| 学ぶこと | シーン遷移の実装 |

### やること

1. `TitleScene` を作る（タイトルロゴ + 「PRESS Z TO START」の点滅）
   - 点滅は `ticks/30%2 == 0` で。**実時間を使わない**
2. `ActionConfirm` で `Replace(NewFieldScene(...))` する
3. 起動時のシーンを `TitleScene` に変える

### 動作確認

```sh
make run
```

- [ ] 起動するとタイトル画面が出る
- [ ] 文字が点滅する
- [ ] Z キーでフィールドに切り替わる

---

## Day 29 — シーン遷移演出（フェード）

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/scene/transition.go` |
| 学ぶこと | **Decorator** 的な発想 / 演出とロジックの分離 |

### やること

1. `Transition` を「シーンを包むシーン」として実装する
   ```go
   // FadeTransition は from から to への切り替えを黒フェードで繋ぐ。
   // 自身も Scene なので、Manager から見れば普通のシーンでしかない。
   type FadeTransition struct {
       from, to Scene
       frames   int
       elapsed  int
   }
   ```
2. 黒い矩形の α を 0→255→0 で変化させる
3. `Manager.Replace` を包んで、遷移が常にフェードするヘルパーを作る

### 動作確認

```sh
make run
```

- [ ] タイトル → フィールドの切り替えが黒フェードで繋がる
- [ ] フェード中にキー入力を受け付けない（多重遷移が起きない）

### ポイント

- **演出をシーン自体に書かない**。「シーンを包むもの」にすることで、
  どのシーン遷移にも同じ演出を使い回せる（Decorator の考え方）

---

## Day 30 — 今週のリファクタリングと通し確認

| 項目 | 内容 |
| --- | --- |
| 目安 | 40分 |

### やること

1. シーン間で共有する状態（`GameState`: 所持金、パーティ、フラグ）の置き場所を決める
   - **`Deps` に含めて注入する**。グローバルにしない
   - まだ中身は空でよい。**置き場所だけ決めておく**
2. `scene` パッケージの doc コメントを整える
3. シーン遷移図を `docs/03-architecture.md` に追記する
4. `make check` を通す

### 動作確認

```sh
make run
```

- [ ] タイトル → フィールド → (F1デバッグ) → 移動 → 壁判定 がすべて動く

### 今週の振り返り（ログに書く）

- `Game` は何行になったか（20行以下なら理想的）
- Scene を足すのに必要な作業は何ステップか
