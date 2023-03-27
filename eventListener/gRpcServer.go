package eventListener

import (
	"goServer/eventListener/gRpcUtils"
	"goServer/eventListener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
	"time"
)

type server struct {
	proto.NewContractServiceServer
	addr string
}

var (
	addrs = []string{":50051", ":50052"}
)

func GRpcServerInit() {
	log.Println("------------- Proto gRPC Server ----------")

	var wg sync.WaitGroup
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			startServer(addr)
		}(addr)
	}
	wg.Wait()
}

func startServer(addr string) {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(3),
		grpc.Creds(credentials.NewTLS(nil)),
		grpc.UnaryInterceptor(gRpcUtils.EnsureValidToken),
		grpc.KeepaliveEnforcementPolicy(gRpcUtils.GetKeepAliveEnforcement(5 * time.Second)),
		grpc.KeepaliveParams(gRpcUtils.GetKeepAliveServerParameters(15*time.Second, 30*time.Second, 5*time.Second, 5*time.Second, 1*time.Second)),
	}

	gRpcClient := grpc.NewServer(opts...)
	proto.RegisterNewContractServiceServer(gRpcClient, &server{addr: addr})

	log.Printf("gRPC Server Start on %s\n", addr)
	reflection.Register(gRpcClient)

	if err := gRpcClient.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
