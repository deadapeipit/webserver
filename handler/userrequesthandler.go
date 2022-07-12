package handler

import (
	"encoding/json"
	"net/http"
	"webserver/entity"

	"github.com/gorilla/mux"
)

type UserRequestHandler struct{}

func InstallUserRequestHandler(r *mux.Router) {
	api := UserRequestHandler{}
	r.HandleFunc("/user", api.UsersRequestHandler)

}

type UserRequestHandlerInterface interface {
	UsersRequestHandler(w http.ResponseWriter, r *http.Request)
}

func NewUserRequestHandler() UserRequestHandlerInterface {
	return &UserRequestHandler{}
}

func (h *UserRequestHandler) UsersRequestHandler(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)

	switch r.Method {
	case http.MethodGet:
		h.getUserRequestHandler(w, r)
	}
}

// getUserRequestHandler
// Method: GET
// Example: localhost/user
func (h *UserRequestHandler) getUserRequestHandler(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("https://random-data-api.com/api/users/random_user?size=10")
	if err != nil {
		WriteJsonResp(w, StatusError, err.Error())
		return
	}
	decoder := json.NewDecoder(res.Body)

	// cara tanpa tampung ke dalam struct sementara:
	// var tesData []map[string]interface{}
	// if err := decoder.Decode(&tesData); err != nil {
	// 	w.Write([]byte("error decoding json body"))
	// 	return
	// }
	// var resp []entity.UserResponse
	// var m []map[string]interface{}
	// for _, data := range tesData {
	// 	tempUser := make(map[string]interface{})
	// 	var floatID float64 = data["id"].(float64)
	// 	var id int = int(floatID)
	// 	tempUser["id"] = id
	// 	tempUser["uid"] = data["uid"].(string)
	// 	tempUser["first_name"] = data["first_name"].(string)
	// 	tempUser["last_name"] = data["last_name"].(string)
	// 	tempUser["user_name"] = data["username"].(string)
	// 	tempUser["address"] = data["address"].(map[string]interface{})
	// 	m = append(m, tempUser)
	// }
	// jsonString, _ := json.Marshal(m)
	// err = json.Unmarshal(jsonString, &resp)
	// if err != nil {
	// 	w.Write([]byte("error decoding json body"))
	// 	return
	// }
	var user []entity.UserRequest
	if err := decoder.Decode(&user); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}
	resp := entity.ToArrayUserResponse(user)
	WriteJsonResp(w, StatusSuccess, resp)
}
