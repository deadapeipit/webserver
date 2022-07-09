package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"webserver/entity"

	"github.com/gorilla/mux"
)

type UserHandlerInterface interface {
	UsersHandler(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	//postgrespool *pgxpool.Pool
}

func NewUserHandler() UserHandlerInterface {
	//return &UserHandler{postgrespool: postgrespool}
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
	users, err := SqlConnect.GetUsers(ctx)
	if err != nil {
		writeJsonResp(w, statusError, err.Error())
		return
	}
	writeJsonResp(w, statusSuccess, users)
}

// getUsersByIDHandler
// Method: GET
// Example: localhost/users/1
func getUsersByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()
	if idInt, err := strconv.Atoi(id); err == nil {
		users, err := SqlConnect.GetUserByID(ctx, idInt)
		if err != nil {
			writeJsonResp(w, statusError, err.Error())
			return
		}
		if idInt == users.Id {
			writeJsonResp(w, statusError, "Data not exists")
			return
		}
		writeJsonResp(w, statusSuccess, users)
	}
}

// createUsersHandler
// Method: POST
// Example: localhost/users
// JSON Body:
// {
//		"id": 1,
//		"user_name": "user1",
//		"email": "user@email.com"
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

	users, err := SqlConnect.CreateUser(ctx, user)
	if err != nil {
		writeJsonResp(w, statusError, err.Error())
		return
	}
	writeJsonResp(w, statusSuccess, users)
}

// updateUserHandler
// Method: POST / PUT
// Example: localhost/users/1
// JSON Body:
// {
//		"id": 1,
//		"user_name": "user1",
//		"email": "user@email.com"
//		"password": "password1",
//		"age": 22
// }
func updateUserHandler(w http.ResponseWriter, r *http.Request, id string) {
	ctx := context.Background()

	if id != "" { // get by id
		if idInt, err := strconv.Atoi(id); err == nil {
			if users, err := SqlConnect.GetUserByID(ctx, idInt); err != nil {
				writeJsonResp(w, statusError, err.Error())
				return
			} else if idInt == users.Id {
				writeJsonResp(w, statusError, "Data not exists")
				return
			} else {
				decoder := json.NewDecoder(r.Body)
				var user entity.User
				if err := decoder.Decode(&user); err != nil {
					w.Write([]byte("error decoding json body"))
					return
				}

				users, err := SqlConnect.UpdateUser(ctx, idInt, user)
				if err != nil {
					writeJsonResp(w, statusError, err.Error())
					return
				}
				writeJsonResp(w, statusSuccess, users)
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
			if users, err := SqlConnect.GetUserByID(ctx, idInt); err != nil {
				writeJsonResp(w, statusError, err.Error())
				return
			} else if idInt == users.Id {
				writeJsonResp(w, statusError, "Data not exists")
				return
			} else {
				users, err := SqlConnect.DeleteUser(ctx, idInt)
				if err != nil {
					writeJsonResp(w, statusError, err.Error())
					return
				}
				writeJsonResp(w, statusSuccess, users)
			}
		}
	}
}
