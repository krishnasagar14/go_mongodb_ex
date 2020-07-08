package handlers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"

	"go_mongodb_ex/db"
	"go_mongodb_ex/proto"
)

func GetUserHandler(resp http.ResponseWriter, req *http.Request) {

}

func UpdateUserHandler(resp http.ResponseWriter, req *http.Request) {
	request := &server_proto.UpdateUserRequest{}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Unable to read request message for update user: %v", err)
	}
	proto.Unmarshal(data, request)
	updated_user := DB.UpdateUser(request)
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
	}
	resp.Write(response)
}

func CreateUserHandler(resp http.ResponseWriter, req *http.Request) {
	request := &server_proto.CreateUserRequest{}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("Unable to read request message for create user: %v", err)
	}
	proto.Unmarshal(data, request)
	new_id, emp_id := DB.InsertNewUser(request)
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
	}
	resp.Write(response)
}
