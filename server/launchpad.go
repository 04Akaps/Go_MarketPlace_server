package server

import (
	"fmt"
	"net/http"
)

func MakeLaunchpad(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Make Launchpad")
}

func GetLaunchpadData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Launchpad")
	//hash := r.URL.Query().Get("id")
	// hash값을 통해서 데이터를 가져 올 예정
}
