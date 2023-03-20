package server

import (
	"errors"
	"fmt"
	"goServer/customError"
	"net/http"
)

const (
	QUERY_ID_EMPTY_STRING = "id값이 빈 문자열 입니다."
)

var (
	queryIdEmpty = errors.New(QUERY_ID_EMPTY_STRING)
)

type LaunchpadController struct {
	ErrorChannel chan error
}

type LaunchpadInterface interface {
	MakeLaunchpad(http.ResponseWriter, *http.Request)
	GetLaunchpadData(http.ResponseWriter, *http.Request)
}

func NewLaunchpadController(channel chan error) LaunchpadInterface {
	return &LaunchpadController{
		ErrorChannel: channel,
	}
}

func (controller *LaunchpadController) MakeLaunchpad(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Make Launchpad")
}

func (controller *LaunchpadController) GetLaunchpadData(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("id")
	// hash값을 통해서 데이터를 가져 올 예정
	if hash == "" {
		controller.ErrorChannel <- queryIdEmpty
		customError.NewHandlerError(w, QUERY_ID_EMPTY_STRING, 200)
		return
	}

	// 데이터 문제 있을 떄
	//customError.NewHandlerError(w, <errorMessage>, 200)

	// 데이터 문제 없을 떄,
	//utils.SuccesResponse(w, data)
}
