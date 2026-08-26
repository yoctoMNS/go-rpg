# Week 11 — セーブ・オーディオ・イベントバス

> **今週のゴール**: **セーブして終了し、再起動して続きから遊べる**。BGM と SE が鳴る。

| | |
| --- | --- |
| 導入するパターン | **Observer / Event Bus** / DTO による永続化 |
| 追加するパッケージ | `save` `audio` `event` |

---

## Day 51 — セーブデータの設計

| 項目 | 内容 |
| --- | --- |
| 目安 | 55分 |
| 触るファイル | `internal/save/data.go` + テスト |
| 学ぶこと | **永続化用の型を分ける（DTO）** / 前方互換性 |

### やること

1. **セーブ用の型を新しく作る**。ゲーム内の型をそのまま JSON にしない
   ```go
   // SaveData はセーブファイルの内容。
   // ゲーム内の型を直接シリアライズすると、内部リファクタリングのたびに
   // 過去のセーブデータが読めなくなる。必ず変換用の型を挟む。
   type SaveData struct {
       Version   int                  `json:"version"`
       SavedAt   time.Time            `json:"savedAt"`
       PlayTime  int                  `json:"playTimeFrames"`
       MapID     string               `json:"mapId"`
       PlayerPos Position             `json:"playerPos"`
       Party     []CharacterSave      `json:"party"`
       Inventory []ItemSave           `json:"inventory"`
       Gold      int                  `json:"gold"`
       Flags     map[string]bool      `json:"flags"`
   }
   ```
2. **`Version` を必ず入れる。** 将来のマイグレーションの唯一の足がかり
3. ゲーム状態 → `SaveData` → ゲーム状態 の変換関数を書く
4. **往復テスト（ラウンドトリップ）を書く**
   ```go
   // ToSaveData → JSON → FromSaveData で、元と同じ状態に戻ること
   ```

### 動作確認

```sh
make test
```

- [ ] ラウンドトリップテストが通る
- [ ] JSON を目で見て、内容が理解できる（デバッグしやすさは重要な設計要件）

---

## Day 52 — セーブとロードの実装

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/save/store.go` |
| 学ぶこと | 安全なファイル書き込み / OS ごとの保存先 |

### やること

1. 保存先を決める（`os.UserConfigDir()` + `go-rpg/saves/`）
   - Windows / macOS / Linux で適切な場所が返る
2. **アトミックに書く**（これを怠るとセーブ中のクラッシュでデータが消える）
   ```go
   // 1. 一時ファイルに書く
   // 2. Sync() でディスクに確実に書き出す
   // 3. os.Rename で本来のパスに移動（同一ファイルシステム上では原子的）
   ```
3. 3スロット対応にする
4. ロード時のエラーハンドリング
   - ファイルが無い → 「データなし」として正常に扱う
   - JSON が壊れている → error。**ゲームを落とさず、UI にメッセージを出す**
5. 破損ファイルを読ませるテストを書く

### 動作確認

```sh
make run
```

- [ ] メニューの「セーブ」でスロットに保存できる
- [ ] ゲームを終了して再起動し、タイトルの「つづきから」で復帰できる
- [ ] 位置・向き・HP・所持金・アイテムがすべて復元される
- [ ] セーブファイルを手で壊してもゲームが落ちず、エラーメッセージが出る

---

## Day 53 — イベントバスを導入する

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/event/bus.go` + テスト |
| 学ぶパターン | **Observer** |

### やること

1. **まず「痛み」を確認する。** 「敵を倒した」ときに反応したいものを数える
   - 経験値付与 / SE 再生 / 撃破演出 / 討伐数カウント / 実績
   - これを呼び出し元が全部知っている状態は依存の爆発
2. `event.Bus` を作る
   ```go
   type Event interface{ isEvent() }

   type EnemyDefeated struct{ EnemyID string; Exp, Gold int }
   type PlayerDamaged struct{ Amount int }
   type LevelUp       struct{ Name string; NewLevel int }

   func (b *Bus) Subscribe(h Handler) (unsubscribe func())
   func (b *Bus) Publish(e Event)
   ```
3. **`Publish` 中の `Subscribe`/`Unsubscribe` で壊れないようにする**
   - ハンドラのスライスをコピーしてから回す。ここは実務でも定番のバグ
4. **同期実行にする**（goroutine を使わない）。ゲームループでは実行順の予測可能性が最優先
5. テスト: 購読 → 発行 → 呼ばれる / 解除 → 呼ばれない / 発行中の解除

### 動作確認

```sh
make run
```

- [ ] 敵撃破時に、経験値付与・SE・演出がすべて動く
- [ ] 撃破処理を呼ぶコードが、それらを一切知らない状態になっている

### ポイント

- **使いすぎに注意。** 「誰がいつ反応するか」が追えなくなると、デバッグ不能になる
- **判断基準**: 1対多で、かつ呼ぶ側が相手を知りたくない場合だけ。それ以外は直接呼ぶ

---

## Day 54 — BGM と SE

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/audio/` `assets/audio/` |
| 学ぶこと | `ebiten/v2/audio` / リソース管理 |

### やること

1. `audio.Context` は **アプリ全体で1つだけ**作る（複数作ると再生できない）
2. `audio.Player` を作る
   ```go
   func (p *Player) PlayBGM(id string)      // 同じ曲なら再開しない
   func (p *Player) StopBGM()
   func (p *Player) FadeOutBGM(frames int)
   func (p *Player) PlaySE(id string)       // 多重再生可
   ```
3. **BGM はストリーミング、SE は事前デコードして使い回す**
   - SE を毎回デコードすると再生が遅れる
4. イベントバスを購読して SE を鳴らす（Day 53 の成果をここで使う）
5. 音量設定を持たせる（後で設定画面から変えられるように）
6. **素材のライセンスを `assets/CREDITS.md` に記録する**

### 動作確認

```sh
make run
```

- [ ] タイトル / フィールド / バトルで BGM が切り替わる
- [ ] シーンが変わっても同じ曲なら**曲が最初から鳴り直さない**
- [ ] 決定 / キャンセル / 攻撃 / 撃破で SE が鳴る
- [ ] SE が連続で鳴っても途切れない

---

## Day 55 — 今週のリファクタリング

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |

### やること

1. イベントバスを使いすぎていないか見直す
   - **直接呼べる関係が Publish 経由になっていたら、直接呼びに戻す**
2. セーブデータに含めるべきものが漏れていないか確認する
   - 実際にセーブ → 終了 → ロードして、失われる情報を探す
3. `audio` の初期化を `Deps` に統合する
4. `make check` を通す

### 動作確認

```sh
make run
```

- [ ] 「起動 → つづきから → 歩く → 戦闘 → セーブ → 終了 → 再起動」が完走する

### 今週の振り返り（ログに書く）

- セーブデータのバージョンを上げる必要が出たとき、どう対応するか手順を書いておく
- イベントバスで「追いづらい」と感じた場面はあったか
