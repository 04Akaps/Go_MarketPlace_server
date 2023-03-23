package crypo

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func NewCryptoClient(nodeUrl string) *ethclient.Client {
	context := context.Background()
	cryptoClient, err := ethclient.DialContext(context, nodeUrl)

	if err != nil {
		log.Println("Error New CryptoClient : ", err)
		return nil
	}

	return cryptoClient
}
