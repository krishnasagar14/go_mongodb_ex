package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go_mongodb_ex/db"
	"go_mongodb_ex/middlewares"
	"go_mongodb_ex/routers"
)

func main() {
	db.ConnectDB("local_db")

	portNo := 9000
	fmt.Println("Starting server on port:", portNo)
	router := main_routes.RegisterRouter()
	router.Use(mw.LoggingMiddleware)

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", portNo),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
