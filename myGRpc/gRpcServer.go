package myGRpc

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"goServer/myGRpc/gRpcUtils"
	"goServer/myGRpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"
)

type server struct {
	proto.NewContractServiceServer
	addr string
	ctx  context.Context
}

var (
	addrs     = []string{":50051", ":50052"}
	abiString string
	client    *ethclient.Client
)

const mumbai = "wss://polygon-mumbai.g.alchemy.com/v2/xhCU8LxKi9YpyWv3Rue9RuZd7OlyMe_T"

func (*server) CreateNewContract(ctx context.Context, req *proto.CreateNewContractRequest) (*proto.CreateNewContractResponse, error) {
	fmt.Println("Create New Contract req")

	newContract := req.GetNewContract()
	contractAddress := common.HexToAddress(newContract.Contract)

	var EventEmitSignature = []byte("TransferFrom(address,address,uint256)")
	var EventEmitTopic = common.BytesToHash(crypto.Keccak256Hash(EventEmitSignature).Bytes())

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(2394201),
		ToBlock:   nil,
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{EventEmitTopic}},
	}

	logs := make(chan types.Log)
	sub, _ := client.SubscribeFilterLogs(context.Background(), query, logs)
	abi, _ := abi.JSON(strings.NewReader(abiString))

	go gRpcUtils.TransferGoRoutine(sub, logs, abi) // 루틴을 생성
	// 해당 값을 DB에 저장을 하며, 발생하는 이벤트에 대해서 루틴을 돌려 주어야 한다.

	return &proto.CreateNewContractResponse{
		NewContract: newContract,
	}, nil
}

func GRpcServerInit() {
	log.Println("------------- Proto gRPC Server ----------")
	ethClient, err := ethclient.Dial(mumbai)

	if err != nil {
		log.Fatal(err)
	}
	client = ethClient

	file, err := os.Open("./myGRpc/abi.json")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		abiString += scanner.Text()
	}

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
	//creds, _ := credentials.NewServerTLSFromFile("cert.pem", "key.pem")

	opts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(100),
		//grpc.Creds(creds),
		//grpc.UnaryInterceptor(gRpcUtils.EnsureValidToken),
		//grpc.KeepaliveEnforcementPolicy(gRpcUtils.GetKeepAliveEnforcement(5 * time.Second)),
		//grpc.KeepaliveParams(gRpcUtils.GetKeepAliveServerParameters(15*time.Second, 30*time.Second, 5*time.Second, 5*time.Second, 1*time.Second)),
	}

	gRpcClient := grpc.NewServer(opts...)
	proto.RegisterNewContractServiceServer(gRpcClient, &server{addr: addr, ctx: context.Background()})

	log.Printf("gRPC Server Start on %s\n", addr)
	reflection.Register(gRpcClient)

	if err := gRpcClient.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
