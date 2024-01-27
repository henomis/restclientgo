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

type updatePostRequest struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func (r *updatePostRequest) Path() (string, error) {
	return "/posts/" + fmt.Sprintf("%d", r.ID), nil
}

func (r *updatePostRequest) Encode() (io.Reader, error) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBytes), nil
}

func (r *updatePostRequest) ContentType() string {
	return "application/json; charset=UTF-8"
}

type UpdatePostResponse struct {
	HTTPStatusCode int    `json:"-"`
	ID             int    `json:"id"`
	UserID         int    `json:"userId"`
	Title          string `json:"title"`
	Body           string `json:"body"`
}

func (r *UpdatePostResponse) Decode(body io.Reader) error {

	return json.NewDecoder(body).Decode(r)
}

func (r *UpdatePostResponse) SetBody(body io.Reader) error {
	return nil
}

func (r *UpdatePostResponse) AcceptContentType() string {
	return "application/json"
}

func (r *UpdatePostResponse) SetStatusCode(code int) error {
	r.HTTPStatusCode = code
	return nil
}

func (r *UpdatePostResponse) SetHeaders(headers restclientgo.Headers) error { return nil }

func main() {

	restClient := restclientgo.New("https://jsonplaceholder.typicode.com")

	restClient.SetRequestModifier(func(req *http.Request) *http.Request {
		req.Header.Set("Accept", "application/json")
		return req
	})

	var updatePostResponse UpdatePostResponse

	err := restClient.Put(
		context.Background(),
		&updatePostRequest{
			ID:     1,
			Title:  "foo",
			Body:   "bar",
			UserID: 1,
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
