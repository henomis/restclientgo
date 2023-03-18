package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/henomis/restclientgo"
)

type createPostRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func (r *createPostRequest) Path() (string, error) {
	return "/posts", nil
}

func (r *createPostRequest) Encode() (io.Reader, error) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBytes), nil
}

func (r *createPostRequest) ContentType() string {
	return "application/json"
}

type CreatePostResponse struct {
	HTTPStatusCode int    `json:"-"`
	ID             int    `json:"id"`
	UserID         int    `json:"userId"`
	Title          string `json:"title"`
	Body           string `json:"body"`
}

func (r *CreatePostResponse) Decode(body io.Reader) error {

	return json.NewDecoder(body).Decode(r)
}

func (r *CreatePostResponse) SetBody(body io.Reader) error {
	return nil
}

func (r *CreatePostResponse) AcceptContentType() string {
	return "application/json"
}

func (r *CreatePostResponse) SetStatusCode(code int) error {
	r.HTTPStatusCode = code
	return nil
}

func main() {

	restClient := restclientgo.New("https://jsonplaceholder.typicode.com")

	restClient.SetRequestModifier(func(req *http.Request) *http.Request {
		req.Header.Set("Accept", "application/json")
		return req
	})

	var createPostResponse CreatePostResponse

	err := restClient.Post(
		context.Background(),
		&createPostRequest{
			Title:  "foo",
			Body:   "bar",
			UserID: 1,
		},
		&createPostResponse,
	)
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := json.MarshalIndent(createPostResponse, "", "  ")

	fmt.Println("Status Code:", createPostResponse.HTTPStatusCode)
	fmt.Println(string(bytes))

}
