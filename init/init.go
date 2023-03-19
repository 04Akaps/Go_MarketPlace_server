package init

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

type EnvData struct {
	HttpServerPort string `mapstructure:"http_server_port"`
}

func InitEnv(path string) EnvData {
	var goConfig EnvData

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("env Read Error : &w", err)
	}

	if err := viper.Unmarshal(&goConfig); err != nil {
		log.Fatal("env Marshal Error : &w", err)
	}

	return goConfig
}

func HttpServerInit(port string) error {
	log.Println(" ------ Server Start ------ ")
	return http.ListenAndServe(port, nil)
}

func HttpErrorChannelInit(channel chan error, logger *log.Logger, file *os.File) {

	go func() {
		for {
			select {
			case httpErr := <-channel:
				log.Println("Error : ", httpErr)
				logger.Println(httpErr)
			}
		}

		defer file.Close() // 어차피 서버가 꺼지기 전에는 defer를 실행을 시켜야 한다
		// 메인 루틴에서 죽으면 해당 루틴도 죽으니, 그떄 함꼐 꺼지게 구성
		// 메인에 적고 싶지는 않아서 서브 함수로 뺴서 작업
	}()

}
