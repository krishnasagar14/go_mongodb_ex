package main

import (
	// "bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	// "github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"go_mongodb_ex/handlers"
	// "go_mongodb_ex/proto"
)

func PrepareRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/assignment/user", handlers.GetUserHandler).Methods("GET")
	router.HandleFunc("/assignment/user", handlers.UpdateUserHandler).Methods("PATCH")
	router.HandleFunc("/assignment/user", handlers.CreateUserHandler).Methods("POST")
	return router
}

func TestGetEmptyUserAPI(t *testing.T) {
	request, _ := http.NewRequest("GET", "/assignment/user", nil)
	response := httptest.NewRecorder()
	PrepareRouter().ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "Correct response code 400 expected")
}

func TestPostEmptyUserAPI(t *testing.T) {
	request, _ := http.NewRequest("POST", "/assignment/user", nil)
	response := httptest.NewRecorder()
	PrepareRouter().ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "Correct response code 400 expected")
}

// func TestPostUserAPI(t *testing.T) {
// 	proto_request := &server_proto.CreateUserRequest{
// 		FirstName:   "Krishna",
// 		LastName:    "Subhedarpage",
// 		Email:       "krishnasagar.subhedarpage@ahs.com",
// 		Designation: "SE3",
// 	}
// 	req, _ := proto.Marshal(proto_request)
// 	act_request, _ := http.NewRequest("POST", "/assignment/user", bytes.NewReader(req))
// 	act_request.Header.Set("Content-Type", "application/x-binary")
// 	response := httptest.NewRecorder()
// 	PrepareRouter().ServeHTTP(response, act_request)
// 	assert.Equal(t, 201, response.Code, "Correct response code 201 expected")
// }
