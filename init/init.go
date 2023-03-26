package init

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type EnvData struct {
	HttpServerPort     string `mapstructure:"http_server_port"`
	DbUserName         string `mapstructure:"db_username"`
	DbPassword         string `mapstructure:"de_password"`
	DbEndPoint         string `mapstructure:"db_endpoint"`
	CryptoNodeUrl      string `mapstructure:"crypto_node_url"`
	GoogleAuthId       string `mapstructure:"google_auth_id"`
	GoogleAuthPassword string `mapstructure:"google_auth_password"`
	AuthKey            string `mapstructure:"auth_key"`
	PaseToKey          string `mapstructure:"paseto_key""`
	RedisAddr          string `mapstructure:"redis_addr""`
	RedisUser          string `mapstructure:"redis_user""`
	RedisPassword      string `mapstructure:"redis_password""`
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

	os.Setenv("redis_addr", goConfig.RedisAddr)
	os.Setenv("redis_user", goConfig.RedisUser)
	os.Setenv("redis_password", goConfig.RedisPassword)

	return goConfig
}
