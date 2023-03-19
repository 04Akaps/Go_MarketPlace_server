package main

import (
	initData "goServer/init"
	"goServer/utils"
	"log"
	"os"
	"os/signal"
)

var envData initData.EnvData
var httpServerErrLog chan error

func init() {
	log.SetFlags(log.Ldate | log.Ltime) // 시간을 로그로 찍음
	httpServerErrLog = make(chan error)
	envData = initData.InitEnv(".")

	logger, file := utils.GetHttpLogFile()
	initData.HttpErrorChannelInit(httpServerErrLog, logger, file)
}

func main() {

	err := initData.HttpServerInit(envData.HttpServerPort)

	if err != nil {
		httpServerErrLog <- err
	}

	// 메인 루틴이 죽으면 모든 루틴이 죽어버리니깐 프로세스에 대한 시그널로 메인 루틴을 안죽게 설정
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
