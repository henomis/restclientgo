package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/henomis/restclientgo"
)

type updatePostRequest struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Body   string `json:"body,omitempty"`
	UserID int    `json:"userId,omitempty"`
}

func (r *updatePostRequest) Path() (string, error) {
	return "/posts/" + fmt.Sprintf("%d", r.ID), nil
}

func (r *updatePostRequest) Encode() (string, error) {
	jsonBytes, err := json.Marshal(r)
	return string(jsonBytes), err
}

func (r *updatePostRequest) ContentType() string {
	return "application/json; charset=UTF-8"
}

type UpdatePostResponse struct {
	HTTPStatusCode int    `json:"-"`
	ID             int    `json:"id,omitempty"`
	UserID         int    `json:"userId,omitempty"`
	Title          string `json:"title,omitempty"`
	Body           string `json:"body,omitempty"`
}

func (r *UpdatePostResponse) Decode(body io.ReadCloser) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(r)
}

func (r *UpdatePostResponse) SetBody(body io.ReadCloser) {
	defer body.Close()
}

func (r *UpdatePostResponse) AcceptContentType() string {
	return "application/json"
}

func (r *UpdatePostResponse) SetStatusCode(code int) {
	r.HTTPStatusCode = code
}

func main() {

	restClient := restclientgo.New("https://jsonplaceholder.typicode.com")

	restClient.SetRequestModifier(func(req *http.Request) *http.Request {
		req.Header.Set("Accept", "application/json")
		return req
	})

	var updatePostResponse UpdatePostResponse

	err := restClient.Patch(
		&updatePostRequest{
			Title: "foo",
		},
		&updatePostResponse,
	)
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := json.MarshalIndent(updatePostResponse, "", "  ")

	fmt.Println("Status Code:", updatePostResponse.HTTPStatusCode)
	fmt.Println(string(bytes))

}
