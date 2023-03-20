package server

import (
	"fmt"
	"net/http"
)

func MakeLaunchpad(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Make Launchpad")
}

func GetLaunchpadData(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("id")
	fmt.Println(hash)
	// hash값을 통해서 데이터를 가져 올 예정

	if hash == "" {
		http.Error(w, "id값이 빈 문자열 입니다.", http.StatusForbidden)
		return
	}
	//
	//w.WriteHeader(http.StatusTemporaryRedirect)
	http.Error(w, "잘못된 요청", http.StatusInternalServerError)

}
