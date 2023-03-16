# RESTclientGo

[![Build Status](https://github.com/henomis/restclientgo/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/henomis/restclientgo/actions/workflows/test.yml?query=branch%3Amain) [![GoDoc](https://godoc.org/github.com/henomis/restclientgo?status.svg)](https://godoc.org/github.com/henomis/restclientgo) [![Go Report Card](https://goreportcard.com/badge/github.com/henomis/restclientgo)](https://goreportcard.com/report/github.com/henomis/restclientgo) [![GitHub release](https://img.shields.io/github/release/henomis/restclientgo.svg)](https://github.com/henomis/restclientgo/releases)


REST client for Go focusing on Request and Response data modeling.

## Rest methods
restclientgo offers the following methods

* GET
* POST
* PUT
* DELETE
* PATCH

## Modeling

### Request
Define your request model and attach restclientgo methods to satisfy the Request interface.

```go
type MyRequest struct {
    // Add your request fields here
    // ID string `json:"id"`
    // ...
}

func (r *MyRequest) Path() (string, error) {
    // Return the path of the request including the query string if any.
}

func (r *MyRequest) Encode() (string, error) {
    // Return the request body as string
}

func (r *MyRequest) ContentType() string {
    // Return the content type of the request
}
```

### Response
Define your response model and attach restclientgo methods to satisfy the Response interface.

```go
type MyResponse struct {
    // Add your response fields here
    // ID string `json:"id"`
    // ...
}

func (r *MyResponse) Decode(body io.Reader) error {
    // Decode the response body into the response model
}

func (r *MyResponse) SetBody(body io.ReadCloser) {
    // Set the response body if needed
}

func (r *MyResponse) StatusCode() int {
    // Handler to set the HTTP status code of the response
}

func (r *MyResponse) AcceptContentType() string {
    // Return the accepted content type of the response
}
```

## Usage
Please referr to the [examples](examples/cmd/) folder for usage examples.