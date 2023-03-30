package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"goServer/myGRpc/proto"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

// RPC Call을 테스트 할 목적
// 나중에는 Launchpad로 이동 예정

type RpcCallTest struct {
	rpc *grpc.ClientConn
	ctx context.Context
}

type RpcCallImpl interface {
	SendRpcCall(http.ResponseWriter, *http.Request)
}

func NewRpcCallTest(grpc *grpc.ClientConn) RpcCallImpl {

	return &RpcCallTest{
		rpc: grpc,
		ctx: context.Background(),
	}
}

func (rpc *RpcCallTest) SendRpcCall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ca := vars["ca"]

	client := proto.NewNewContractServiceClient(rpc.rpc)

	newContract := &proto.NewContract{
		Contract: ca,
	}

	response, err := client.CreateNewContract(context.Background(), &proto.CreateNewContractRequest{
		NewContract: newContract,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)

}
