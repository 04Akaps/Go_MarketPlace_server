package utils

import (
	"encoding/json"
	"goServer/customError"
	goServer "goServer/mysql/sqlc"
	p "goServer/paseto"
	"goServer/redis"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"
)

func GetLogFile(path string) *log.Logger {
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

func RedisLaunchpadMiddleWare(next http.Handler, redis redis.RedisImpl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		path := r.URL.Path
		parts := strings.Split(path, "/")
		redisKey := ""

		if len(parts) > 2 && parts[1] == "chainId" {
			redisKey = "launchpads:chainId:" + parts[2]
		} else {
			// 모든 데이터를 가져 오는 행위 일 떄,
			redisKey = "launchpads"
		}

		val, err := redis.GetDataFromRedis(redisKey)

		if err != nil {
			// 에러가 뜨는 경우에는 데이터가 없는 경우에 뜨기 떄문에 controller로 돌려줌
			recoder := httptest.NewRecorder()
			next.ServeHTTP(recoder, r)

			if recoder.Code == http.StatusOK {
				_, err = redis.SetDataToRedis(redisKey, recoder.Body.String())

				if err != nil {
					// 데이터 저장을 실패 하였 을 떄,
					log.Println(err)
					customError.NewHandlerError(w, "데이터 저장 실패", 200)
					return
				}
			}
		} else {
			var launchpads goServer.Launchpad
			err := json.Unmarshal(val, &launchpads)

			if err != nil {
				log.Println(err)
				customError.NewHandlerError(w, "데이터 마샬링 실패", 200)
				return
			}
			SuccesResponse(w, launchpads)
			return
		}

	}
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
