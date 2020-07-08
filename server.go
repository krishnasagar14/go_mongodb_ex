package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go_mongodb_ex/db"
	"go_mongodb_ex/handlers"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.RequestURI)
		next.ServeHTTP(resp, req)
	})
}

func main() {
	DB.ConnectDB()

	port_no := 9000
	fmt.Println("Starting server on port:", port_no)
	router := mux.NewRouter().StrictSlash(true)

	router.Use(LoggingMiddleware)

	router.HandleFunc("/assignment/user", handlers.GetUserHandler).Methods("GET").Queries("proto_body", "{proto_body}")
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
