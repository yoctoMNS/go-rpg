# Week 09 — バトルの中身（Command と Strategy）

> **今週のゴール**: 「たたかう / まほう / どうぐ / にげる」が選べ、**敵が状況に応じて
> 行動を変える**。素早さ順に行動が回る。

| | |
| --- | --- |
| 導入するパターン | **Command**（行動）/ **Strategy**（敵AI） |

---

## Day 41 — 行動を Command にする

| 項目 | 内容 |
| --- | --- |
| 目安 | 55分 |
| 触るファイル | `internal/battle/action.go` + テスト |
| 学ぶパターン | **Command** |

### やること

1. `Action` インターフェースを定義する
   ```go
   // Action は戦闘中の1回の行動。
   // Perform は行動を実行し、画面に出すためのログを返す。
   type Action interface {
       Actor() *Combatant
       Perform(f *Field) []LogEntry
   }
   ```
2. 実装を作る
   - `AttackAction{From, To *Combatant}`
   - `DefendAction{From *Combatant}` — このターンの被ダメージ半減
   - `EscapeAction{From *Combatant}` — 素早さ差で成功率が決まる
3. **`LogEntry` を返すのが重要**。ドメインは「何が起きたか」を返すだけで、
   それをどう演出するかは画面側が決める
   ```go
   type LogEntry struct {
       Kind    LogKind // Damage / Heal / Miss / Message
       Target  *Combatant
       Amount  int
       Message string
   }
   ```
4. 各 Action のテストを書く（防御中のダメージが半分になること、など）

### 動作確認

```sh
make run
```

- [ ] 「たたかう」「ぼうぎょ」「にげる」がコマンドウィンドウに並ぶ
- [ ] 「ぼうぎょ」を選ぶと、そのターンの被ダメージが減る
- [ ] 「にげる」で戦闘から離脱できる（失敗することもある）

---

## Day 42 — ターン順（素早さ）とアクションキュー

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/battle/turn.go` + テスト |
| 学ぶこと | 安定ソート / 実行順序の決定 |

### やること

1. 全員の `Action` を集めてキューに積む
2. **素早さの降順にソートする**
   - `sort.SliceStable` を使う（同速のとき順序がランダムに変わると、リプレイが再現しない）
   - 同速の場合の順序ルールを決めてコメントに書く（例: 味方優先）
3. キューから1つずつ取り出して `Perform` する
4. **行動前に「まだ生きているか」を必ず確認する**（先に倒された敵が行動しないように）
5. ターン順のテストを書く（素早さ同値、途中で死亡、など）

### 動作確認

```sh
make run
```

- [ ] 素早さの高いキャラから行動する
- [ ] 先に倒した敵が行動してこない

---

## Day 43 — 敵AIを Strategy にする

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/battle/ai.go` + テスト |
| 学ぶパターン | **Strategy** |

### やること

1. `AI` インターフェースを定義する
   ```go
   // AI は敵1体のこのターンの行動を決める。
   type AI interface {
       Decide(self *Combatant, f *Field) Action
   }
   ```
2. 実装を3つ作る
   - `AggressiveAI` — 常に攻撃。HP が一番低い相手を狙う
   - `RandomAI` — ランダムな相手を攻撃（`*rand.Rand` を注入）
   - `CautiousAI` — 自分の HP が 30% 以下なら防御、それ以外は攻撃
3. **テストを書く**（`CautiousAI` が HP 30% 以下で `DefendAction` を返すこと）
   - 乱数を固定シードにすれば、`RandomAI` も決定的にテストできる

### 動作確認

```sh
make run
```

- [ ] 敵の種類によって行動が変わる
- [ ] `CautiousAI` の敵は瀕死になると防御してくる

### ポイント

- W3 の `entity.Behavior` とまったく同じ構造。**同じ問題には同じ形が効く**という実感を持つ
- W10 では、敵 JSON に `"ai": "cautious"` と書くだけで挙動が決まるようにする

---

## Day 44 — 魔法とMP

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/battle/skill.go` + テスト |
| 学ぶこと | 効果の合成 / 「定義」と「実行」の分離 |

### やること

1. `Skill` を **定義（データ）** として作る
   ```go
   // SkillDef は魔法・特技の定義。マスタデータ。書き換えない。
   type SkillDef struct {
       ID       string
       Name     string
       Cost     int
       Target   TargetKind // 敵単体 / 敵全体 / 味方単体 / 味方全体
       Effects  []EffectDef
   }
   ```
2. `SkillAction` が `SkillDef` を受け取って実行する
   - MP 不足なら実行できない（`CanPerform` で事前チェック）
3. 効果を種類ごとに分ける（`EffectDamage` / `EffectHeal`）
4. 対象選択（敵単体 / 全体）を実装する
5. テスト: MP 消費、MP 不足時、全体攻撃で全員にダメージ

### 動作確認

```sh
make run
```

- [ ] 「まほう」でスキル一覧が出る
- [ ] 実行すると MP が減り、効果が出る
- [ ] MP が足りないスキルは選べない（グレーアウト）

---

## Day 45 — 戦闘演出と今週のリファクタリング

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/scene/battle/` |

### やること

1. `LogEntry` を演出に変える
   - ダメージ数値をポップアップ表示（上に飛んで消える）
   - 被弾時に対象を白く点滅させる（`ColorScale` を使う）
   - 撃破時にフェードアウト
2. **演出は必ずプレゼンテーション層に置く。** `battle` パッケージには一切書かない
3. `BattleScene` のフェーズごとの `Update` を整理する
4. `make check` を通す

### 動作確認

```sh
make run
```

- [ ] ダメージ数値が飛び出して消える
- [ ] 攻撃を受けたキャラが点滅する
- [ ] 一連の戦闘が最後まで気持ちよく流れる

### 今週の振り返り（ログに書く）

- 新しい行動（例: 「アイテムを盗む」）を足すには何ファイル触ることになるか
- 1ファイルで済むなら Command パターンが効いている証拠
