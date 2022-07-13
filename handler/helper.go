package handler

import (
	"encoding/json"
	"net/http"
	"webserver/database"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type helper struct {
	Tesdb database.DatabaseIface
}

var Helper helper
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

type response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

const (
	StatusSuccess int = 0
	StatusError   int = 1
)

func EncryptPassword(pwd string) (string, error) {
	// Hashing the password with the default cost of 10
	securePassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(securePassword), nil
}

func WriteJsonResp(w http.ResponseWriter, status int, obj interface{}) {

	resp := response{
		Status: status,
		Data:   obj,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
