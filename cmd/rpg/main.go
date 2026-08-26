// Command rpg is the entry point of the game.
//
// 責務はここでは「起動」だけに絞る。ゲームのロジックは internal 配下に置き、
// main は「設定を組み立てて実行する」ことのみを行う（Composition Root）。
package main

import (
	"log/slog"
	"os"

	"github.com/yoctoMNS/go-rpg/internal/config"
	"github.com/yoctoMNS/go-rpg/internal/game"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Default()

	if err := game.Run(cfg); err != nil {
		logger.Error("game exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}
