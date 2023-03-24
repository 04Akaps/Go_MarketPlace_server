package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"goServer/customError"
	"goServer/paseto"
	"net/http"
)

const (
	OAUTH_ACTION_INVALID         = "OAuth에 적합한 URL가 아님"
	OAUTH_PROVIDER_NOT_FOUND     = "지원하지 않는 Provider형태"
	OAUTH_CREDENTIALS_ERROR      = "OAuth의 Auth에러"
	OAUTH_GET_USER_ERROR         = "불분명한 USER"
	OAUTH_GET_CALLBACK_URL_ERROR = "callBack URL을 찾지 못하였습니다."
	PASETO_TOKEN_CREATE_FAILED   = "PaseToken 생성 실패"
)

type AuthController struct {
	ctx     context.Context
	channel chan error
	paseto  paseto.PasetoInterface
}

type AuthInterface interface {
	Auth(http.ResponseWriter, *http.Request)
	Logout(http.ResponseWriter, *http.Request)
}

func NewAuthController(channel chan error, paseto paseto.PasetoInterface) AuthInterface {
	context := context.Background()
	return &AuthController{
		ctx:     context,
		channel: channel,
		paseto:  paseto,
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
			controller.channel <- err
			customError.NewHandlerError(w, OAUTH_PROVIDER_NOT_FOUND, http.StatusBadRequest)
			return
		}

		loginUrl, err := provider.GetBeginAuthURL(nil, nil)

		if err != nil {
			controller.channel <- err
			customError.NewHandlerError(w, OAUTH_GET_CALLBACK_URL_ERROR, http.StatusInternalServerError)
			return
		}

		// Location을 이동 시켜야 하나...?

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":

		provider, err := gomniauth.Provider(provider)

		if err != nil {
			controller.channel <- err
			customError.NewHandlerError(w, OAUTH_PROVIDER_NOT_FOUND, http.StatusBadRequest)
			return
		}

		credentials, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			controller.channel <- err
			customError.NewHandlerError(w, OAUTH_CREDENTIALS_ERROR, http.StatusInternalServerError)
			return
		}

		user, err := provider.GetUser(credentials)

		if err != nil {
			controller.channel <- err
			customError.NewHandlerError(w, OAUTH_GET_USER_ERROR, http.StatusInternalServerError)
			return
		}

		newToken, err := controller.paseto.CreateToken(user.Name())

		if err != nil {
			controller.channel <- err
			customError.NewHandlerError(w, PASETO_TOKEN_CREATE_FAILED, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: newToken,
			Path:  "/",
		})

		w.Header().Set("Location", "http://localhost:3000/")
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		customError.NewHandlerError(w, OAUTH_ACTION_INVALID, 200)
	}

}

func (controller *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout")

	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.Header().Set("Location", "http://localhost:3000")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
