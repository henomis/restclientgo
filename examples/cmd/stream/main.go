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
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

func (r *createPostRequest) Path() (string, error) {
	return "/generate", nil
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
	Model          string `json:"model"`
	CreatedAt      string `json:"created_at"`
	Response       string `json:"response"`
	Done           bool   `json:"done"`
}

func (r *CreatePostResponse) Decode(body io.Reader) error {

	return json.NewDecoder(body).Decode(r)
}

func (r *CreatePostResponse) SetBody(body io.Reader) error {
	return nil
}

func (r *CreatePostResponse) AcceptContentType() string {
	return "application/x-ndjson"
}

func (r *CreatePostResponse) SetStatusCode(code int) error {
	r.HTTPStatusCode = code
	return nil
}

func (r *CreatePostResponse) SetHeaders(headers restclientgo.Headers) error { return nil }

func main() {

	var response string
	restClient := restclientgo.New("http://localhost:11434/api")

	restClient.SetStreamCallback(
		func(data []byte) error {
			var createPostResponse CreatePostResponse

			err := json.Unmarshal(data, &createPostResponse)
			if err != nil {
				return err
			}

			response += createPostResponse.Response
			fmt.Printf(createPostResponse.Response)

			return nil
		},
	)

	restClient.SetRequestModifier(func(req *http.Request) *http.Request {
		req.Header.Set("Accept", "application/json")
		return req
	})

	var createPostResponse CreatePostResponse

	err := restClient.Post(
		context.Background(),
		&createPostRequest{
			Model:  "llama2",
			Prompt: "Why is the sky blue?",
		},
		&createPostResponse,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status Code:", createPostResponse.HTTPStatusCode)
}
