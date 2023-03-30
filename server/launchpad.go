package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"goServer/customError"
	"goServer/myGRpc/proto"
	sqlc "goServer/mysql/sqlc"
	"goServer/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strings"
)

const (
	QUERY_ID_EMPTY_STRING     = "id값이 빈 문자열 입니다."
	CODE_AT_FAILED            = "CODE At에 실해 했습니다."
	INVALID_CONTRACT_ADDRESS  = "잘못된 Contract주소 입니다."
	INSERT_DB_ERRROR          = "Insert가 실패하였습니다."
	GET_LAUNCHPAD_BY_HASH_ID  = "Hash ID를 통한 Get 요청 실패"
	HASH_EMPRY_STRING         = "hash값이 빈 문자열 입니다."
	GET_LAUNCHPAD_BY_CHAIN_ID = "Chain ID를 통한 Get 요청 실패"
)

type LaunchpadController struct {
	ErrorChannel chan error
	DBClient     *sqlc.Queries
	ctx          context.Context
	cryptoClient *ethclient.Client
	gRPC         *grpc.ClientConn
}

type LaunchpadInterface interface {
	MakeLaunchpad(http.ResponseWriter, *http.Request)
	GetLaunchpadByHashValue(http.ResponseWriter, *http.Request)
	GetLaunchpadsByChainId(http.ResponseWriter, *http.Request)
}

func NewLaunchpadController(channel chan error, dbClient *sqlc.Queries, cryptoClient *ethclient.Client) LaunchpadInterface {
	context := context.Background()

	//cred, err := credentials.NewClientTLSFromFile("cert.pem", "example.com")
	//
	//if err != nil {
	//	log.Fatal("TlsFromFile Error ", err)
	//}

	opts := []grpc.DialOption{
		//insecure.NewCredentials()
		//grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 인증서 없이 테스트 하기 위함
		//grpc.WithPerRPCCredentials(gRpcUtils.GetTokenSource()), // 매 요청마다 인증서 확인
		//grpc.WithKeepaliveParams(gRpcUtils.GetKeepAliveClientParameters(10*time.Second, time.Second)),
	}

	roundBin, _ := grpc.Dial(
		"localhost:50051",
		opts...,
	)

	return &LaunchpadController{
		ErrorChannel: channel,
		DBClient:     dbClient,
		ctx:          context,
		cryptoClient: cryptoClient,
		gRPC:         roundBin,
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

	hashValue := req.CaAddress[:7]

	whiteListJSON, _ := json.Marshal(req.WhiteListAddress)
	airdropListJson, _ := json.Marshal(req.AirdropAddress)

	dbReq := sqlc.CreateNewLaunchpadParams{
		HashValue:        hashValue,
		FirstOwnerEmail:  req.FirstOwnerEmail,
		CaAddress:        req.CaAddress,
		ChainID:          req.ChainId,
		Price:            int32(req.Price),
		AirdropAddress:   whiteListJSON,
		WhitelistAddress: airdropListJson,
	}

	_, err = controller.DBClient.CreateNewLaunchpad(controller.ctx, dbReq)

	if err != nil {
		controller.ErrorChannel <- err
		customError.NewHandlerError(w, INSERT_DB_ERRROR, 200)
		return
	}

	client := proto.NewNewContractServiceClient(controller.gRPC)
	newContract := &proto.NewContract{
		Contract: req.CaAddress,
	}

	_, err = client.CreateNewContract(context.Background(), &proto.CreateNewContractRequest{
		NewContract: newContract,
	})

	if err != nil {
		log.Fatal(err)
	}

	customError.NewHandlerError(w, "Success", 200)
}

func (controller *LaunchpadController) GetLaunchpadByHashValue(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("id")
	if len(strings.TrimSpace(hash)) == 0 {
		controller.ErrorChannel <- errors.New(QUERY_ID_EMPTY_STRING)
		customError.NewHandlerError(w, QUERY_ID_EMPTY_STRING, 200)
		return
	}

	launchpad, err := controller.DBClient.GetLaunchpadByHash(controller.ctx, hash)

	if err != nil {
		controller.ErrorChannel <- err
		customError.NewHandlerError(w, GET_LAUNCHPAD_BY_HASH_ID, 200)
		return
	}

	utils.SuccesResponse(w, launchpad)
}

func (controller *LaunchpadController) GetLaunchpadsByChainId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chainId := vars["chainId"]

	if len(strings.TrimSpace(chainId)) == 0 {
		controller.ErrorChannel <- errors.New(HASH_EMPRY_STRING)
		customError.NewHandlerError(w, HASH_EMPRY_STRING, 200)
		return
	}

	launchpads, err := controller.DBClient.GetLaunchpadByChainId(controller.ctx, chainId)

	if err != nil {
		controller.ErrorChannel <- err
		customError.NewHandlerError(w, GET_LAUNCHPAD_BY_CHAIN_ID, 200)
		return
	}

	utils.SuccesResponse(w, launchpads)
}
