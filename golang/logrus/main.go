package main

import (
	"fmt"
	"net/http"
	"os"

	"io"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Setting up Logrus with a multi-writer
	logger := log.New()
	mw := io.MultiWriter(os.Stdout) // MultiWriter to log to stdout and file
	logger.SetOutput(mw)
	// logger.Formatter = &log.JSONFormatter{}
	logger.Level = log.DebugLevel // Setting the log level to debug

	http.HandleFunc("/", handleIndex(logger))
	http.HandleFunc("/log", handleLog(logger))
	http.HandleFunc("/data", handleData(logger))
	http.HandleFunc("/error", handleError(logger))

	fmt.Println("Server starting on <http://localhost:8080>")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func handleIndex(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(log.Fields{
			"method": r.Method,
		}).Info("Accessing index page")
		fmt.Fprintln(w, "Welcome to the Go Application!")
	}
}

func handleLog(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fields := log.Fields{
			"Method": r.Method,
			"Path":   r.URL.Path,
		}
		switch r.Method {
		case "GET":
			logger.WithFields(fields).Info("Handled GET request on /log")
			fmt.Fprintln(w, "Received a GET request at /log.")
		case "POST":
			logger.WithFields(fields).Info("Handled POST request on /log")
			fmt.Fprintln(w, "Received a POST request at /log.")
		default:
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
		}
	}
}

func handleData(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(log.Fields{
			"method":     r.Method,
			"endpoint":   "/data",
			"request_id": fmt.Sprintf("%d", os.Getpid()),
		}).Info("Data endpoint hit")
		fmt.Fprintln(w, "This is the data endpoint. Method used:", r.Method)
	}
}

func handleError(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.WithFields(log.Fields{
			"method":   r.Method,
			"endpoint": "/error",
			"error_id": "error123",
		}).Error("Error endpoint accessed")
		http.Error(w, "You have reached the error endpoint", http.StatusInternalServerError)
	}
}
