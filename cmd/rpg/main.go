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
