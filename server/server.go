package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go_mongodb_ex/db"
	"go_mongodb_ex/middlewares"
	"go_mongodb_ex/routers"
)

func main() {
	err := db.ConnectDB("local_db")
	if err != nil {
		os.Exit(1)
	}

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
