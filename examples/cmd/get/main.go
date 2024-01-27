package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/henomis/restclientgo"
)

type todoRequest struct {
	ID string
}

func (r *todoRequest) Path() (string, error) {
	return "/todos/" + r.ID, nil
}

func (r *todoRequest) Encode() (io.Reader, error) {
	return nil, nil
}

func (r *todoRequest) ContentType() string {
	return ""
}

type TodoResponse struct {
	HTTPStatusCode int    `json:"-"`
	ID             int    `json:"id"`
	UserID         int    `json:"userId"`
	Title          string `json:"title"`
	Comleted       bool   `json:"completed"`
}

func (r *TodoResponse) Decode(body io.Reader) error {

	return json.NewDecoder(body).Decode(r)
}

func (r *TodoResponse) SetBody(body io.Reader) error {
	return nil
}

func (r *TodoResponse) AcceptContentType() string {
	return "application/json"
}

func (r *TodoResponse) SetStatusCode(code int) error {
	r.HTTPStatusCode = code
	return nil
}

func (r *TodoResponse) SetHeaders(headers restclientgo.Headers) error { return nil }

func main() {

	restClient := restclientgo.New("https://jsonplaceholder.typicode.com")

	restClient.SetRequestModifier(func(req *http.Request) *http.Request {
		req.Header.Set("Accept", "application/json")
		return req
	})

	var todoResponse TodoResponse

	err := restClient.Get(
		context.Background(),
		&todoRequest{ID: "1"},
		&todoResponse,
	)
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := json.MarshalIndent(todoResponse, "", "  ")

	fmt.Println("Status Code:", todoResponse.HTTPStatusCode)
	fmt.Println(string(bytes))

}
