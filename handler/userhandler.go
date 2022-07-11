package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"webserver/entity"

	"github.com/gorilla/mux"
)

type UserHandler struct{}

func InstalUsersHandler(r *mux.Router) {
	api := UserHandler{}
	r.HandleFunc("/users", api.UsersHandler)
	r.HandleFunc("/users/{id}", api.UsersHandler)
}

type UserHandlerInterface interface {
	UsersHandler(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler() UserHandlerInterface {
	return &UserHandler{}
}

func (h *UserHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	switch r.Method {
	case http.MethodGet:
		if id != "" { // get by id
			getUsersByIDHandler(w, r, id)
		} else { // get all
			h.getUsersHandler(w, r)
		}
	case http.MethodPost:
		if id != "" {
			updateUserHandler(w, r, id)
		} else {
			createUsersHandler(w, r)
		}
	case http.MethodPut:
		updateUserHandler(w, r, id)
	case http.MethodDelete:
		deleteUserHandler(w, r, id)
	}
}

// getUsersHandler
// Method: GET
// Example: localhost/users
func (h *UserHandler) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	users, err := Helper.Tesdb.GetUsers(ctx)
	if err != nil {
		WriteJsonResp(w, StatusError, err.Error())
		return
	}
	WriteJsonResp(w, StatusSuccess, users)
}

// getUsersByIDHandler
// Method: GET
// Example: localhost/users/1
func getUsersByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if idInt, err := strconv.Atoi(id); err == nil {
		users, err := Helper.Tesdb.GetUserByID(ctx, idInt)
		if err != nil {
			WriteJsonResp(w, StatusError, err.Error())
			return
		}
		if idInt == users.Id {
			WriteJsonResp(w, StatusError, "Data not exists")
			return
		}
		WriteJsonResp(w, StatusSuccess, users)
	}
}

// createUsersHandler
// Method: POST
// Example: localhost/users
// JSON Body:
// {
//		"id": 1,
//		"user_name": "user1",
//		"email": "user@email.com",
//		"password": "password1",
//		"age": 22
// }
func createUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	decoder := json.NewDecoder(r.Body)
	var user entity.User
	if err := decoder.Decode(&user); err != nil {
		w.Write([]byte("error decoding json body"))
		return
	}

	users, err := Helper.Tesdb.CreateUser(ctx, user)
	if err != nil {
		WriteJsonResp(w, StatusError, err.Error())
		return
	}
	WriteJsonResp(w, StatusSuccess, users)
}

// updateUserHandler
// Method: POST / PUT
// Example: localhost/users/1
// JSON Body:
// {
//		"id": 1,
//		"user_name": "user1",
//		"email": "user@email.com",
//		"password": "password1",
//		"age": 22
// }
func updateUserHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()

	if id != "" { // get by id
		if idInt, err := strconv.Atoi(id); err == nil {
			if users, err := Helper.Tesdb.GetUserByID(ctx, idInt); err != nil {
				WriteJsonResp(w, StatusError, err.Error())
				return
			} else if idInt == users.Id {
				WriteJsonResp(w, StatusError, "Data not exists")
				return
			} else {
				decoder := json.NewDecoder(r.Body)
				var user entity.User
				if err := decoder.Decode(&user); err != nil {
					w.Write([]byte("error decoding json body"))
					return
				}

				users, err := Helper.Tesdb.UpdateUser(ctx, idInt, user)
				if err != nil {
					WriteJsonResp(w, StatusError, err.Error())
					return
				}
				WriteJsonResp(w, StatusSuccess, users)
			}
		}
	}
}

// deleteUserHandler
// Method: DELETE
// Example: localhost/users/1
func deleteUserHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if id != "" { // get by id
		if idInt, err := strconv.Atoi(id); err == nil {
			if users, err := Helper.Tesdb.GetUserByID(ctx, idInt); err != nil {
				WriteJsonResp(w, StatusError, err.Error())
				return
			} else if idInt == users.Id {
				WriteJsonResp(w, StatusError, "Data not exists")
				return
			} else {
				users, err := Helper.Tesdb.DeleteUser(ctx, idInt)
				if err != nil {
					WriteJsonResp(w, StatusError, err.Error())
					return
				}
				WriteJsonResp(w, StatusSuccess, users)
			}
		}
	}
}
