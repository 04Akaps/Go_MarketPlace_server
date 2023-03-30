package gRpcUtils

import (
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
)

type TransferFromEvent struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

func TransferGoRoutine(sub ethereum.Subscription, logs chan types.Log, abi abi.ABI) {
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:

			transferEvent := new(TransferFromEvent)
			err := abi.UnpackIntoInterface(transferEvent, "TransferFrom", vLog.Data)

			if err != nil {
				log.Println(err)
			}

			fmt.Println(transferEvent.From.Hex())
			fmt.Println(transferEvent.To.Hex())
			fmt.Println(transferEvent.TokenId.Int64())

		}
	}
}
