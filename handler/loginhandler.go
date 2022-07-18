package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"webserver/entity"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct{}

func InstallLoginHandler(r *mux.Router) {
	api := LoginHandler{}
	r.HandleFunc("/login", api.login)

}

type LoginHandlerInterface interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
}

func NewJWTHandler() LoginHandlerInterface {
	return &LoginHandler{}
}

func (h *LoginHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)

	switch r.Method {
	case http.MethodGet:
		h.login(w, r)
	}
}

// getJWTToken
// Method: GET
// Example: localhost/login
func (h *LoginHandler) login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	decoder := json.NewDecoder(r.Body)
	var user entity.UserLogin
	if err := decoder.Decode(&user); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}
	u, err := Helper.Tesdb.Login(ctx, user.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteJsonResp(w, StatusError, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		WriteJsonResp(w, StatusError, "invalid username / password")
		return
	}
	claims := entity.MyClaims{
		Iat: int(time.Now().UnixMilli()),
		Exp: int(time.Now().Add(time.Second * time.Duration(60)).UnixMilli()),
		Uid: user.Username,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	retVal, err := token.SignedString([]byte(Config.SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteJsonResp(w, StatusError, "BAD_REQUEST")
		return
	}
	WriteJsonResp(w, StatusSuccess, retVal)
}
