package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"goServer/customError"
	initData "goServer/init"
	redis2 "goServer/redis"
	"goServer/utils"
	"log"
	"time"
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

	option := &redis.Options{
		DB:              0,
		ClientName:      "NFT_Market_go",
		ConnMaxIdleTime: 1 * time.Second,
		ConnMaxLifetime: 1 * time.Second,
		MaxIdleConns:    1000,
		PoolSize:        25,
	}

	redis2.NewRedisClient(option, context.Background())

	err := initData.HttpServerInit(envData, httpServerErrLog.HttpServerErrLog)

	if err != nil {
		log.Fatal(err) // 굳이 서버를 안끌 필요가 없으니 그냥 바로 Fatal
	}

}
