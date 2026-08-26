# Week 04 — タイルマップを表示する

> **今週のゴール**: JSON で定義したマップが画面に描かれ、その上をキャラが歩ける
> （まだ壁はすり抜ける）。

| | |
| --- | --- |
| 導入するパターン | **Flyweight** / データ駆動の第一歩 |
| 追加するパッケージ | `tilemap` |

---

## Day 16 — タイル定義を作る（Flyweight）

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/tilemap/tile.go` `internal/tilemap/tileset.go` |
| 学ぶパターン | **Flyweight** |

### やること

1. **共有される定義**と**個々の配置**を分ける
   ```go
   // TileDef は1種類のタイルの性質。種類ごとに1つだけ存在し、共有される（内在状態）。
   type TileDef struct {
       ID       int
       Name     string
       Passable bool // 通れるか
   }

   // Tileset は TileDef を ID で引ける集合 + タイル画像。
   type Tileset struct {
       sheet *sprite.Sheet
       defs  map[int]TileDef
   }
   ```
2. `Tileset.Def(id int) (TileDef, bool)` と `Tileset.Image(id int) *ebiten.Image` を作る
3. **画像はあらかじめ全部 `SubImage` して slice にキャッシュする**（毎フレーム切り出さない）

### 動作確認

```sh
make run
```

- [ ] タイル画像を ID 0〜5 まで横一列に並べて描画し、すべて正しく出る（確認用の一時コード）

### ポイント

- タイルが 1000 個あっても `TileDef` は種類数（10個程度）しか存在しない。これが Flyweight
- **配置データは `int` の slice だけ**。メモリもキャッシュ効率も段違いに良くなる

---

## Day 17 — マップデータを JSON で読み込む

| 項目 | 内容 |
| --- | --- |
| 目安 | 50分 |
| 触るファイル | `internal/tilemap/map.go` `assets/maps/town.json` |
| 学ぶこと | データとコードの分離 / `encoding/json` / バリデーション |

### やること

1. マップの JSON スキーマを自分で決める（**まずは Tiled を使わず、最小限を自作する**）
   ```json
   {
     "name": "town",
     "width": 20,
     "height": 15,
     "tileSize": 16,
     "layers": [
       { "name": "ground", "tiles": [1, 1, 1, ...] },
       { "name": "objects", "tiles": [0, 0, 4, ...] }
     ]
   }
   ```
2. `tilemap.Load(fs fs.FS, path string) (*Map, error)` を書く
3. **必ずバリデーションする**（`len(tiles) != width*height` なら error）
   - データの不正は起動時に落とす。ゲームループ中に発覚させない
4. `Map.TileAt(layer, col, row int) int` を書き、範囲外の扱いを決めてテストする

### 動作確認

```sh
make test
```

- [ ] 正常な JSON が読める
- [ ] タイル数が合わない JSON で `error` が返る（エラーメッセージにファイル名が入っている）

### ポイント

- **自作スキーマから始める理由**: Tiled の JSON は仕様が大きく、いきなり読むと本質が見えない。
  自分で決めた最小スキーマを一度通してから Tiled に移行する（Day 20）
- レイヤーを最初から複数対応にしておく（地面 / 装飾 / 上に被さるもの）

---

## Day 18 — マップを描画する

| 項目 | 内容 |
| --- | --- |
| 目安 | 45分 |
| 触るファイル | `internal/tilemap/renderer.go` |
| 学ぶこと | 描画の最適化 / 可視範囲カリング |

### やること

1. `Renderer` を作り、レイヤーを順に描く
2. **`DrawImageOptions` を使い回す**（毎タイルで `new` しない）
3. **画面に映る範囲だけ描く**（カリング）
   ```go
   // 画面外のタイルを描かない。マップが大きくなるほど効く。
   startCol := int(cam.X) / tileSize
   endCol   := (int(cam.X) + screenW) / tileSize + 1
   ```
4. プレイヤーをマップの上に重ねて描く

### 動作確認

```sh
make run
```

- [ ] タイルマップが画面いっぱいに描かれる
- [ ] その上をキャラが歩ける（まだ壁をすり抜ける）
- [ ] `TPS: 60.0` を維持している（カリングが効いている）

### ポイント

- タイルの継ぎ目に線が見えるときは、拡大倍率が整数でないか、`FilterLinear` になっているのが原因。
  ドット絵では **`ebiten.FilterNearest`** を使う

---

## Day 19 — マップエディタ代わりの確認機能

| 項目 | 内容 |
| --- | --- |
| 目安 | 40分 |
| 触るファイル | `internal/game/debug.go` |
| 学ぶこと | デバッグ機能への投資 |

### やること

1. F1 でデバッグ表示をトグルする
2. デバッグ表示の内容
   - タイルのグリッド線
   - プレイヤーのタイル座標 `(col, row)` と実座標
   - マウス位置のタイル ID
3. **デバッグ描画は `debug.go` に隔離する**（本番描画に混ぜない）

### 動作確認

```sh
make run
```

- [ ] F1 でグリッドと座標表示が出たり消えたりする
- [ ] マウスをタイルに乗せると、その ID が表示される

### ポイント

- **デバッグ機能は「作業時間の投資」**。ここで作ったグリッド表示は、W5 の当たり判定デバッグで
  何倍にもなって返ってくる
- 30分の実装が、この先の何時間ものデバッグを節約する

---

## Day 20 — Tiled 形式に対応する（または自作形式を確定する）

| 項目 | 内容 |
| --- | --- |
| 目安 | 55分 |
| 触るファイル | `internal/tilemap/tiled.go` |
| 学ぶこと | 外部フォーマットへの依存の切り離し |

### やること

1. [Tiled Map Editor](https://www.mapeditor.org/) をインストールし、簡単なマップを作って
   JSON でエクスポートする
2. **Tiled の JSON を直接ゲームで使わない。** 読み込んで **自分の `tilemap.Map` に変換する**
   ```go
   // tiled.go は「Tiled の形」を知っている唯一の場所にする。
   func LoadTiled(fs fs.FS, path string) (*Map, error)
   ```
   - これで、将来エディタを変えても影響が `tiled.go` 1ファイルに閉じる（**腐敗防止層**）
3. Tiled の GID は 1 始まり・0 は空白。ここで必ずハマるので、変換にテストを書く
4. 時間が足りなければ **自作形式のままでよい**。その判断もログに書く

### 動作確認

```sh
make run
```

- [ ] Tiled で作ったマップがゲーム内に表示される
- [ ] Tiled で編集して再エクスポートすると、ゲームに反映される

### 今週の振り返り（ログに書く）

- 外部フォーマットを直接使わず変換した理由を、自分の言葉で書けるか
- マップが大きくなったときの描画負荷はどうだったか
