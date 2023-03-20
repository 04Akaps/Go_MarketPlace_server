package init

import (
	"github.com/gorilla/mux"
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

	registerSubRouterTest(router)

	logMux := utils.LoggingMiddleware(router)
	// 라우팅 관련해서는 Mux쓰는 것이 훨씬 깔끔하고 좋다고 생각하기 떄문에 Mux로 관리

	return logMux
}

func registerSubRouterTest(router *mux.Router) {
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
