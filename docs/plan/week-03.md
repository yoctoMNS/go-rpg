# Week 03 — State パターンとエンティティ設計

> **今週のゴール**: キャラの状態管理が `if` の山から **State パターン**に置き換わり、
> NPC が画面に立ち、勝手にうろうろする。

| | |
| --- | --- |
| 導入するパターン | **State パターン** / インターフェースによる多態 |
| 追加するパッケージ | `entity/state`（または `entity` 内） |

---

## Day 11 — State パターンを導入する

| 項目 | 内容 |
| --- | --- |
| 目安 | 55分 |
| 触るファイル | `internal/entity/state.go` `internal/entity/player.go` |
| 学ぶパターン | **State** |

### やること

1. まず **Day 09 の `Player.Update` を読み返す**。`if` がいくつあるか数えてログに書く（動機の記録）
2. `State` インターフェースを定義する
   ```go
   // State はプレイヤーの1つの状態を表す。
   // Update は次フレームの状態を返す。状態が変わらないなら自分自身を返す。
   type State interface {
       Enter(*Player)
       Update(*Player, input.Provider) State
       Exit(*Player)
   }
   ```
3. `IdleState` と `WalkState` を実装する
   - `Enter` でアニメ再生開始、`Update` で入力を見て遷移判定
4. `Player` は `state State` を1つ持つだけにする
   ```go
   func (p *Player) Update(in input.Provider) {
       next := p.state.Update(p, in)
       if next != p.state {
           p.state.Exit(p)
           p.state = next
           next.Enter(p)
       }
   }
   ```
5. **遷移テストを書く**（`Fake` 入力で「押す→Walk」「離す→Idle」を検証）

### 動作確認

```sh
make run
```

- [ ] Day 09 とまったく同じ動作をする（内部だけが変わった）
- [ ] `make test` で状態遷移テストが通る

### ポイント

- **状態が自分の遷移先を知っている**のがミソ。新しい状態を足しても既存の状態を触らずに済む
- `Enter` / `Exit` があることで「状態に入った瞬間だけやること」（アニメリセット、無敵開始など）を
  取りこぼさない

---

## Day 12 — ダッシュ状態を足して、パターンの効果を体感する

| 項目 | 内容 |
| --- | --- |
| 目安 | 35分 |
| 触るファイル | `internal/entity/state_dash.go` |
| 学ぶこと | 開放閉鎖の原則（OCP）の実感 |

### やること

1. `DashState` を追加する（`ActionCancel` = X キー押しっぱなしで速度2倍）
2. `WalkState.Update` に「X が押されたら `DashState` へ」の1行だけを足す
3. **既存の `IdleState` を1行も触らずに済んだこと**をログに記録する

### 動作確認

```sh
make run
```

- [ ] X を押しながら移動すると速く動く
- [ ] X を離すと通常速度に戻る
- [ ] 止まると待機アニメになる（Idle への遷移が壊れていない）

### ポイント

- **これがパターンを使う理由**。Day 09 のまま `if` を足していたら、既存の分岐すべてに
  影響が出ていないか確認する必要があった
- 追加コスト（新ファイル1つ + 既存1行）を実感としてログに残す

---

## Day 13 — エンティティの共通部分を切り出す

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/entity/entity.go` `internal/entity/npc.go` |
| 学ぶこと | **継承ではなく合成** / Component 指向の入口 |

### やること

1. `Player` から「位置・向き・スプライト・アニメータ」を `Actor` として切り出す
   ```go
   // Actor はマップ上に存在して絵を持つものの共通部分。
   type Actor struct {
       Pos      geom.Vec2
       Dir      Direction
       Sheet    *sprite.Sheet
       Animator *sprite.Animator
   }
   ```
2. `Player` は `Actor` を **埋め込みではなくフィールドで持つ**
   ```go
   type Player struct {
       Actor Actor // 埋め込み(Actor)にしない
       state State
   }
   ```
   - **なぜ埋め込まないか**: 埋め込みは「継承っぽい」が、メソッドが暗黙に昇格して
     どこで何が起きているか追えなくなる。**明示的な合成を選ぶ**
3. 同じ `Actor` を使う `NPC` 型を作り、マップに1体立たせる

### 動作確認

```sh
make run
```

- [ ] プレイヤーの挙動が変わっていない
- [ ] NPC が画面のどこかに立っている（アニメは待機のみ）

### ポイント

- ここでの `Actor` が W12 の Component 化の原型になる
- 「共通部分を親クラスに上げる」のではなく「共通部分を部品として持たせる」

---

## Day 14 — NPC を自律的に動かす

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/entity/npc.go` `internal/entity/behavior.go` |
| 学ぶパターン | **Strategy**（の予行演習） |

### やること

1. `Behavior` インターフェースを定義する
   ```go
   // Behavior は NPC の行動パターンを表す。
   type Behavior interface {
       Update(*NPC)
   }
   ```
2. 実装を2つ作る
   - `StandStill{}` — その場で待機
   - `RandomWalk{ Interval int; rng *rand.Rand }` — 一定間隔でランダムな方向へ歩く
3. **乱数は必ず `*rand.Rand` を注入する**（グローバル `rand` を使わない）
   → シード固定でテストが再現可能になる
4. NPC を3体置き、それぞれ違う `Behavior` を持たせる

### 動作確認

```sh
make run
```

- [ ] 1体はその場に立ち、2体はうろうろ歩き回る
- [ ] 歩いている NPC のアニメーションと向きが正しい

### ポイント

- `Behavior` は W9 の敵AI（`battle.AI`）とまったく同じ構造。**同じ問題には同じ形が効く**
- 乱数の注入は「テスト可能性のための設計」の代表例

---

## Day 15 — 今週のリファクタリングと整理

| 項目 | 内容 |
| --- | --- |
| 目安 | 40分 |
| 触るファイル | 今週触った全ファイル |

### やること

1. `entity` パッケージのファイル構成を見直す（1ファイル200行を超えていたら分割）
2. `State` と `Behavior` の doc コメントを充実させる(**なぜその設計か**を書く)
3. 状態遷移図を `docs/03-architecture.md` に追記する（Mermaid の `stateDiagram-v2` が便利）
4. `make check` を通す

### 動作確認

```sh
make run
```

- [ ] 先週から今週までの機能がすべて壊れていない

### 今週の振り返り（ログに書く）

- State パターン導入前後で `Player.Update` は何行から何行になったか
- 「これも状態にできるな」と思ったものは何か（メニュー？会話中？→ W6, W7 の伏線）
