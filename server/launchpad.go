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
		http.Error(w, QUERY_ID_EMPTY_STRING, http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(customError.NewHandlerError("abc", 200)) // 후에 데이터를 내려 줘야함
}
