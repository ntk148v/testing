package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
)

func main() {
	logger := log.WithFields(log.Fields{
		"app": "myapp",
		"env": "prod",
	})

	logger.Info("Hello from Apex logger")

	// add logger
	stdout := logfmt.New(os.Stdout)

	log.SetHandler(stdout)
	logger = log.WithFields(log.Fields{
		"app": "myapp",
		"env": "prod",
	})

	logger.Info("Hello from Apex logger")
	// apex support tracing
	// defer l.Trace("fetching random quote").Stop(&err)
}
