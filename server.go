package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go_mongodb_ex/db"
	"go_mongodb_ex/handlers"
	"go_mongodb_ex/middlewares"
)

func main() {
	db.ConnectDB("local_db")

	port_no := 9000
	fmt.Println("Starting server on port:", port_no)
	router := mux.NewRouter().StrictSlash(true)

	router.Use(mw.LoggingMiddleware)

	router.HandleFunc("/assignment/user", handlers.GetUserHandler).Methods("GET")
	router.HandleFunc("/assignment/user", handlers.UpdateUserHandler).Methods("PATCH")
	router.HandleFunc("/assignment/user", handlers.CreateUserHandler).Methods("POST")

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port_no),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
