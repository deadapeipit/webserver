package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"webserver/entity"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type JWTHandler struct{}

func InstallJWTHandler(r *mux.Router) {
	api := JWTHandler{}
	r.HandleFunc("/jwt", api.JWTHandler)

}

type JWTHandlerInterface interface {
	JWTHandler(w http.ResponseWriter, r *http.Request)
}

func NewJWTHandler() JWTHandlerInterface {
	return &JWTHandler{}
}

func (h *JWTHandler) JWTHandler(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)

	switch r.Method {
	case http.MethodGet:
		h.getJWTToken(w, r)
	}
}

// getJWTToken
// Method: GET
// Example: localhost/jwt
func (h *JWTHandler) getJWTToken(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user entity.UserLogin
	if err := decoder.Decode(&user); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}
	claims := entity.MyClaims{
		Iat: int(time.Now().UnixMilli()),
		Exp: int(time.Now().Add(time.Second * time.Duration(60)).UnixMilli()),
		Uid: user.Username,
		Pwd: user.Password,
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
