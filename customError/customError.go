package customError

import (
	"encoding/json"
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
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code"`
}

func NewHandlerError(message string, errorCode int) []byte {
	// 일단 두기는 했는데.. 나중에 사용하면 사용하고 아니면 삭제 예정
	customError := CustomError{
		Message:   message,
		ErrorCode: errorCode,
	}

	byteData, err := json.Marshal(customError)

	if err != nil {
		log.Println("NewHandlerError Marshaling Error")
		return []byte("파싱 에러 서버 확인 필요")
	}

	return byteData
}
