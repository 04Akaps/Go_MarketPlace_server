package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	initData "goServer/init"
	gRPC "goServer/myGRpc"
	"log"
	"math/big"
	"strings"
)

var envData initData.EnvData

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // 시간을 로그로 찍음
	envData = initData.InitEnv(".")

}

const mumbai = "wss://polygon-mumbai.g.alchemy.com/v2/xhCU8LxKi9YpyWv3Rue9RuZd7OlyMe_T"
const contractAbi = "[\n\t{\n\t\t\"anonymous\": false,\n\t\t\"inputs\": [\n\t\t\t{\n\t\t\t\t\"indexed\": false,\n\t\t\t\t\"internalType\": \"address\",\n\t\t\t\t\"name\": \"sender\",\n\t\t\t\t\"type\": \"address\"\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"indexed\": false,\n\t\t\t\t\"internalType\": \"uint256\",\n\t\t\t\t\"name\": \"currentNumber\",\n\t\t\t\t\"type\": \"uint256\"\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"indexed\": false,\n\t\t\t\t\"internalType\": \"uint256\",\n\t\t\t\t\"name\": \"nextNumber\",\n\t\t\t\t\"type\": \"uint256\"\n\t\t\t}\n\t\t],\n\t\t\"name\": \"EventEmit\",\n\t\t\"type\": \"event\"\n\t},\n\t{\n\t\t\"inputs\": [],\n\t\t\"name\": \"plus\",\n\t\t\"outputs\": [],\n\t\t\"stateMutability\": \"nonpayable\",\n\t\t\"type\": \"function\"\n\t},\n\t{\n\t\t\"inputs\": [],\n\t\t\"name\": \"number\",\n\t\t\"outputs\": [\n\t\t\t{\n\t\t\t\t\"internalType\": \"uint256\",\n\t\t\t\t\"name\": \"\",\n\t\t\t\t\"type\": \"uint256\"\n\t\t\t}\n\t\t],\n\t\t\"stateMutability\": \"view\",\n\t\t\"type\": \"function\"\n\t}\n]"

type EventEmit struct {
	Sender        common.Address
	CurrentNumber *big.Int
	NextNumber    *big.Int
}

func main() {
	go gRPC.GRpcServerInit()

	client, err := ethclient.Dial(mumbai)
	contractAddress := common.HexToAddress("0x3c5e85f16c755699E72868e4d995CE9D628Acb4d")
	//
	//result, err := client.BlockNumber(context.TODO())

	var EventEmitSignature = []byte("EventEmit(address,uint256,uint256)")
	var EventEmitTopic = common.BytesToHash(crypto.Keccak256Hash(EventEmitSignature).Bytes())

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(2394201),
		ToBlock:   nil,
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{EventEmitTopic}},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	abi, err := abi.JSON(strings.NewReader(contractAbi))

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:

			fmt.Println("vLOGGGGGGG", vLog) // 캐치한 이벤트를 출력합니다.
			//
			//type Log struct {
			//	Address     common.Address `json:"address" gencodec:"required"`
			//	Topics      []common.Hash  `json:"topics" gencodec:"required"`
			//	Data        []byte         `json:"data" gencodec:"required"`
			//	BlockNumber uint64         `json:"blockNumber"`
			//	TxHash      common.Hash    `json:"transactionHash" gencodec:"required"`
			//	TxIndex     uint           `json:"transactionIndex"`
			//	BlockHash   common.Hash    `json:"blockHash"`
			//	Index       uint           `json:"logIndex"`
			//	Removed     bool           `json:"removed"`
			//}

			fmt.Println(vLog.Topics)

			eventEmit := new(EventEmit)
			err := abi.UnpackIntoInterface(eventEmit, "EventEmit", vLog.Data)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(eventEmit.Sender)
			fmt.Println(eventEmit.CurrentNumber.Int64())
			fmt.Println(eventEmit.NextNumber.Int64())

			txHash := common.HexToHash(vLog.TxHash.String())
			tx, _, err := client.TransactionByHash(context.TODO(), txHash)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(tx.To().Hex())
		}
	}

	err = initData.HttpServerInit(envData)

	if err != nil {
		log.Fatal(err) // 굳이 서버를 안끌 필요가 없으니 그냥 바로 Fatal
	}
}
