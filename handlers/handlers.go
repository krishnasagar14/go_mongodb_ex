package handlers

import (
	b64 "encoding/base64"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"

	"go_mongodb_ex/db"
	"go_mongodb_ex/proto"
)

func GetUserHandler(resp http.ResponseWriter, req *http.Request) {
	proto_body_raw, _ := req.URL.Query()["proto_body"]
	if proto_body_raw == nil || len(proto_body_raw[0]) < 1 {
		log.Println("Proto body is missing")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	proto_body, err := b64.URLEncoding.DecodeString(proto_body_raw[0])
	if err != nil {
		log.Println("Proto body decoding problem found: %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	request := &server_proto.GetUserRequest{}
	proto.Unmarshal(proto_body, request)
	user_id := request.GetUserId()
	if user_id == "" {
		log.Println("No user id found")
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	user_record, err := db.GetUserViaId(user_id)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	res := &server_proto.UserDetailsResponse{
		UserId:      user_record.Id.Hex(),
		EmployeeId:  user_record.EmployeeId,
		FirstName:   user_record.FirstName,
		LastName:    user_record.LastName,
		Email:       user_record.Email,
		Designation: user_record.Designation,
	}
	response, err := proto.Marshal(res)
	if err != nil {
		log.Println("Unable to marshal response for get user: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.Write(response)

}

func UpdateUserHandler(resp http.ResponseWriter, req *http.Request) {
	request := &server_proto.UpdateUserRequest{}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Unable to read request message for update user: %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	proto.Unmarshal(data, request)
	updated_user, err := db.UpdateUser(request)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	res := &server_proto.UserDetailsResponse{
		UserId:      updated_user.Id.Hex(),
		EmployeeId:  updated_user.EmployeeId,
		FirstName:   updated_user.FirstName,
		LastName:    updated_user.LastName,
		Email:       updated_user.Email,
		Designation: updated_user.Designation,
	}
	response, err := proto.Marshal(res)
	if err != nil {
		log.Println("Unable to marshal response for update user: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write(response)
}

func CreateUserHandler(resp http.ResponseWriter, req *http.Request) {
	request := &server_proto.CreateUserRequest{}
	if req.Body == nil {
		log.Println("Unable to read request message for create user")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Unable to read request message for create user: %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	proto.Unmarshal(data, request)
	new_id, emp_id := db.InsertNewUser(request)
	if new_id == "" || emp_id == "" {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := &server_proto.UserDetailsResponse{
		UserId:      new_id,
		EmployeeId:  emp_id,
		FirstName:   request.GetFirstName(),
		LastName:    request.GetLastName(),
		Email:       request.GetEmail(),
		Designation: request.GetDesignation(),
	}
	response, err := proto.Marshal(res)
	if err != nil {
		log.Println("Unable to marshal response for create user: %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusCreated)
	resp.Write(response)
}
