package main

import (
	"fmt"
	"net/http"
	"os"

	"io"

	"github.com/rs/zerolog"
)

func main() {
	// Configure zerolog
	mw := io.MultiWriter(os.Stdout)

	logger := zerolog.New(mw).With().Timestamp().Logger()

	// Optional: pretty console output
	logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// We pass a copy of logger to handlers
	http.HandleFunc("/", handleIndex(logger))
	http.HandleFunc("/log", handleLog(logger))
	http.HandleFunc("/data", handleData(logger))
	http.HandleFunc("/error", handleError(logger))

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleIndex(logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info().
			Str("method", r.Method).
			Msg("Accessing index page")

		fmt.Fprintln(w, "Welcome to the Go Application!")
	}
}

func handleLog(logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path)

		switch r.Method {
		case http.MethodGet:
			l.Msg("Handled GET request on /log")
			fmt.Fprintln(w, "Received a GET request at /log.")
		case http.MethodPost:
			l.Msg("Handled POST request on /log")
			fmt.Fprintln(w, "Received a POST request at /log.")
		default:
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		}
	}
}

func handleData(logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info().
			Str("method", r.Method).
			Str("endpoint", "/data").
			Str("request_id", fmt.Sprintf("%d", os.Getpid())).
			Msg("Data endpoint hit")

		fmt.Fprintf(w, "This is the data endpoint. Method used: %s\n", r.Method)
	}
}

func handleError(logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Error().
			Str("method", r.Method).
			Str("endpoint", "/error").
			Str("error_id", "error123").
			Msg("Error endpoint accessed")

		http.Error(w, "You have reached the error endpoint", http.StatusInternalServerError)
	}
}
