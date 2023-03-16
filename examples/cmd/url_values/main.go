package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/henomis/restclientgo"
)

type commentsRequest struct {
	PostID string
}

func (r *commentsRequest) Path() (string, error) {

	urlValues := restclientgo.NewURLValues()
	urlValues.Add("postId", &r.PostID)

	return "/comments?" + urlValues.Encode(), nil
}

func (r *commentsRequest) Encode() (string, error) {
	return "", nil
}

func (r *commentsRequest) ContentType() string {
	return ""
}

type CommentsResponse struct {
	HTTPStatusCode int `json:"-"`
	Data           []struct {
		ID     int    `json:"id"`
		PostID int    `json:"postId"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Body   string `json:"body"`
	} `json:"data"`
}

func (r *CommentsResponse) Decode(body io.ReadCloser) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(&r.Data)
}

func (r *CommentsResponse) SetBody(body io.ReadCloser) {
	defer body.Close()
}

func (r *CommentsResponse) AcceptContentType() string {
	return "application/json"
}

func (r *CommentsResponse) SetStatusCode(code int) {
	r.HTTPStatusCode = code
}

func main() {

	restClient := restclientgo.New("https://jsonplaceholder.typicode.com")

	restClient.SetRequestModifier(func(req *http.Request) {
		req.Header.Set("Accept", "application/json")
	})

	var commentsResponse CommentsResponse

	err := restClient.Get(
		&commentsRequest{PostID: "1"},
		&commentsResponse,
	)
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := json.MarshalIndent(commentsResponse, "", "  ")

	fmt.Println("Status Code:", commentsResponse.HTTPStatusCode)
	fmt.Println(string(bytes))

}
