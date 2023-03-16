package main

import (
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

func (r *createPostRequest) Encode() (string, error) {
	jsonBytes, err := json.Marshal(r)
	return string(jsonBytes), err
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

func (r *CreatePostResponse) Decode(body io.ReadCloser) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(r)
}

func (r *CreatePostResponse) SetBody(body io.ReadCloser) {
	defer body.Close()
}

func (r *CreatePostResponse) AcceptContentType() string {
	return "application/json"
}

func (r *CreatePostResponse) SetStatusCode(code int) {
	r.HTTPStatusCode = code
}

func main() {

	restClient := restclientgo.New("https://jsonplaceholder.typicode.com")

	restClient.SetRequestModifier(func(req *http.Request) {
		req.Header.Set("Accept", "application/json")
	})

	var createPostResponse CreatePostResponse

	err := restClient.Post(
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
