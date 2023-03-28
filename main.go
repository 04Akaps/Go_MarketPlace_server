package main

import (
	initData "goServer/init"
	gRPC "goServer/myGRpc"
	"log"
)

var envData initData.EnvData

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // 시간을 로그로 찍음
	envData = initData.InitEnv(".")

}

func main() {
	go gRPC.GRpcServerInit()

	err := initData.HttpServerInit(envData)

	if err != nil {
		log.Fatal(err) // 굳이 서버를 안끌 필요가 없으니 그냥 바로 Fatal
	}
}
