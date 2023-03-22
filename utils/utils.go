package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func GetHttpLogFile(path string) *log.Logger {
	t := time.Now()
	startTime := t.Format("2006-01-02 15:04:05")
	logFile, err := os.Create("serverLog/" + path + startTime + ".log")
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(logFile, "", log.LstdFlags)

	return logger
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "" {
			r.Header.Set("Content-Type", "application/json")
		}
		w.Header().Set("Content-Type", "application/json")
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func SuccesResponse(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
