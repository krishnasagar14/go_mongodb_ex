package main

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/golang/protobuf/proto"
	"go_mongodb_ex/proto"
)

func MakePostRequest() *server_proto.UserDetailsResponse {

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
	resp_bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Unable to read bytes from request : %v", err)
	}

	resp_obj := &server_proto.UserDetailsResponse{}
	proto.Unmarshal(resp_bytes, resp_obj)
	return resp_obj

}

func MakePatchRequest(user_id string) *server_proto.UserDetailsResponse {

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
	resp_bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Fatalf("Unable to read bytes from request : %v", err)
	}

	resp_obj := &server_proto.UserDetailsResponse{}
	proto.Unmarshal(resp_bytes, resp_obj)
	return resp_obj

}

func MakeGetRequest(user_id string) *server_proto.UserDetailsResponse {
	request := &server_proto.GetUserRequest{
		UserId: user_id,
	}
	p_body, _ := proto.Marshal(request)
	proto_body := b64.URLEncoding.EncodeToString(p_body)

	api_url := "http://0.0.0.0:9000/assignment/user"
	u, _ := url.Parse(api_url)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("proto_body", proto_body)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatalf("Unable to read from the server : %v", err)
	}
	resp_bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Unable to read bytes from request : %v", err)
	}

	resp_obj := &server_proto.UserDetailsResponse{}
	proto.Unmarshal(resp_bytes, resp_obj)
	return resp_obj
}

func main() {
	resp := MakePostRequest()
	fmt.Printf("Response from POST API is : %v\n", resp)
	patch_resp := MakePatchRequest(resp.UserId)
	fmt.Printf("Response from PATCH API is : %v\n", patch_resp)
	get_resp := MakeGetRequest(resp.UserId)
	fmt.Printf("Response from GET API is : %v\n", get_resp)
}
