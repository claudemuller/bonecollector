package main

import (
	"log/slog"
	"os"

	engine "bonecollector/pkg/engine"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	game, err := engine.New("test", 800, 600)
	if err != nil {
		slog.Error("engine", "error", err)
		return
	}
	defer game.Destroy()

	game.Run()
}
