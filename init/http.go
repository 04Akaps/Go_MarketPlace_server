package init

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"goServer/crypo"
	sqlc "goServer/mysql/sqlc"
	"goServer/paseto"
	r "goServer/redis"
	"goServer/server"
	"goServer/utils"
	"log"
	"net/http"
	"time"
)

func HttpServerInit(envData EnvData, channel chan error) error {
	log.Println(" ------ Server Start ------ ")

	dbClient := NewDBClient("mysql", envData.DbUserName, envData.DbPassword, "launchpad", envData.DbEndPoint, "3306")
	initOAuth(envData)
	return http.ListenAndServe(envData.HttpServerPort, registerHttpRouter(channel, dbClient, envData))
}

func registerHttpRouter(channel chan error, dbClient *sqlc.Queries, envData EnvData) http.Handler {
	// 라우팅 관련해서는 Mux쓰는 것이 훨씬 깔끔하고 좋다고 생각하기 떄문에 Mux로 관리
	router := mux.NewRouter()

	newPasto := paseto.NewPasetoMaker(envData.PaseToKey)

	logMux := utils.LoggingMiddleware(router, newPasto) // 들어오는 요청에 대해서 로그 설정
	c := cors.AllowAll()                                // 개발 편의상을 위해 전체 수용

	corsRouter := c.Handler(logMux)

	registerTestRouter(router)

	registerLaunchpadRouter(router, channel, dbClient, envData)
	registerAuthRouter(router, channel, newPasto)

	return corsRouter
}

func registerTestRouter(router *mux.Router) {
	// 나중에 health Check용으로 사용 할 수도??
	testRouter := router.PathPrefix("/test").Subrouter()
	testRouter.HandleFunc("", server.NewGetTest().ServeHTTP).Methods("GET")
}

func registerLaunchpadRouter(router *mux.Router, channel chan error, dbClient *sqlc.Queries, envData EnvData) {

	launchpadRouter := router.PathPrefix("/launchpad").Subrouter()
	controller := server.NewLaunchpadController(channel, dbClient, crypo.NewCryptoClient(envData.CryptoNodeUrl))

	option := &redis.Options{
		DB:              0,
		ClientName:      "NFT_Market_go",
		ConnMaxIdleTime: 30 * time.Minute,
		ConnMaxLifetime: 1 * time.Minute,
		MaxIdleConns:    1000,
		PoolSize:        25,
	}
	client := r.NewRedisClient(option, context.Background())
	client.SetRedisLoggerFile(utils.GetLogFile("redisLog/"), make(chan error))

	launchpadRouter.HandleFunc("", utils.RedisLaunchpadMiddleWare(http.HandlerFunc(controller.GetLaunchpadByHashValue), client)).Methods("GET")
	launchpadRouter.HandleFunc("/chainId/{chainId}", utils.RedisLaunchpadMiddleWare(http.HandlerFunc(controller.GetLaunchpadsByChainId), client)).Methods("GET")
	launchpadRouter.HandleFunc("", controller.MakeLaunchpad).Methods("POST")
}

func registerAuthRouter(router *mux.Router, channel chan error, paseto paseto.PasetoInterface) {
	authRouter := router.PathPrefix("/auth").Subrouter()

	authController := server.NewAuthController(channel, paseto)

	authRouter.HandleFunc("/{action}/{provider}", authController.Auth)
	authRouter.HandleFunc("/logout", authController.Logout)
}
