# Week 08 — バトルの骨格（有限状態機械）

> **今週のゴール**: フィールドを歩いているとエンカウントし、**バトル画面に切り替わって
> ターンが回り、勝敗が決まってフィールドに戻る**（行動は「たたかう」のみ）。

| | |
| --- | --- |
| 導入するパターン | **FSM（有限状態機械）** / ドメインとプレゼンテーションの分離 |
| 追加するパッケージ | `battle` `scene/battle` |

---

## Day 36 — 戦闘ドメインを ebiten 抜きで作る

| 項目 | 内容 |
| --- | --- |
| 目安 | 55分 |
| 触るファイル | `internal/battle/combatant.go` `internal/battle/damage.go` + テスト |
| 学ぶこと | **Domain 層に描画ライブラリを持ち込まない** |

### やること

1. `battle` パッケージに **`ebiten` を一切 import しない**と決める
2. 型を定義する
   ```go
   // Stats は戦闘に使う能力値。マスタデータ由来で、戦闘中に変わらない。
   type Stats struct {
       MaxHP, MaxMP     int
       Attack, Defense  int
       Speed            int
   }

   // Combatant は戦闘に参加する1体。HP/MP は戦闘中に変化する。
   type Combatant struct {
       Name  string
       Stats Stats
       HP    int
       MP    int
       Alive bool
   }
   ```
3. ダメージ計算を **純粋関数**で書く
   ```go
   // Damage は攻撃力と防御力からダメージ量を返す。最低保証は 1。
   func Damage(attack, defense int) int
   ```
4. **テーブル駆動テストを厚く書く**（通常 / 防御が上回る / 0 / 極端に大きい値）

### 動作確認

```sh
make test
```

- [ ] `internal/battle` のテストが通る
- [ ] `go list -deps ./internal/battle | grep ebiten` が **何も出さない**（依存が無い証拠）

### ポイント

- ドメインが独立していると、**GPU 無しの CI でテストが回る**、AI に自動対戦させられる、
  バランス調整をスクリプトで回せる、といった利点がすべて手に入る

---

## Day 37 — 戦闘の状態機械を作る

| 項目 | 内容 |
| --- | --- |
| 目安 | 55分 |
| 触るファイル | `internal/battle/battle.go` `internal/battle/phase.go` + テスト |
| 学ぶパターン | **FSM** |

### やること

1. フェーズを型付き定数で定義する
   ```go
   type Phase int

   const (
       PhaseStart      Phase = iota // 開始演出
       PhaseCommand                 // プレイヤーのコマンド入力待ち
       PhaseAction                  // 行動の実行
       PhaseTurnEnd                 // ターン終了処理
       PhaseWin
       PhaseLose
   )
   ```
2. `Battle` 型が現在のフェーズと参加者を持つ
3. `Battle.Step()` で1ステップ進める（**フレームではなく論理ステップ**）
   - 演出の待ち時間はプレゼンテーション層の責務。ドメインは「次に何が起きるか」だけ決める
4. **遷移表をテストする**
   - 敵を全滅させたら `PhaseWin`
   - 味方が全滅したら `PhaseLose`
   - 全員行動したら `PhaseTurnEnd` → `PhaseCommand`

### 動作確認

```sh
make test
```

- [ ] 「攻撃 → 敵撃破 → 勝利」の流れがテストで再現できる
- [ ] 画面はまだ無くてよい（ドメインだけで戦闘が完結している）

---

## Day 38 — バトル画面を作る

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/scene/battle/battle.go` |
| 学ぶこと | ドメインを画面に載せる |

### やること

1. `BattleScene` を作り、`battle.Battle` を1つ持たせる
2. 描画する
   - 背景 + 敵グラフィック
   - 味方のステータスウィンドウ（名前 / HP / MP）
   - コマンドウィンドウ（「たたかう」のみ）
3. `battle.Phase` に応じて画面表示を切り替える
4. **ドメインの `Step()` を呼ぶタイミングを、演出の区切りに合わせる**

### 動作確認

```sh
make run
```

- [ ] バトル画面が表示される（一時的に起動直後にバトルへ飛ばして確認する）
- [ ] 「たたかう」でダメージが出て、敵の HP が減る
- [ ] 敵を倒すと勝利メッセージが出る

---

## Day 39 — エンカウントとシーン遷移

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/scene/field/encounter.go` |
| 学ぶこと | シーン間のデータ受け渡し / 乱数の扱い |

### やること

1. `Encounter` を作る（歩数ベースの判定にする）
   ```go
   // 歩数ベースにする理由: フレームベースだと「その場で足踏み」でも遭遇して不自然。
   // また、歩数なら「N歩以内は必ず出ない」という保証（初期猶予）が作りやすい。
   type Encounter struct {
       stepsUntilNext int
       minSteps       int
       maxSteps       int
       rng            *rand.Rand
   }
   ```
2. 1タイル歩くごとにカウントを減らし、0 でエンカウント
3. `Manager.Push(NewBattleScene(...))` でバトルを**重ねる**（Replace ではない）
   - Push にすると、バトル終了時に `Pop` するだけでフィールドの状態が完全に戻る
4. バトルの結果（勝敗、獲得経験値）をフィールドに返す仕組みを作る
   - コールバック関数を渡すのがシンプル: `NewBattleScene(enemies, func(result Result) {...})`

### 動作確認

```sh
make run
```

- [ ] フィールドを歩いているとランダムにバトルが始まる
- [ ] バトルに勝つと、**元いた場所・向きのまま**フィールドに戻る
- [ ] その場で足踏みしてもエンカウントしない

---

## Day 40 — 今週のリファクタリング

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |

### やること

1. `battle` パッケージが本当に `ebiten` から独立しているか再確認する
   ```sh
   go list -deps ./internal/battle | grep -c ebiten   # 0 であること
   ```
2. `BattleScene` の `Draw` が長くなっていたら、描画パーツごとに分割する
   （`drawEnemies` / `drawStatus` / `drawCommandWindow`）
3. ダメージ計算式のテストケースを追加する（レベル差が極端な場合など）
4. `make check` を通す

### 今週の振り返り（ログに書く）

- ドメイン（`battle`）と画面（`scene/battle`）の境界で迷った箇所はどこか
- Push / Replace のどちらを使うべきか、判断基準を自分の言葉で書く
