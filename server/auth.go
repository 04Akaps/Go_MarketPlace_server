package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"goServer/customError"
	"net/http"
)

const (
	OAUTH_ACTION_INVALID         = "OAuth에 적합한 URL가 아님"
	OAUTH_PROVIDER_NOT_FOUND     = "지원하지 않는 Provider형태"
	OAUTH_CREDENTIALS_ERROR      = "OAuth의 Auth에러"
	OAUTH_GET_USER_ERROR         = "불분명한 USER"
	OAUTH_GET_CALLBACK_URL_ERROR = "callBack URL을 찾지 못하였습니다."
)

type AuthController struct {
	ctx context.Context
}

type AuthInterface interface {
	Auth(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
}

func NewAuthController() AuthInterface {
	context := context.Background()
	return &AuthController{
		ctx: context,
	}
}

func (controller *AuthController) Auth(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	action := vars["action"]
	provider := vars["provider"]

	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			customError.NewHandlerError(w, OAUTH_PROVIDER_NOT_FOUND, http.StatusBadRequest)
			return
		}

		loginUrl, err := provider.GetBeginAuthURL(nil, nil)

		if err != nil {
			customError.NewHandlerError(w, OAUTH_GET_CALLBACK_URL_ERROR, http.StatusInternalServerError)
			return
		}

		// Location을 이동 시켜야 하나...?

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":

		provider, err := gomniauth.Provider(provider)

		if err != nil {
			customError.NewHandlerError(w, OAUTH_PROVIDER_NOT_FOUND, http.StatusBadRequest)
			return
		}

		credentials, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			customError.NewHandlerError(w, OAUTH_CREDENTIALS_ERROR, http.StatusInternalServerError)
			return
		}

		user, err := provider.GetUser(credentials)

		if err != nil {
			customError.NewHandlerError(w, OAUTH_GET_USER_ERROR, http.StatusInternalServerError)
			return
		}

		fmt.Println(user.Name())

		// Paseto를 통해 JWT토큰을 만들지 말지 고민 중
		// 만든뒤에 middleWare에 추가하여, Token을 검증 하는 것이 좋을까 고민 중

		w.Header().Set("Location", "http://localhost:3000/")
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		customError.NewHandlerError(w, OAUTH_ACTION_INVALID, 200)
	}

}

func (controller *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout")
}
