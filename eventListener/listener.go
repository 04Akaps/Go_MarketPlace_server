package eventListener

import (
	"fmt"
	"goServer/eventListener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct {
	proto.NewContractServiceServer
}

type gRpcOption struct{}

func (g *gRpcOption) apply(opts *grpc.ServerOption) grpc.ServerOption {
	fmt.Println(g)
	return nil
}

func GRpcServer() {
	log.Println("------------- Proto gRPC Client Server ----------")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Error  : ", err)
	}

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(3),
		grpc.Creds(credentials.NewTLS(nil)),
	}

	gRpcClient := grpc.NewServer(opts...)
	proto.RegisterNewContractServiceServer(gRpcClient, &server{})
	reflection.Register(gRpcClient)

	log.Println(" ------------- gRPC Server Start ------------- ")

	if err := gRpcClient.Serve(lis); err != nil {
		log.Fatal(" : Error is ocured : ", err)
	}
}
