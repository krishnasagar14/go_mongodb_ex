package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/golang/protobuf/proto"
	"go_mongodb_ex/proto"
)

func makePostRequest() *server_proto.UserDetailsResponse {

	request := &server_proto.CreateUserRequest{
		FirstName:   "Krishna",
		LastName:    "Subhedarpage",
		Email:       "krishnasagar.subhedarpage@ahs.com",
		Designation: "SE3",
	}
	req, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("Unable to marshal request : %v", err)
	}

	resp, err := http.Post("http://0.0.0.0:9000/assignment/user", "application/x-binary", bytes.NewReader(req))
	if err != nil {
		log.Fatalf("Unable to read from the server : %v", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Unable to read bytes from request : %v", err)
	}

	respObj := &server_proto.UserDetailsResponse{}
	proto.Unmarshal(respBytes, respObj)
	return respObj

}

func makePatchRequest(user_id string) *server_proto.UserDetailsResponse {

	request := &server_proto.UpdateUserRequest{
		Email:  "krishnasagar.subhedarpage@frontdoor.com",
		UserId: user_id,
	}
	req, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("Unable to marshal request : %v", err)
	}

	act_req, _ := http.NewRequest("PATCH", "http://0.0.0.0:9000/assignment/user", bytes.NewReader(req))
	act_req.Header.Set("Content-Type", "application/x-binary")

	client := &http.Client{}
	resp, err := client.Do(act_req)
	if err != nil {
		log.Fatalf("Unable to read from the server : %v", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Fatalf("Unable to read bytes from request : %v", err)
	}

	respObj := &server_proto.UserDetailsResponse{}
	proto.Unmarshal(respBytes, respObj)
	return respObj

}

func main() {
	resp := makePostRequest()
	fmt.Printf("Response from API is : %v\n", resp)
	patch_resp := makePatchRequest(resp.UserId)
	fmt.Printf("Response from API is : %v\n", patch_resp)
}
