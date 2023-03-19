package init

import (
	"goServer/server"
	"log"
	"net/http"
)

func HttpServerInit(port string) error {
	log.Println(" ------ Server Start ------ ")

	return http.ListenAndServe(port, registerHttpRouter())
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

func registerHttpRouter() *http.ServeMux {
	router := &http.ServeMux{}

	router.Handle("/test", server.NewGetTest())

	return router
}
