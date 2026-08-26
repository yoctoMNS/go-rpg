# Week 10 — データ駆動と成長要素

> **今週のゴール**: 敵・アイテム・スキルが **JSON で定義され、コードを書き換えずに追加できる**。
> 経験値でレベルが上がり、アイテムを持てる。

| | |
| --- | --- |
| 導入するパターン | **Repository** / **Factory** / データ駆動設計 |
| 追加するパッケージ | `gamedata` `inventory` |

---

## Day 46 — マスタデータを JSON に出す

| 項目 | 内容 |
| --- | --- |
| 目安 | 55分 |
| 触るファイル | `internal/gamedata/` `assets/data/enemies.json` |
| 学ぶパターン | **Repository** |

### やること

1. `assets/data/enemies.json` を書く
   ```json
   [
     {
       "id": "slime",
       "name": "スライム",
       "stats": { "maxHP": 12, "attack": 6, "defense": 3, "speed": 4 },
       "ai": "aggressive",
       "exp": 3,
       "gold": 5,
       "sprite": "enemies/slime.png"
     }
   ]
   ```
2. `gamedata.Repository` を作る
   ```go
   // Repository はマスタデータの読み取り専用アクセスを提供する。
   type Repository[T any] interface {
       Find(id string) (T, error)
       All() []T
   }
   ```
   - Go のジェネリクスで敵・アイテム・スキルを共通化できる
3. **起動時に一括ロードして検証する**
   - ID の重複、存在しないスプライトパス、未知の `ai` 名 → **起動時に error**
   - 「実行中に初めて壊れる」のが最悪。データの不正は起動時に全部見つける
4. バリデーションのテストを書く

### 動作確認

```sh
make run
```

- [ ] 敵のステータスが JSON から読まれている（JSON の HP を変えると反映される）
- [ ] JSON をわざと壊すと、起動時に**分かりやすいエラーメッセージ**で落ちる

### ポイント

- **`EnemyDef`（定義）と `Combatant`（実体）を必ず別の型にする**
- 混ぜると「マスタデータの HP が減っていて、2回目の戦闘で敵が瀕死」という致命的バグになる

---

## Day 47 — Factory で実体を作る

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/gamedata/factory.go` |
| 学ぶパターン | **Factory** |

### やること

1. `NewCombatant(def EnemyDef) *battle.Combatant` を書く
   - **`Stats` は値コピー**する（定義を共有しない）
2. `ai` 文字列から `battle.AI` を返す関数を書く
   ```go
   func newAI(name string, rng *rand.Rand) (battle.AI, error) {
       switch name {
       case "aggressive": return battle.AggressiveAI{}, nil
       case "cautious":   return battle.CautiousAI{}, nil
       case "random":     return battle.RandomAI{Rng: rng}, nil
       default:           return nil, fmt.Errorf("unknown ai %q", name)
       }
   }
   ```
   - **この switch は許容する**。文字列から型への変換は必ずどこかで1回必要で、
     ここ1か所に閉じているなら健全（境界の switch）
3. 敵グループ（「スライム x2 + こうもり」）を JSON で定義できるようにする

### 動作確認

```sh
make run
```

- [ ] JSON に新しい敵を追加すると、**Goのコードを1行も書かずに**戦闘に出てくる
- [ ] 敵グループが JSON どおりに編成される

### ポイント

- **これがデータ駆動の到達点**。ここから先、コンテンツ追加はデータ編集だけで済む

---

## Day 48 — 経験値とレベルアップ

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/party/growth.go` `assets/data/levels.json` + テスト |
| 学ぶこと | 成長曲線をデータで持つ |

### やること

1. レベルごとの必要経験値と能力値を JSON で持つ
   - **計算式をコードに埋めない**。バランス調整のたびに再ビルドするのは苦痛
2. `AddExp(c *Character, exp int) []LevelUpResult` を書く
   - **一度に複数レベル上がる場合に対応する**（ボス撃破でよくある）
3. テストを厚く書く（1レベルアップ、複数レベル同時、最大レベルで打ち止め）

### 動作確認

```sh
make run
```

- [ ] 戦闘後に経験値が入り、必要量に達するとレベルアップメッセージが出る
- [ ] 大量の経験値で 2 レベル以上一気に上がる

---

## Day 49 — アイテムとインベントリ

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/inventory/` `assets/data/items.json` + テスト |
| 学ぶこと | コレクションの設計 / 上限の扱い |

### やること

1. `items.json` を作る（消耗品 / 武器 / 防具）
2. `inventory.Inventory` を作る
   ```go
   func (inv *Inventory) Add(itemID string, count int) error   // 所持上限を超えたら error
   func (inv *Inventory) Remove(itemID string, count int) error
   func (inv *Inventory) Count(itemID string) int
   func (inv *Inventory) Items() []Slot                        // 表示用（ソート済み）
   ```
3. **所持数上限・所持種類上限を最初から実装する**（後から入れるとUI全体に影響する）
4. 戦闘とフィールド両方から「どうぐ」を使えるようにする
5. テスト: 上限超過、存在しないアイテムの削除、個数0でスロットが消えること

### 動作確認

```sh
make run
```

- [ ] 戦闘で「どうぐ」から回復薬を使うと HP が回復し、所持数が減る
- [ ] 所持数が 0 になるとリストから消える

---

## Day 50 — メニュー画面

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/scene/menu/` |
| 学ぶこと | Scene Stack の実用 / リスト UI の共通化 |

### やること

1. フィールドで `ActionMenu`（Esc / Enter）を押すとメニューを **`Push`** する
   - Push なので、後ろにフィールドが見えたまま重なる
2. メニュー項目: 「どうぐ」「そうび」「ステータス」「セーブ」（セーブは W11 で中身を入れる）
3. **リスト UI を `ui.List` として共通化する**
   - コマンドウィンドウ、アイテム一覧、スキル一覧で使い回す
   - カーソル移動、スクロール、ページング、決定/キャンセルを1か所に持つ
4. `ui.List` のカーソル移動ロジックにテストを書く（端でのループ、スクロール境界）

### 動作確認

```sh
make run
```

- [ ] フィールドで Esc を押すとメニューが開き、後ろにフィールドが見える
- [ ] カーソルが動き、決定で下位メニューに入り、キャンセルで戻る
- [ ] メニューを閉じると、**元の位置・向きのまま**歩ける

### 今週の振り返り（ログに書く）

- 新しい敵1体を追加するのに必要な作業は何か（JSON に数行、で済んでいるか）
- `ui.List` を共通化したことで、どれだけ重複が消えたか
