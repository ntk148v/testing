package main

import (
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func main() {
	// Configure Zap (console encoder)
	config := zap.NewDevelopmentConfig() // human-readable output
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Sugar logger is easier, but we will use the structured logger like Zerolog
	sugar := logger.Sugar()

	// Register handlers, passing logger
	http.HandleFunc("/", handleIndex(logger))
	http.HandleFunc("/log", handleLog(logger))
	http.HandleFunc("/data", handleData(logger))
	http.HandleFunc("/error", handleError(logger))

	sugar.Infof("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleIndex(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Accessing index page",
			zap.String("method", r.Method),
		)

		fmt.Fprintln(w, "Welcome to the Go Application!")
	}
}

func handleLog(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.With(
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		switch r.Method {
		case http.MethodGet:
			l.Info("Handled GET request on /log")
			fmt.Fprintln(w, "Received a GET request at /log.")
		case http.MethodPost:
			l.Info("Handled POST request on /log")
			fmt.Fprintln(w, "Received a POST request at /log.")
		default:
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		}
	}
}

func handleData(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Data endpoint hit",
			zap.String("method", r.Method),
			zap.String("endpoint", "/data"),
			zap.String("request_id", fmt.Sprintf("%d", os.Getpid())),
		)

		fmt.Fprintf(w, "This is the data endpoint. Method used: %s\n", r.Method)
	}
}

func handleError(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Error("Error endpoint accessed",
			zap.String("method", r.Method),
			zap.String("endpoint", "/error"),
			zap.String("error_id", "error123"),
		)

		http.Error(w, "You have reached the error endpoint", http.StatusInternalServerError)
	}
}
