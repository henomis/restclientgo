package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/henomis/restclientgo"
)

type deletePostRequest struct {
	ID int `json:"-"`
}

func (r *deletePostRequest) Path() (string, error) {
	return "/posts/" + fmt.Sprintf("%d", r.ID), nil
}

func (r *deletePostRequest) Encode() (string, error) {
	return "", nil
}

func (r *deletePostRequest) ContentType() string {
	return ""
}

type DeletePostResponse struct {
	HTTPStatusCode int `json:"-"`
}

func (r *DeletePostResponse) Decode(body io.ReadCloser) error {
	defer body.Close()
	return nil
}

func (r *DeletePostResponse) SetBody(body io.ReadCloser) {
	defer body.Close()
}

func (r *DeletePostResponse) AcceptContentType() string {
	return ""
}

func (r *DeletePostResponse) SetStatusCode(code int) {
	r.HTTPStatusCode = code
}

func main() {

	restClient := restclientgo.New("https://jsonplaceholder.typicode.com")

	restClient.SetRequestModifier(func(req *http.Request) {
		req.Header.Set("Accept", "application/json")
	})

	var deletePostResponse DeletePostResponse

	err := restClient.Delete(
		&deletePostRequest{ID: 1},
		&deletePostResponse,
	)
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := json.MarshalIndent(deletePostResponse, "", "  ")

	fmt.Println("Status Code:", deletePostResponse.HTTPStatusCode)
	fmt.Println(string(bytes))

}
