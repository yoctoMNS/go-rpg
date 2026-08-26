# go-rpg 開発用タスク集
#
# 「毎日30分」の開発では、コマンドを思い出す時間すら惜しい。
# よく使う操作は必ず make のターゲットとして 1 行で叩けるようにしておく。

GO      ?= go
PKG     := ./...
CMD     := ./cmd/rpg
BIN_DIR := bin
BIN     := $(BIN_DIR)/rpg

.DEFAULT_GOAL := help

## help: 利用可能なターゲット一覧を表示する
.PHONY: help
help:
	@grep -E '^## [a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed -e 's/^## //' | awk -F': ' '{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

## run: ゲームを起動する（毎日の動作確認はこれ）
.PHONY: run
run:
	$(GO) run $(CMD)

## build: 実行ファイルを bin/ に生成する
.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) $(CMD)

## test: テストを実行する
.PHONY: test
test:
	$(GO) test -race -shuffle=on $(PKG)

## cover: カバレッジを計測して HTML で開く
.PHONY: cover
cover:
	$(GO) test -coverprofile=coverage.out $(PKG)
	$(GO) tool cover -html=coverage.out

## vet: 静的解析をかける
.PHONY: vet
vet:
	$(GO) vet $(PKG)

## fmt: コードを整形する
.PHONY: fmt
fmt:
	$(GO) fmt $(PKG)

## tidy: go.mod / go.sum を整理する
.PHONY: tidy
tidy:
	$(GO) mod tidy

## check: コミット前の一括チェック（fmt + vet + test）
.PHONY: check
check: fmt vet test

## log: 今日の作業ログファイルを docs/logs/ に作成する
.PHONY: log
log:
	@scripts/new-log.sh

## clean: 生成物を削除する
.PHONY: clean
clean:
	rm -rf $(BIN_DIR) coverage.out
