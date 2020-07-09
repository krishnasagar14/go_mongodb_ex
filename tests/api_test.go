package main

import (
	"bytes"
	b64 "encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"go_mongodb_ex/db"
	"go_mongodb_ex/handlers"
	"go_mongodb_ex/proto"
)

const DB_NAME string = "test_local_db"

var test_router *mux.Router

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
	test_router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code, "Correct response code 400 expected")
}

func TestPostEmptyUserAPI(t *testing.T) {
	request, _ := http.NewRequest("POST", "/assignment/user", nil)
	response := httptest.NewRecorder()
	test_router.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code, "Correct response code 400 expected")
}

func make_post_request() *httptest.ResponseRecorder {
	proto_request := &server_proto.CreateUserRequest{
		FirstName:   "Krishna",
		LastName:    "Subhedarpage",
		Email:       "krishnasagar.subhedarpage@ahs.com",
		Designation: "SE3",
	}
	req, _ := proto.Marshal(proto_request)
	act_request, _ := http.NewRequest("POST", "/assignment/user", bytes.NewReader(req))
	act_request.Header.Set("Content-Type", "application/x-binary")
	response := httptest.NewRecorder()
	test_router.ServeHTTP(response, act_request)
	return response
}

func TestPostUserAPI(t *testing.T) {
	response := make_post_request()
	assert.Equal(t, http.StatusCreated, response.Code, "Correct response code 201 expected")
}

func TestPatchUserAPI(t *testing.T) {
	response := make_post_request()
	data, _ := ioutil.ReadAll(response.Body)
	user_details := &server_proto.UserDetailsResponse{}
	proto.Unmarshal(data, user_details)

	proto_req := &server_proto.UpdateUserRequest{
		Email:  "krishnasagar@ahs.com",
		UserId: user_details.GetUserId(),
	}
	req, _ := proto.Marshal(proto_req)
	act_request, _ := http.NewRequest("PATCH", "/assignment/user", bytes.NewReader(req))
	act_request.Header.Set("Content-Type", "application/x-binary")
	patch_response := httptest.NewRecorder()
	test_router.ServeHTTP(patch_response, act_request)
	assert.Equal(t, http.StatusOK, patch_response.Code, "Correct response code 200 expected")
}

func TestGetUserAPI(t *testing.T) {
	response := make_post_request()
	data, _ := ioutil.ReadAll(response.Body)
	user_details := &server_proto.UserDetailsResponse{}
	proto.Unmarshal(data, user_details)

	request := &server_proto.GetUserRequest{
		UserId: user_details.GetUserId(),
	}
	p_body, _ := proto.Marshal(request)
	proto_body := b64.URLEncoding.EncodeToString(p_body)
	api_url := "/assignment/user"
	u, _ := url.Parse(api_url)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("proto_body", proto_body)
	u.RawQuery = q.Encode()

	act_request, _ := http.NewRequest("GET", u.String(), nil)
	get_response := httptest.NewRecorder()
	test_router.ServeHTTP(get_response, act_request)
	assert.Equal(t, http.StatusOK, get_response.Code, "Correct response code 200 expected")
}

func TestMain(m *testing.M) {
	db.ConnectDB(DB_NAME)
	test_router = PrepareRouter()
	m.Run()
	db.DropDB(DB_NAME)
	os.Exit(0)
}
