# assets — 素材の仕様

現在このディレクトリにある素材は **Dungeon Tileset II v1.7**（作者: 0x72 / CC0）のみ。
ライセンスの詳細は [`CREDITS.md`](./CREDITS.md) を参照。

読み込みは Day 06 で `//go:embed` により実行ファイルへ埋め込む予定
（[`docs/plan/week-02.md`](../docs/plan/week-02.md)）。

---

## ディレクトリ構成

```
assets/
├── CREDITS.md
├── README.md                      このファイル
└── images/
    └── dungeon-tileset-ii/
        ├── 0x72_DungeonTilesetII_v1.7.png   全部入り統合シート（512x512）
        ├── atlas_floor-16x16.png            床・仕掛け（112x112、16x16グリッド）
        ├── atlas_walls_low-16x16.png        組み立て済み部屋の見本（192x64、16x16グリッド）
        ├── atlas_walls_high-16x32.png       壁パーツ集（384x128、不揃い）
        ├── tile_list_v1.7                  【公式】統合シート内の座標一覧（後述）
        ├── README-vendor.txt               【公式】配布元 README（グリッド仕様の原文）
        └── frames/                          統合シートを1コマずつ切り出し済み（370枚）
```

**`frames/` と `tile_list_v1.7` が最重要。** 自分でスプライトシートを切り出す必要がなく、
名前とピクセル座標が最初から確定している。Week 02 の実装では基本的に `frames/` の
個別ファイルをそのまま使えばよく、統合シートからの切り出しは学習目的（Day 07）以外では不要。

---

## 1. グリッドの基本ルール（公式READMEの要約）

配布元 `README-vendor.txt` に明記されている設計:

- **基本グリッドは 16x16px。** 床・低い壁（`atlas_walls_low-16x16.png`）はそのまま
  16x16 で配置できる
- **高い壁（`atlas_walls_high-16x32.png`）は 16x32px だが、配置先のグリッドは 16x16 のまま。**
  Y方向に **基本 +8px、一部 -2px** のオフセットをかけて基準グリッドに乗せる
  （どちらのオフセットかはタイルごとに見て判断する、との記載。差が10pxあれば
  見分けられるとのこと）
- オートタイル設計は [3x3 minimal 方式](https://github.com/godotengine/godot-docs/issues/3316)
  に準拠

これは前回、別素材（RF Catacombs / Pixel Lands）を目視で解析したときに立てた
「壁は下端合わせで描く」という仮説と一致する。**このパックでは公式に明文化されている。**

---

## 2. `tile_list_v1.7` — 統合シート座標一覧

`0x72_DungeonTilesetII_v1.7.png`（512x512）内の各コマについて、
`名前 x y w h` の形式で1行ずつ記載されている。

```
big_demon_idle_anim_f0 16 428 32 36
wizzard_f_idle_anim_f0 128 132 16 28
coin_anim_f0 289 385 6 7
```

`frames/` 配下の同名ファイルは、この座標でクロップ済みのものと同一内容。
**統合シートを直接切り出すコードを書くより、`frames/` の個別ファイルを読む方が
実装がシンプルになる**（Week 02 ではこちらを使う）。

### フレームサイズの分布（370エントリ）

| サイズ | 件数 | 用途 |
| --- | --- | --- |
| 16x16 | 150 | 床・低い壁・アイテム・小物 |
| 16x28 | 90 | 人型キャラクター（wizzard, knight, elf, lizard, dwarf） |
| 16x23 | 64 | ゴブリン系・小型モンスター |
| 32x36 | 24 | 大型モンスター（big_demon, big_zombie） |
| その他 | 少数 | 剣・矢・コイン・ハートなど不定形アイテム |

**キャラクター・モンスターは看板サイズが16x16ではない。** 見た目の縦幅が
キャラごとに違うため、描画時は「足元（下端）をタイル境界に合わせる」のが正しい。
左上合わせで描くと頭が切れたり浮いたりする。

---

## 3. 命名規則

### キャラクター・モンスター

```
<名前>_<状態>_anim_f<N>.png
```

| 要素 | 例 | 意味 |
| --- | --- | --- |
| 状態 | `idle` `run` `hit` | アニメーションの種類。`hit`は1コマのみの個体が多い |
| `_f` / `_m` | `elf_f` `elf_m` | 性別バリエーション（該当キャラのみ） |
| `_anim_f<N>` | `_f0`〜`_f3` | フレーム番号。基本4コマ |

一部（`necromancer` `zombie` `slug` `swampy` `muddy` など）は状態区分がなく
`<名前>_anim_f<N>.png` のみ。

### プレイヤー候補になりうる人型キャラクター（`_f`/`_m` があるもの）

`dwarf` `elf` `knight` `lizard` `wizzard` の5種 × 男女2種 = 10種。
いずれも idle/run 4コマ + hit 1コマの構成で統一されている。
**Week 02〜03（アニメーション・State パターン）の題材として最も扱いやすい。**

### そのほかのモンスター

`angel` `big_demon` `big_zombie` `chort` `doc`（Halloween限定キャラ）
`goblin` `ice_zombie` `imp` `masked_orc` `muddy` `necromancer` `ogre`
`orc_shaman` `orc_warrior` `pumpkin_dude`（Halloween限定）`skelet` `slug`
`swampy` `tiny_slug` `tiny_zombie` `wogol` `zombie`

### タイル・小物

- `floor_1`〜`floor_8`: 床のバリエーション（ひび割れ等）
- `wall_*`: 壁の各パーツ（`atlas_walls_high` の切り出しと同一）
- `wall_banner_*`: 壁掛けの旗（4色）
- `wall_fountain_*`: 噴水（3コマアニメ + 台座）
- `chest_empty_open` / `chest_full_open` / `chest_mimic_open`: 宝箱（開閉3コマ）
- `door_*` / `doors_leaf_*`: 扉（開閉2種 + 枠3パーツ）
- `lever_left` / `lever_right`: レバー
- `floor_spikes`: 床のトゲ（アニメあり）
- `floor_ladder` / `floor_stairs` / `hole`: 移動関連ギミック
- `button_*_up` / `button_*_down`: スイッチ（赤・青、押下状態あり）
- `bomb_f0`〜`f2`: 爆弾の点火アニメ
- `flask_*` / `flask_big_*`: ポーション（4色 × 2サイズ）
- `coin_anim_f0`〜`f3`: コインの回転アニメ
- `ui_heart_full` / `ui_heart_half` / `ui_heart_empty`: HP表示アイコン
- `weapon_*`: 武器アイコン（20種以上。ドロップ表示・装備アイコンに使える）
- `column` / `column_wall` / `crate` / `skull`: 静的な装飾オブジェクト

---

## 4. 素材を追加するときのルール

1. `assets/images/` 配下に、パックごとのディレクトリを作って置く
2. **ライセンスを確認してから取り込む。** 確認前の素材は取り込まない
   （実際に別パック2種は再配布規約の懸念があり、ユーザー判断で見送った経緯がある）
3. `CREDITS.md` にその場で追記する（ライセンス全文または要点を引用する）
4. 配布元にREADMEや座標リストが同梱されていれば、リネームせずそのまま置く
   （このパックの `tile_list_v1.7` のように、後から座標を再解析する手間が省ける）
