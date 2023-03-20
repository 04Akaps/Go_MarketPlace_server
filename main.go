package main

import (
	"goServer/customError"
	initData "goServer/init"
	"goServer/utils"
	"log"
)

var envData initData.EnvData
var httpServerErrLog customError.HttpServerLog

func init() {
	log.SetFlags(log.Ldate | log.Ltime) // 시간을 로그로 찍음
	envData = initData.InitEnv(".")

	httpServerErrLog = customError.HttpServerLog{
		HttpServerErrLog: make(chan error),
		Logger:           utils.GetHttpLogFile("httpErrorLog/"),
	}
	httpServerErrLog.HttpErrorChannelInit()

}

func main() {

	err := initData.HttpServerInit(envData.HttpServerPort, httpServerErrLog.HttpServerErrLog)

	if err != nil {
		httpServerErrLog.HttpServerErrLog <- err
	}

	//// 메인 루틴이 죽으면 모든 루틴이 죽어버리니깐 프로세스에 대한 시그널로 메인 루틴을 안죽게 설정
	//stop := make(chan os.Signal, 1)
	//signal.Notify(stop, os.Interrupt)
	//<-stop
	// 테스트 용으로 작성한 것이기 떄문에 일단 주석
}
