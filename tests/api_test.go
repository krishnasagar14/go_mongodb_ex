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
	"go_mongodb_ex/proto"
	"go_mongodb_ex/routers"
)

const DbName string = "test_local_db"

var testRouter *mux.Router

func TestGetEmptyUserAPI(t *testing.T) {
	request, _ := http.NewRequest("GET", "/assignment/user", nil)
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code, "Correct response code 400 expected")
}

func TestPostEmptyUserAPI(t *testing.T) {
	request, _ := http.NewRequest("POST", "/assignment/user", nil)
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code, "Correct response code 400 expected")
}

func makePostRequest() *httptest.ResponseRecorder {
	protoRequest := &server_proto.CreateUserRequest{
		FirstName:   "Krishna",
		LastName:    "Subhedarpage",
		Email:       "krishnasagar.subhedarpage@ahs.com",
		Designation: "SE3",
	}
	req, _ := proto.Marshal(protoRequest)
	actRequest, _ := http.NewRequest("POST", "/assignment/user", bytes.NewReader(req))
	actRequest.Header.Set("Content-Type", "application/x-binary")
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, actRequest)
	return response
}

func TestPostUserAPI(t *testing.T) {
	response := makePostRequest()
	assert.Equal(t, http.StatusCreated, response.Code, "Correct response code 201 expected")
}

func TestPatchUserAPI(t *testing.T) {
	response := makePostRequest()
	data, _ := ioutil.ReadAll(response.Body)
	userDetails := &server_proto.UserDetailsResponse{}
	proto.Unmarshal(data, userDetails)

	protoReq := &server_proto.UpdateUserRequest{
		Email:  "krishnasagar@ahs.com",
		UserId: userDetails.GetUserId(),
	}
	req, _ := proto.Marshal(protoReq)
	actRequest, _ := http.NewRequest("PATCH", "/assignment/user", bytes.NewReader(req))
	actRequest.Header.Set("Content-Type", "application/x-binary")
	patchResponse := httptest.NewRecorder()
	testRouter.ServeHTTP(patchResponse, actRequest)
	assert.Equal(t, http.StatusOK, patchResponse.Code, "Correct response code 200 expected")
}

func TestGetUserAPI(t *testing.T) {
	response := makePostRequest()
	data, _ := ioutil.ReadAll(response.Body)
	userDetails := &server_proto.UserDetailsResponse{}
	proto.Unmarshal(data, userDetails)

	request := &server_proto.GetUserRequest{
		UserId: userDetails.GetUserId(),
	}
	pBody, _ := proto.Marshal(request)
	protoBody := b64.URLEncoding.EncodeToString(pBody)
	apiUrl := "/assignment/user"
	u, _ := url.Parse(apiUrl)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("proto_body", protoBody)
	u.RawQuery = q.Encode()

	actRequest, _ := http.NewRequest("GET", u.String(), nil)
	getResponse := httptest.NewRecorder()
	testRouter.ServeHTTP(getResponse, actRequest)
	assert.Equal(t, http.StatusOK, getResponse.Code, "Correct response code 200 expected")
}

func TestMain(m *testing.M) {
	db.ConnectDB(DbName)
	testRouter = main_routes.RegisterRouter()
	m.Run()
	db.DropDB()
	os.Exit(0)
}
