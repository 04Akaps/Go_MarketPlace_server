package utils

import (
	"encoding/json"
	"goServer/customError"
	p "goServer/paseto"
	"log"
	"net/http"
	"os"
	"strings"
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

func LoggingMiddleware(next http.Handler, paseto p.PasetoInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Content-Type") != "" {
			r.Header.Set("Content-Type", "application/json")
		}
		w.Header().Set("Content-Type", "application/json")

		if !strings.Contains(r.URL.String(), "auth") {
			authCookie, err := r.Cookie("auth")

			if err == http.ErrNoCookie {
				customError.NewHandlerError(w, "Auth Token이 없습니다. ", 200)
				return
			}

			err = paseto.VerifyToken(authCookie.Value)

			if err != nil {
				customError.NewHandlerError(w, "Auth 인증 실패...!!", 200)
				return
			}
		}
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func BodyParserDecoder(w http.ResponseWriter, r *http.Request) *json.Decoder {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // body의 최대 크기
	dec := json.NewDecoder(r.Body)                   // decoder
	dec.DisallowUnknownFields()
	return dec
}

func SuccesResponse(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
