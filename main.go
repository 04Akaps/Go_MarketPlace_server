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
		Logger:           utils.GetLogFile("httpErrorLog/"),
	}
	httpServerErrLog.HttpErrorChannelInit()

}

func main() {

	err := initData.HttpServerInit(envData, httpServerErrLog.HttpServerErrLog)

	if err != nil {
		log.Fatal(err) // 굳이 서버를 안끌 필요가 없으니 그냥 바로 Fatal
	}

}
