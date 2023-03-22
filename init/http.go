package init

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"goServer/server"
	"goServer/utils"
	"log"
	"net/http"
)

func HttpServerInit(port string, channel chan error) error {
	log.Println(" ------ Server Start ------ ")
	return http.ListenAndServe(port, registerHttpRouter(channel))
}

func registerHttpRouter(channel chan error) http.Handler {
	// 라우팅 관련해서는 Mux쓰는 것이 훨씬 깔끔하고 좋다고 생각하기 떄문에 Mux로 관리
	router := mux.NewRouter()

	logMux := utils.LoggingMiddleware(router) // 들어오는 요청에 대해서 로그 설정
	c := cors.AllowAll()                      // 일단 개발 편의상을 위해 전체 수용

	registerTestRouter(router)
	registerLaunchpadRouter(router, channel)

	corsRouter := c.Handler(logMux)
	return corsRouter
}

func registerTestRouter(router *mux.Router) {
	// 나중에 health Check용으로 사용 할 수도??
	testRouter := router.PathPrefix("/test").Subrouter()
	testRouter.HandleFunc("", server.NewGetTest().ServeHTTP).Methods("GET")
}

func registerLaunchpadRouter(router *mux.Router, channel chan error) {
	// MarketPlace에서 모든 블록을 계속 패칭하는 것은 개인 개발상으로 어렵고, 리소스 낭비가 너무 하다고정생각이 들기 떄문에
	// Launchpad에서 만들어지는 NFT를 거래하는 부분만 다룰 예정
	launchpadRouter := router.PathPrefix("/launchpad").Subrouter()

	controller := server.NewLaunchpadController(channel)

	launchpadRouter.HandleFunc("", controller.GetLaunchpadData).Methods("GET")
	launchpadRouter.HandleFunc("", controller.MakeLaunchpad).Methods("POST")

}
