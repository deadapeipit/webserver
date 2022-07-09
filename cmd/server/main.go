package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"webserver/database"
	"webserver/handler"

	"github.com/gorilla/mux"
)

// Replace with your own connection parameters
var sqlserver = "localhost"
var sqlport = 1433
var sqldbName = "tes"
var sqluser = "david"
var sqlpassword = "david"

func main() {
	// Create connection string
	// connString := fmt.Sprintf("server=%s;database=%s;port=%d;trusted_connection=yes",
	// 	sqlserver, sqldbName, sqlport)

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		sqluser, sqlpassword, sqlserver, sqlport, sqldbName)

	sql := database.NewSqlConnection(connString)
	handler.SqlConnect = sql
	r := mux.NewRouter()

	userHandler := handler.NewUserHandler()
	orderHandler := handler.NewOrderHandler()
	r.HandleFunc("/users", userHandler.UsersHandler)
	r.HandleFunc("/users/{id}", userHandler.UsersHandler)
	r.HandleFunc("/orders", orderHandler.OrdersHandler)
	r.HandleFunc("/orders/{id}", orderHandler.OrdersHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Printf("Link: http://%s \n", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}
