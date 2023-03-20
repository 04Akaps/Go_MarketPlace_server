package customError

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpServerLog struct {
	HttpServerErrLog chan error
	Logger           *log.Logger
}

func (httpServerLog *HttpServerLog) HttpErrorChannelInit() {
	go func() {
		for {
			select {
			case httpErr := <-httpServerLog.HttpServerErrLog:
				log.Println("Error : ", httpErr)
				httpServerLog.Logger.Println(httpErr)
			}
		}
	}()

}

type Launchpad struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// 기타 필드
}

type CustomError struct {
	Message string `json:"message"`
}

func NewHandlerError(w http.ResponseWriter, message string, errorCode int) {
	w.WriteHeader(errorCode)
	customError := CustomError{
		Message: message,
	}
	json.NewEncoder(w).Encode(customError)
}
