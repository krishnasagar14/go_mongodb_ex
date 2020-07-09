package main_routes

import (
	"github.com/gorilla/mux"

	"go_mongodb_ex/handlers"
)

func RegisterRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/assignment/user", handlers.GetUserHandler).Methods("GET")
	router.HandleFunc("/assignment/user", handlers.UpdateUserHandler).Methods("PATCH")
	router.HandleFunc("/assignment/user", handlers.CreateUserHandler).Methods("POST")
	return router
}
