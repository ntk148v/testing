package main

import (
	"log/slog"
	"os"
)

func main() {
	slog.Debug("Debug message")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Debug("Debug message")
	logger.Info("Info message")
}
