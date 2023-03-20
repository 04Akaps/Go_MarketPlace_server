package customError

import (
	"encoding/json"
	"fmt"
	"log"
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

type CustomError struct {
	message   string
	errorCode int
}

func NewHandlerError(message string, errorCode int) string {
	// 일단 두기는 했는데.. 나중에 사용하면 사용하고 아니면 삭제 예정
	customError := &CustomError{
		message:   message,
		errorCode: errorCode,
	}

	fmt.Println(customError)
	byteData, err := json.Marshal(customError)

	fmt.Println(message, errorCode)
	if err != nil {
		log.Println("NewHandlerError Marshaling Error")
		return "Error Chcek Server"
	}

	fmt.Println("====", string(byteData))

	return string(byteData)
}
