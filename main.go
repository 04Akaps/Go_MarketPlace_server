package main

import (
	"fmt"
	gRPC "goServer/eventListener"
	gRPCUtils "goServer/eventListener/gRpcUtils"
	initData "goServer/init"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

var envData initData.EnvData

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // 시간을 로그로 찍음
	envData = initData.InitEnv(".")

}

const (
	gRpcScheme      = "market"
	gRpcServiceName = "lb.market.grpc.io"
)

func main() {
	go gRPC.GRpcServerInit()

	cred, err := credentials.NewClientTLSFromFile("cert.pem", "")

	if err != nil {
		log.Fatal("TlsFromFile Error ", err)
	}

	opts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithTransportCredentials(cred),                    // 인증서 없이 테스트 하기 위함
		grpc.WithPerRPCCredentials(gRPCUtils.GetTokenSource()), // 매 요청마다 인증서 확인
		grpc.WithKeepaliveParams(gRPCUtils.GetKeepAliveClientParameters(10*time.Second, time.Second)),
	}

	roundBin, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", gRpcScheme, gRpcServiceName),
		opts...,
	)

	if err != nil {
		log.Fatal(err) // 굳이 서버를 안끌 필요가 없으니 그냥 바로 Fatal
	}

	defer roundBin.Close()

	err = initData.HttpServerInit(envData)

	if err != nil {
		log.Fatal(err) // 굳이 서버를 안끌 필요가 없으니 그냥 바로 Fatal
	}
}
