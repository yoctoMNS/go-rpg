# CLAUDE.md

このリポジトリで作業するときの指針。

## プロジェクトの目的

Go + Ebitengine で 2D RPG の基盤を作る **学習プロジェクト**。
オーナーは「1日30〜60分 × 週5日 × 12週（全60日）」で自分の手で実装しながら進めている。

そのため、**先回りして大量のコードを書かない**こと。原則として **その日の Day のスコープだけ** を実装する。

## 必ず最初に読むもの

| ファイル | 内容 |
| --- | --- |
| `docs/00-roadmap.md` | 全体計画とパターン導入計画 |
| `docs/01-rules.md` | 1日のサイクル、記録・Git・テストのルール |
| `docs/02-coding-standards.md` | コーディング規約 |
| `docs/03-architecture.md` | パッケージ構成とレイヤ設計 |
| `docs/PROGRESS.md` | いまどの Day まで進んでいるか |
| `docs/plan/week-NN.md` | 今日やる作業の詳細 |
| `assets/README.md` | 素材のタイルサイズ・オフセット規則・命名規則（切り出し前に必読） |
| `assets/CREDITS.md` | 素材の出典・ライセンス |

## 作業の進め方

1. `docs/PROGRESS.md` で現在の Day を確認する
2. 対応する `docs/plan/week-NN.md` の当日分を読む
3. **その Day のスコープだけ**を実装する（次の Day を先取りしない）
4. `make check` を通す
5. `docs/logs/YYYY-MM-DD.md` を書く（`make log` で雛形を作る）
6. `docs/PROGRESS.md` のチェックを埋める

## コードを書くときの必須事項

- **`Update` は状態更新のみ、`Draw` は描画のみ。** `Draw` で状態を変えない
- **時間はフレーム数で数える。** `time.Now()` をゲームロジックに使わない
- **`Update` / `Draw` の中でメモリを確保しない。** バッファ・`DrawImageOptions` は使い回す
- **ドメイン層（`battle` / `entity` のロジック部）に `ebiten` を import しない**
- **グローバルな可変状態を作らない。** 依存はコンストラクタで注入する
- **公開する型・関数には doc コメント。** 「何を」ではなく「なぜ」を書く
- **インターフェースは使う側のパッケージに定義する**
- **痛みが出る前にパターンを入れない。** 計画書に書かれた週まで待つ

## アセット

- 素材は `assets/images/dungeon-tileset-ii/` に配置済み（Dungeon Tileset II v1.7 / 0x72 / CC0）
- **`frames/` 配下に1コマずつ切り出し済みのPNGがある。** 統合シートを自前で切り出す必要はない
- タイルは16x16、キャラは16x28前後（16x16グリッドの下端に接地）。詳細は `assets/README.md`
- 素材を追加したら `assets/CREDITS.md` にその場で追記する

## テスト

- テーブル駆動テスト、外部テストパッケージ（`package xxx_test`）、`t.Parallel()`
- 必ずテストする: ダメージ計算 / 当たり判定 / 座標変換 / 状態遷移 / セーブ入出力
- テストしない: `Draw` / 入力デバイスの直接読み取り

## コマンド

```sh
make run     # 起動
make check   # fmt + vet + test
make log     # 今日の作業ログを作成
make help    # 一覧
```

## Git

- ブランチ: `feat/dayNN-<内容>`
- コミット: Conventional Commits（`feat:` `fix:` `refactor:` `test:` `docs:` `chore:` `wip:`）
- 要約は日本語で書く

## 環境の注意

Linux で Ebitengine をビルドするには X11 / OpenGL の開発ヘッダが必要。

```sh
sudo apt-get install -y libx11-dev libxrandr-dev libxcursor-dev \
    libxinerama-dev libxi-dev libxxf86vm-dev libgl1-mesa-dev
```

ヘッダを入れられない環境（ヘッドレスコンテナ等）では、コンパイル確認だけなら
`GOOS=js GOARCH=wasm go build ./...` が使える（cgo を通らないため）。
