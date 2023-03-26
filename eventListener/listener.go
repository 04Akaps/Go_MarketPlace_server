package main

import (
	"fmt"
	"goServer/eventListener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct {
	proto.NewContractServiceServer
}

func main() {
	log.Println("------------- Proto gRPC CLient Server ----------")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Error  : ", err)
	}

	var opts []grpc.ServerOption

	gRpcClient := grpc.NewServer(opts...)
	proto.RegisterNewContractServiceServer(gRpcClient, &server{})
	reflection.Register(gRpcClient)

	fmt.Println("gRPC Server Start")

	if err := gRpcClient.Serve(lis); err != nil {
		log.Fatal(" : Error is ocured : ", err)
	}
}
