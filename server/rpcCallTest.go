package server

import (
	"context"
	"github.com/gorilla/mux"
	"goServer/myGRpc/proto"
	"google.golang.org/grpc"
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

	client.CreateNewContract(context.Background(), &proto.CreateNewContractRequest{
		NewContract: newContract,
	})

}
