# go-rpg

Go + [Ebitengine](https://ebitengine.org/) で 2D RPG の基盤を作るプロジェクト。

**1日30〜60分 × 週5日 × 12週（全60日）で、RPGとして遊べる基盤を完成させる**ことを目標にしている。

---

## いまの状態

| | |
| --- | --- |
| 進捗 | **Day 01 / 60 完了** |
| 動くもの | ウィンドウが開き、TPS/FPS とフレーム数が表示される |
| 素材 | 配置済み（Dungeon Tileset II v1.7 / 0x72 / CC0）。仕様は [`assets/README.md`](assets/README.md) |

進捗の詳細は [`docs/PROGRESS.md`](docs/PROGRESS.md) を参照。

---

## 動かす

### 必要なもの

- Go 1.25 以上
- Linux の場合のみ、Ebitengine のビルドに X11 / OpenGL の開発ヘッダが必要

  ```sh
  sudo apt-get install -y libx11-dev libxrandr-dev libxcursor-dev \
      libxinerama-dev libxi-dev libxxf86vm-dev libgl1-mesa-dev
  ```

  macOS / Windows では追加インストールは不要。

### 起動

```sh
make run
```

### そのほかのコマンド

```sh
make help    # ターゲット一覧
make check   # fmt + vet + test（コミット前に必ず通す）
make test    # テスト（-race -shuffle=on）
make build   # bin/rpg を生成
make log     # 今日の作業ログを docs/logs/ に作成
```

---

## ドキュメント

| ファイル | 内容 |
| --- | --- |
| [`docs/00-roadmap.md`](docs/00-roadmap.md) | 3ヶ月・12週・60日の全体計画とパターン導入計画 |
| [`docs/01-rules.md`](docs/01-rules.md) | 1日のサイクル、記録・Git・テスト・リファクタのルール |
| [`docs/02-coding-standards.md`](docs/02-coding-standards.md) | コーディング規約（ゲーム開発特有の項目を含む） |
| [`docs/03-architecture.md`](docs/03-architecture.md) | パッケージ構成、レイヤ設計、デザインパターン適用計画 |
| [`docs/plan/`](docs/plan/) | 週ごとの日次計画（Day 01〜60） |
| [`docs/logs/`](docs/logs/) | 日々の作業ログ |
| [`docs/PROGRESS.md`](docs/PROGRESS.md) | 進捗チェックリスト |

---

## ディレクトリ構成

```
go-rpg/
├── cmd/rpg/          # エントリポイント（起動のみ。ロジックを書かない）
├── internal/
│   ├── config/       # 起動時設定
│   └── game/         # ebiten.Game 実装（ゲームループ）
├── assets/               # 画像・音・データ（Day 06 で //go:embed する）
│   ├── README.md         #   ★ タイルサイズ・命名規則の仕様
│   ├── CREDITS.md        #   ★ 出典とライセンス
│   └── images/
│       └── dungeon-tileset-ii/  # キャラ・モンスター・タイル一式（CC0）
├── docs/                 # 計画書とログ
└── scripts/              # 開発補助スクリプト
```

最終的な構成（12週後の姿）は [`docs/03-architecture.md`](docs/03-architecture.md) を参照。

---

## 設計の方針

**動くものを作る → 痛みが出る → パターンで解決する**、の順で進める。
最初から完璧な設計を狙わないことで、「なぜそのパターンが必要か」を体で理解することを重視している。

| 週 | 出てくる痛み | 導入するパターン |
| --- | --- | --- |
| W1 | 入力処理がロジックに散らばる | Input Abstraction (Adapter) |
| W3 | キャラの状態分岐が `if` の山になる | State / Strategy |
| W4 | 同じタイル画像でメモリを食う | Flyweight |
| W6 | `Game` が全状態を抱えて肥大化する | State (Scene) + Scene Stack |
| W7 | 会話の分岐がベタ書きになる | Command / Composite |
| W8 | バトル進行が巨大な switch になる | 有限状態機械 (FSM) |
| W10 | コンテンツ追加のたび再ビルドする | データ駆動 + Repository + Factory |
| W11 | 1つの出来事に多数が反応したがる | Observer / Event Bus |
| W12 | 構造体が肥大化する | Component 指向 |

---

## ライセンス

コードは [LICENSE](LICENSE) を参照。

アセット（Dungeon Tileset II / 0x72）は **CC0**（パブリックドメイン相当）。
詳細と引用元は [`assets/CREDITS.md`](assets/CREDITS.md) を参照。
