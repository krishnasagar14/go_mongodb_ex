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
	protoBodyRaw, _ := req.URL.Query()["proto_body"]
	if protoBodyRaw == nil || len(protoBodyRaw[0]) < 1 {
		log.Println("Proto body is missing")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	protoBody, err := b64.URLEncoding.DecodeString(protoBodyRaw[0])
	if err != nil {
		log.Println("Proto body decoding problem found: %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	request := &server_proto.GetUserRequest{}
	proto.Unmarshal(protoBody, request)
	user_id := request.GetUserId()
	if user_id == "" {
		log.Println("No user id found")
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	userRecord, err := db.GetUserViaId(user_id)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	res := &server_proto.UserDetailsResponse{
		UserId:      userRecord.Id.Hex(),
		EmployeeId:  userRecord.EmployeeId,
		FirstName:   userRecord.FirstName,
		LastName:    userRecord.LastName,
		Email:       userRecord.Email,
		Designation: userRecord.Designation,
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
	if req.Body == nil {
		log.Println("Unable to read request message for update user")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Unable to read request message for update user: %v", err)
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	proto.Unmarshal(data, request)
	updatedUser, err := db.UpdateUser(request)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	res := &server_proto.UserDetailsResponse{
		UserId:      updatedUser.Id.Hex(),
		EmployeeId:  updatedUser.EmployeeId,
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		Email:       updatedUser.Email,
		Designation: updatedUser.Designation,
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
	newId, empId := db.InsertNewUser(request)
	if newId == "" || empId == "" {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	res := &server_proto.UserDetailsResponse{
		UserId:      newId,
		EmployeeId:  empId,
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
