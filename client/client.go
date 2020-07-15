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

func MakePatchRequest(userId string) *server_proto.UserDetailsResponse {

	request := &server_proto.UpdateUserRequest{
		Email:  "krishnasagar.subhedarpage@frontdoor.com",
		UserId: userId,
	}
	req, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("Unable to marshal request : %v", err)
	}

	actReq, _ := http.NewRequest("PATCH", "http://0.0.0.0:9000/assignment/user", bytes.NewReader(req))
	actReq.Header.Set("Content-Type", "application/x-binary")

	client := &http.Client{}
	resp, err := client.Do(actReq)
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

func MakeGetRequest(userId string) *server_proto.UserDetailsResponse {
	request := &server_proto.GetUserRequest{
		UserId: userId,
	}
	pBody, _ := proto.Marshal(request)
	protoBody := b64.URLEncoding.EncodeToString(pBody)

	apiUrl := "http://0.0.0.0:9000/assignment/user"
	u, _ := url.Parse(apiUrl)
	q, _ := url.ParseQuery(u.RawQuery)
	q.Add("proto_body", protoBody)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
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

func main() {
	postResp := MakePostRequest()
	fmt.Printf("Response from POST API is : %v\n", postResp)
	patchResp := MakePatchRequest(postResp.UserId)
	fmt.Printf("Response from PATCH API is : %v\n", patchResp)
	getResp := MakeGetRequest(postResp.UserId)
	fmt.Printf("Response from GET API is : %v\n", getResp)
}
