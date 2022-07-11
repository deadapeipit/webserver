package handler

import (
	"encoding/json"
	"net/http"
	"webserver/database"
)

type helper struct {
	Tesdb database.DatabaseIface
}

var Helper helper

type response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

const (
	StatusSuccess int = 0
	StatusError   int = 1
)

func WriteJsonResp(w http.ResponseWriter, status int, obj interface{}) {

	resp := response{
		Status: status,
		Data:   obj,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
