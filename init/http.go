package init

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"goServer/server"
	"goServer/utils"
	"log"
	"net/http"
)

func HttpServerInit(port string) error {
	log.Println(" ------ Server Start ------ ")

	return http.ListenAndServe(port, registerHttpRouter())
}

func registerHttpRouter() http.Handler {
	router := mux.NewRouter()
	c := cors.AllowAll() // 일단 개발 편의상을 위해 전체 수용

	registerSubRouterTest(router)

	logMux := utils.LoggingMiddleware(router)
	// 라우팅 관련해서는 Mux쓰는 것이 훨씬 깔끔하고 좋다고 생각하기 떄문에 Mux로 관리

	corsRouter := c.Handler(logMux)
	return corsRouter
}

func registerSubRouterTest(router *mux.Router) {
	// 나중에 health Check용으로 사용 할 수도??
	testRouter := router.PathPrefix("/test").Subrouter()
	testRouter.HandleFunc("", server.NewGetTest().ServeHTTP).Methods("GET")
}

func HttpErrorChannelInit(channel chan error, logger *log.Logger) {

	go func() {
		for {
			select {
			case httpErr := <-channel:
				log.Println("Error : ", httpErr)
				logger.Println(httpErr)
			}
		}
	}()

}
