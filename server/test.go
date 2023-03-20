package server

import (
	"fmt"
	"net/http"
)

type GetTest struct {
	Name string `json:"name"`
}

func NewGetTest() Handler {
	return &GetTest{
		Name: "test",
	}
}

func (t GetTest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("들어옴")
}
