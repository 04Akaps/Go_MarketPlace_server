package crypo

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func NewCryptoClient(nodeUrl string) *ethclient.Client {
	context := context.Background()
	cryptoClient, err := ethclient.DialContext(context, nodeUrl)
	if err != nil {
		log.Println("Error New CryptoClient : ", err)
		fmt.Println(err)
		return nil
	}

	return cryptoClient
}
