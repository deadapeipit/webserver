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

// // Replace with your own connection parameters
// var sqlserver = "localhost"
// var sqlport = 1433
// var sqldbName = "tes"
// var sqluser = "david"
// var sqlpassword = "david"

func main() {
	handler.GetConfig()
	// Create connection string
	// connString := fmt.Sprintf("server=%s;database=%s;port=%d;trusted_connection=yes",
	// 	sqlserver, sqldbName, sqlport)

	// connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
	// 	sqluser, sqlpassword, sqlserver, sqlport, sqldbName)

	connString := handler.Config.TesConnectionString
	//fmt.Printf("Connection string: %s \n", connString)

	sql := database.NewSqlConnection(connString)
	handler.Helper.Tesdb = sql
	defer handler.Helper.Tesdb.CloseConnection()

	r := mux.NewRouter()
	handler.InstallOrderAPI(r)
	handler.InstalUsersHandler(r)
	handler.InstallUserRequestHandler(r)
	r.Use(handler.SecureMiddleware)

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
