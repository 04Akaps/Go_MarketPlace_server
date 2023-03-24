package init

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"net"
	"os"
)

func initOAuth(envData EnvData) {
	authKey := os.Getenv("AUTH_KEY") // bashrc에 export한 값을 가져 온다.

	if len(authKey) == 0 {
		authKey = envData.AuthKey
	}

	var baseUri string

	// 현재 IP주소를 가져오는 코드
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					baseUri = ipnet.IP.String()
				}
			}
		}
	}

	os.Setenv("BASE_URL", baseUri)
	//baseUri = strings.Join([]string{baseUri, "/auth"}, "")

	//192.168.219.101/auth -> 이런 형태가 된다.

	gomniauth.SetSecurityKey(authKey)
	gomniauth.WithProviders(
		google.New(envData.GoogleAuthId, envData.GoogleAuthPassword, "http://localhost/auth/callback/google"),
	)
}
