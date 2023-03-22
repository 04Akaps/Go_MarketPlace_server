package server

import (
	"encoding/json"
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

type MakeLaunchpadReq struct {
	FirstOwnerEmail  string   `json:"first_owner_email"`
	CaAddress        string   `json:"ca_address"`
	ChainId          string   `json:"chain_id"`
	Price            int64    `json:"price"`
	AirdropAddress   []string `json:"airdrop_address"`
	WhiteListAddress []string `json:"whitelist_address"`
}

func (controller *LaunchpadController) MakeLaunchpad(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Make Launchpad")

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // body의 최대 크기
	dec := json.NewDecoder(r.Body)                   // decoder
	dec.DisallowUnknownFields()                      //

	var req MakeLaunchpadReq
	err := dec.Decode(&req)

	fmt.Println(req.CaAddress)

	if err != nil {
		customError.NewHandlerError(w, customError.NewPostErrorHandler(err), 200)
		return
	}

	hashValue := req.CaAddress[:7]
	fmt.Println(hashValue)
	// SQL Insert 필요
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
	//customError.NewHandlerError(w, "ac", 200)

	// 데이터 문제 없을 떄,
	//utils.SuccesResponse(w, data)
}
