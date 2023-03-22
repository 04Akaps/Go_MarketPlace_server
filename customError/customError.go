package customError

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func NewPostErrorHandler(err error) string {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		return fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		return fmt.Sprintf("Request body contains badly-formed JSON")
	case errors.As(err, &unmarshalTypeError):
		return fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
	case errors.Is(err, io.EOF):
		return "Request body must not be empty"
	case err.Error() == "http: request body too large":
		return "Request body must not be larger than 1MB"
	default:
		return "Check Your Post Body Keys"
	}
}
