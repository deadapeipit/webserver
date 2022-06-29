package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"webserver/entity"
	"webserver/handler"
	"webserver/service"

	"github.com/gorilla/mux"
)

var PORT = ":8888"

func main() {
	r := mux.NewRouter()

	userHandler := handler.NewUserHandler()
	//fmt.Println("Hello World")
	//http.HandleFunc("/greet", greet)
	r.HandleFunc("/register", userRegister).Methods("POST")
	r.HandleFunc("/users", userHandler.UsersHandler)
	r.HandleFunc("/users/{id}", userHandler.UsersHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	//http.ListenAndServe(PORT, nil)
	log.Fatal(srv.ListenAndServe())
}

// func greet(w http.ResponseWriter, r *http.Request) {
// 	msg := "Hello world"
// 	fmt.Fprint(w, msg)
// }

func userRegister(w http.ResponseWriter, r *http.Request) {
	userSvc := service.NewUserService()
	// newUser := &entity.User{
	// 	Id:        1,
	// 	Username:  "david123",
	// 	Email:     "david123@gmail.com",
	// 	Password:  "Passdav!d",
	// 	Age:       17,
	// 	CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(),
	// }

	decoder := json.NewDecoder(r.Body)
	var newUser entity.User
	if err := decoder.Decode(&newUser); err != nil {
		w.WriteHeader(201)
		w.Write([]byte("error decoding json body"))
		return
	}

	if user, err := userSvc.Register(&newUser); err != nil {
		fmt.Printf("Error when register user: %+v \n", err)
		w.WriteHeader(201)
		w.Write([]byte("Error when register user"))
		return
	} else {
		m, err := json.Marshal(user)
		if err != nil {
			fmt.Printf("Error when register user: %+v \n", err)
			w.WriteHeader(201)
			w.Write([]byte("Error when register user"))
		}

		fmt.Printf("Success register user: %+v \n", user)
		fmt.Println("----------------------------------")
		w.Header().Add("Content-Type", "application/json")
		w.Write(m)
	}
}
