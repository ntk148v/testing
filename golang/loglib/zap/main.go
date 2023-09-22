package main

import (
	"os"

	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	if os.Getenv("APP_ENV") == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}

	zap.ReplaceGlobals(logger)

	zap.L().Info("Hello from zap!")
}
