package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	//"github.com/gorilla/mux"

	"github.com/go_mongodb_ex/proto/server_proto"
)

func GetUserHandler(resp http.ResponseWriter, req *http.Request) {

}

func UpdateUserHandler(resp http.ResponseWriter, req *http.Request) {

}

func CreateUserHandler(resp http.ResponseWriter, req *http.Request) {
	request := &server_proto.CreateUserRequest{}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("Unable to read request message for create user: %v", err)
	}
	proto.Unmarshal(data, request)
	res := &server_proto.UserDetailsResponse{
		user_id:     "1",
		employee_id: "ftdr-1",
		first_name:  request.GetFirstName(),
		last_name:   request.GetLastName(),
		email:       request.GetEmail(),
		designation: request.GetDesignation(),
	}
	response, err := proto.Marshal(res)
	if err != nil {
		log.Fatalf("Unable to marshal response for create user: %v", err)
	}
	resp.Write(response)
}
