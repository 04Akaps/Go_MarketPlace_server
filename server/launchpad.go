package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"goServer/customError"
	sqlc "goServer/mysql/sqlc"
	"goServer/utils"
	"net/http"
)

const (
	QUERY_ID_EMPTY_STRING    = "id값이 빈 문자열 입니다."
	CODE_AT_FAILED           = "CODE At에 실해 했습니다."
	INVALID_CONTRACT_ADDRESS = "잘못된 Contract주소 입니다."
)

var (
	queryIdEmpty = errors.New(QUERY_ID_EMPTY_STRING)
)

type LaunchpadController struct {
	ErrorChannel chan error
	DBClient     *sqlc.Queries
	ctx          context.Context
	cryptoClient *ethclient.Client
}

type LaunchpadInterface interface {
	MakeLaunchpad(http.ResponseWriter, *http.Request)
	GetLaunchpadByHashValue(http.ResponseWriter, *http.Request)
	GetLaunchpadsByChainId(http.ResponseWriter, *http.Request)
}

func NewLaunchpadController(channel chan error, dbClient *sqlc.Queries, cryptoClient *ethclient.Client) LaunchpadInterface {
	context := context.Background()
	return &LaunchpadController{
		ErrorChannel: channel,
		DBClient:     dbClient,
		ctx:          context,
		cryptoClient: cryptoClient,
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
	var req MakeLaunchpadReq

	dec := utils.BodyParserDecoder(w, r)
	err := dec.Decode(&req)

	if err != nil {
		errorMsg := customError.NewPostErrorHandler(err)
		controller.ErrorChannel <- errors.New(errorMsg)
		customError.NewHandlerError(w, errorMsg, 200)
		return
	}

	address := common.HexToAddress(req.CaAddress)
	byteCode, err := controller.cryptoClient.CodeAt(controller.ctx, address, nil)

	if err != nil {
		controller.ErrorChannel <- errors.New(CODE_AT_FAILED)
		customError.NewHandlerError(w, CODE_AT_FAILED, 200)
		return
	}

	isContractAddress := len(byteCode) > 0

	if !isContractAddress {
		controller.ErrorChannel <- errors.New(INVALID_CONTRACT_ADDRESS)
		customError.NewHandlerError(w, INVALID_CONTRACT_ADDRESS, 200)
		return
	}

	fmt.Println(address)

	hashValue := req.CaAddress[:7]
	fmt.Println(hashValue)
	// SQL Insert 필요
}

func (controller *LaunchpadController) GetLaunchpadByHashValue(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("id")
	// hash값을 통해서 데이터를 가져 올 예정
	if hash == "" {
		controller.ErrorChannel <- queryIdEmpty
		customError.NewHandlerError(w, QUERY_ID_EMPTY_STRING, 200)
		return
	}

	fmt.Println(hash)

	// 데이터 문제 있을 떄
	//customError.NewHandlerError(w, "ac", 200)

	// 데이터 문제 없을 떄,
	//utils.SuccesResponse(w, data)
}

func (controller *LaunchpadController) GetLaunchpadsByChainId(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	chainId := vars["chainId"]
	//r.URL.pa

	fmt.Println(chainId)
}
