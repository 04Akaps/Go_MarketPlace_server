package init

import (
	"github.com/spf13/viper"
	"log"
)

type EnvData struct {
	HttpServerPort string `mapstructure:"http_server_port"`
	DbUserName     string `mapstructure:"db_username"`
	DbPassword     string `mapstructure:"de_password"`
	DbEndPoint     string `mapstructure:"db_endpoint"`
	CryptoNodeUrl  string `mapstructure:"crypto_node_url"`
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
