package main

import (
	"errors"

	log "github.com/inconshreveable/log15"
)

func main() {
	log.Info("Hello from Log15", "name", "John", "age", 20)

	logger := log.New("env", "prod", "go_version", "1.20")

	logger.Info("Hello from Log15", "name", "John", "age", 20)
	logger.Error("Something unexpected happened", log.Ctx{
		"error": errors.New("an error"),
	})
}
