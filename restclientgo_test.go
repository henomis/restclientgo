package restclientgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// GET /todos/1
type todoRequest struct{ ID string }

type TodoResponse struct {
	HTTPStatusCode int    `json:"-"`
	ID             int    `json:"id"`
	UserID         int    `json:"userId"`
	Title          string `json:"title"`
	Completed      bool   `json:"completed"`
}

func (r *todoRequest) Path() (string, error)      { return "/todos/" + r.ID, nil }
func (r *todoRequest) Encode() (io.Reader, error) { return nil, nil }
func (r *todoRequest) ContentType() string        { return "" }
func (r *TodoResponse) Decode(body io.Reader) error {
	return json.NewDecoder(body).Decode(r)
}
func (r *TodoResponse) SetBody(body io.Reader) error {
	_ = body
	return nil
}
func (r *TodoResponse) AcceptContentType() string { return "application/json" }
func (r *TodoResponse) SetStatusCode(code int) error {
	r.HTTPStatusCode = code
	return nil
}
func (r *TodoResponse) SetHeaders(headers Headers) error {
	_ = headers
	return nil
}

//---------------------------------------------

// DELETE

type deletePostRequest struct {
	ID int `json:"-"`
}
type DeletePostResponse struct {
	HTTPStatusCode int `json:"-"`
}

func (r *deletePostRequest) Path() (string, error) { return "/posts/" + fmt.Sprintf("%d", r.ID), nil }

func (r *deletePostRequest) Encode() (io.Reader, error) { return nil, nil }
func (r *deletePostRequest) ContentType() string        { return "" }
func (r *DeletePostResponse) Decode(body io.Reader) error {
	_ = body
	return nil
}
func (r *DeletePostResponse) SetBody(body io.Reader) error {
	_ = body
	return nil
}
func (r *DeletePostResponse) AcceptContentType() string { return "" }
func (r *DeletePostResponse) SetStatusCode(code int) error {
	r.HTTPStatusCode = code
	return nil
}
func (r *DeletePostResponse) SetHeaders(headers Headers) error {
	_ = headers
	return nil
}

// ---------------------------------------------

// PATCH
type updatePostRequest struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title,omitempty"`
	Body   string `json:"body,omitempty"`
	UserID int    `json:"userId,omitempty"`
}
type UpdatePostResponse struct {
	HTTPStatusCode int    `json:"-"`
	ID             int    `json:"id,omitempty"`
	UserID         int    `json:"userId,omitempty"`
	Title          string `json:"title,omitempty"`
	Body           string `json:"body,omitempty"`
}

func (r *updatePostRequest) Path() (string, error) { return "/posts/" + fmt.Sprintf("%d", r.ID), nil }
func (r *updatePostRequest) Encode() (io.Reader, error) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBytes), nil
}
func (r *updatePostRequest) ContentType() string { return "application/json; charset=UTF-8" }
func (r *UpdatePostResponse) Decode(body io.Reader) error {
	return json.NewDecoder(body).Decode(r)
}
func (r *UpdatePostResponse) SetBody(body io.Reader) error {
	_ = body
	return nil
}
func (r *UpdatePostResponse) AcceptContentType() string { return "application/json" }
func (r *UpdatePostResponse) SetStatusCode(code int) error {
	r.HTTPStatusCode = code
	return nil
}
func (r *UpdatePostResponse) SetHeaders(headers Headers) error {
	_ = headers
	return nil
}

// ---------------------------------------------

// POST

type createPostRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

type CreatePostResponse struct {
	HTTPStatusCode int    `json:"-"`
	ID             int    `json:"id"`
	UserID         int    `json:"userId"`
	Title          string `json:"title"`
	Body           string `json:"body"`
}

func (r *createPostRequest) Path() (string, error) { return "/posts", nil }
func (r *createPostRequest) Encode() (io.Reader, error) {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBytes), nil
}
func (r *createPostRequest) ContentType() string { return "application/json" }
func (r *CreatePostResponse) Decode(body io.Reader) error {
	return json.NewDecoder(body).Decode(r)
}
func (r *CreatePostResponse) SetBody(body io.Reader) error {
	_ = body
	return nil
}
func (r *CreatePostResponse) AcceptContentType() string { return "application/json" }
func (r *CreatePostResponse) SetStatusCode(code int) error {
	r.HTTPStatusCode = code
	return nil
}
func (r *CreatePostResponse) SetHeaders(headers Headers) error {
	_ = headers
	return nil
}

// ---------------------------------------------

func TestRestClient_Get(t *testing.T) {
	type fields struct {
		httpClient      *http.Client
		endpoint        string
		requestModifier func(*http.Request) *http.Request
	}
	type args struct {
		ctx      context.Context
		request  Request
		response Response
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantResponse Response
	}{
		{
			name: "test get",
			fields: fields{
				httpClient: http.DefaultClient,
				endpoint:   "https://jsonplaceholder.typicode.com",
				requestModifier: func(req *http.Request) *http.Request {
					req.Header.Set("Accept", "application/json")
					return req
				},
			},
			args: args{
				ctx: context.Background(),
				request: &todoRequest{
					ID: "1",
				},
				response: &TodoResponse{},
			},
			wantErr: false,
			wantResponse: &TodoResponse{
				ID:        1,
				UserID:    1,
				Title:     "delectus aut autem",
				Completed: false,
			},
		},
		{
			name: "test get",
			fields: fields{
				httpClient: http.DefaultClient,
				endpoint:   "https://jsonplaceholder.typicode.com",
				requestModifier: func(req *http.Request) *http.Request {
					req.Header.Set("Accept", "application/json")
					return req
				},
			},
			args: args{
				ctx: context.Background(),
				request: &todoRequest{
					ID: "1",
				},
				response: &TodoResponse{},
			},
			wantErr: false,
			wantResponse: &TodoResponse{
				ID:        1,
				UserID:    1,
				Title:     "delectus aut autem",
				Completed: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestClient{
				httpClient:      tt.fields.httpClient,
				endpoint:        tt.fields.endpoint,
				requestModifier: tt.fields.requestModifier,
			}
			if err := r.Get(tt.args.ctx, tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RestClient.Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			response, ok := tt.args.response.(*TodoResponse)
			if !ok {
				t.Errorf("Cannot cast response to TodoResponse")
			}
			wantResponse, ok := tt.wantResponse.(*TodoResponse)
			if !ok {
				t.Errorf("Cannot cast response to TodoResponse")
			}

			jsonReponse, err := json.Marshal(*response)
			if err != nil {
				t.Errorf("RestClient.Get() error = %v", err)
			}

			jsonWantResponse, err := json.Marshal(*wantResponse)
			if err != nil {
				t.Errorf("RestClient.Get() error = %v", err)
			}

			if string(jsonReponse) != string(jsonWantResponse) {
				t.Errorf("RestClient.Get() = %s, want %s", jsonReponse, jsonWantResponse)
			}
		})
	}
}

func TestRestClient_Delete(t *testing.T) {
	type fields struct {
		httpClient      *http.Client
		endpoint        string
		requestModifier func(*http.Request) *http.Request
	}
	type args struct {
		ctx      context.Context
		request  Request
		response Response
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test delete",
			fields: fields{
				httpClient: http.DefaultClient,
				endpoint:   "https://jsonplaceholder.typicode.com",
				requestModifier: func(req *http.Request) *http.Request {
					req.Header.Set("Accept", "application/json")
					return req
				},
			},
			args: args{
				ctx: context.Background(),
				request: &deletePostRequest{
					ID: 1,
				},
				response: &DeletePostResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestClient{
				httpClient:      tt.fields.httpClient,
				endpoint:        tt.fields.endpoint,
				requestModifier: tt.fields.requestModifier,
			}
			if err := r.Delete(tt.args.ctx, tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RestClient.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestClient_Patch(t *testing.T) {
	type fields struct {
		httpClient      *http.Client
		endpoint        string
		requestModifier func(*http.Request) *http.Request
	}
	type args struct {
		ctx      context.Context
		request  Request
		response Response
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantResponse Response
	}{
		{
			name: "test patch",
			fields: fields{
				httpClient: http.DefaultClient,
				endpoint:   "https://jsonplaceholder.typicode.com",
				requestModifier: func(req *http.Request) *http.Request {
					req.Header.Set("Accept", "application/json")
					return req
				},
			},
			args: args{
				ctx: context.Background(),
				request: &updatePostRequest{
					Title: "foo",
				},
				response: &UpdatePostResponse{},
			},
			wantErr: false,
			wantResponse: &UpdatePostResponse{
				Title: "foo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestClient{
				httpClient:      tt.fields.httpClient,
				endpoint:        tt.fields.endpoint,
				requestModifier: tt.fields.requestModifier,
			}
			if err := r.Patch(tt.args.ctx, tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RestClient.Patch() error = %v, wantErr %v", err, tt.wantErr)
			}

			response, ok := tt.args.response.(*UpdatePostResponse)
			if !ok {
				t.Errorf("Cannot cast response to UpdatePostResponse")
			}
			wantResponse, ok := tt.wantResponse.(*UpdatePostResponse)
			if !ok {
				t.Errorf("Cannot cast wantResponse to UpdatePostResponse")
			}

			jsonReponse, err := json.Marshal(*response)
			if err != nil {
				t.Errorf("RestClient.Get() error = %v", err)
			}

			jsonWantResponse, err := json.Marshal(*wantResponse)
			if err != nil {
				t.Errorf("RestClient.Get() error = %v", err)
			}

			if string(jsonReponse) != string(jsonWantResponse) {
				t.Errorf("RestClient.Get() = %s, want %s", jsonReponse, jsonWantResponse)
			}
		})
	}
}

func TestRestClient_Post(t *testing.T) {
	type fields struct {
		httpClient      *http.Client
		endpoint        string
		requestModifier func(*http.Request) *http.Request
	}
	type args struct {
		ctx      context.Context
		request  Request
		response Response
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantResponse Response
	}{
		{
			name: "test post",
			fields: fields{
				httpClient: http.DefaultClient,
				endpoint:   "https://jsonplaceholder.typicode.com",
				requestModifier: func(req *http.Request) *http.Request {
					req.Header.Set("Accept", "application/json")
					return req
				},
			},
			args: args{
				ctx: context.Background(),
				request: &createPostRequest{
					Title:  "foo",
					Body:   "bar",
					UserID: 1,
				},
				response: &CreatePostResponse{},
			},
			wantErr: false,
			wantResponse: &CreatePostResponse{
				ID:     101,
				Title:  "foo",
				Body:   "bar",
				UserID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestClient{
				httpClient:      tt.fields.httpClient,
				endpoint:        tt.fields.endpoint,
				requestModifier: tt.fields.requestModifier,
			}
			if err := r.Post(tt.args.ctx, tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RestClient.Post() error = %v, wantErr %v", err, tt.wantErr)
			}

			response, ok := tt.args.response.(*CreatePostResponse)
			if !ok {
				t.Errorf("Cannot cast response to CreatePostResponse")
			}
			wantResponse, ok := tt.wantResponse.(*CreatePostResponse)
			if !ok {
				t.Errorf("Cannot cast wantResponse to CreatePostResponse")
			}

			jsonReponse, err := json.Marshal(*response)
			if err != nil {
				t.Errorf("RestClient.Get() error = %v", err)
			}

			jsonWantResponse, err := json.Marshal(*wantResponse)
			if err != nil {
				t.Errorf("RestClient.Get() error = %v", err)
			}

			if string(jsonReponse) != string(jsonWantResponse) {
				t.Errorf("RestClient.Get() = %s, want %s", jsonReponse, jsonWantResponse)
			}
		})
	}
}

func TestRestClient_Put(t *testing.T) {
	type fields struct {
		httpClient      *http.Client
		endpoint        string
		requestModifier func(*http.Request) *http.Request
	}
	type args struct {
		ctx      context.Context
		request  Request
		response Response
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      bool
		wantResponse Response
	}{
		{
			name: "test put",
			fields: fields{
				httpClient: http.DefaultClient,
				endpoint:   "https://jsonplaceholder.typicode.com",
				requestModifier: func(req *http.Request) *http.Request {
					req.Header.Set("Accept", "application/json")
					return req
				},
			},
			args: args{
				ctx: context.Background(),
				request: &updatePostRequest{
					ID:     1,
					Title:  "foo",
					Body:   "bar",
					UserID: 1,
				},
				response: &UpdatePostResponse{},
			},
			wantErr: false,
			wantResponse: &UpdatePostResponse{
				ID:     1,
				Title:  "foo",
				Body:   "bar",
				UserID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestClient{
				httpClient:      tt.fields.httpClient,
				endpoint:        tt.fields.endpoint,
				requestModifier: tt.fields.requestModifier,
			}
			if err := r.Put(tt.args.ctx, tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RestClient.Put() error = %v, wantErr %v", err, tt.wantErr)
			}

			response, ok := tt.args.response.(*UpdatePostResponse)
			if !ok {
				t.Errorf("Cannot cast response to UpdatePostResponse")
			}
			wantResponse, ok := tt.wantResponse.(*UpdatePostResponse)
			if !ok {
				t.Errorf("Cannot cast wantResponse to UpdatePostResponse")
			}

			jsonReponse, err := json.Marshal(*response)
			if err != nil {
				t.Errorf("RestClient.Get() error = %v", err)
			}

			jsonWantResponse, err := json.Marshal(*wantResponse)
			if err != nil {
				t.Errorf("RestClient.Get() error = %v", err)
			}

			if string(jsonReponse) != string(jsonWantResponse) {
				t.Errorf("RestClient.Get() = %s, want %s", jsonReponse, jsonWantResponse)
			}
		})
	}
}

func TestRestClient_SetEndpoint(t *testing.T) {
	type fields struct {
		httpClient      *http.Client
		endpoint        string
		requestModifier func(*http.Request) *http.Request
	}
	type args struct {
		endpoint string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test set endpoint",
			fields: fields{
				httpClient: http.DefaultClient,
				endpoint:   "https://example.com",
				requestModifier: func(req *http.Request) *http.Request {
					req.Header.Set("Accept", "application/json")
					return req
				},
			},
			args: args{
				endpoint: "https://another.com",
			},
			want: "https://another.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestClient{
				httpClient:      tt.fields.httpClient,
				endpoint:        tt.fields.endpoint,
				requestModifier: tt.fields.requestModifier,
			}
			r.SetEndpoint(tt.args.endpoint)

			if tt.want != r.endpoint {
				t.Errorf("RestClient.Get() = %s, want %s", tt.want, r.endpoint)
			}
		})
	}
}

func TestRestClient_Endpoint(t *testing.T) {
	type fields struct {
		httpClient      *http.Client
		endpoint        string
		requestModifier func(*http.Request) *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test get endpoint",
			fields: fields{
				httpClient: http.DefaultClient,
				endpoint:   "https://example.com",
				requestModifier: func(req *http.Request) *http.Request {
					req.Header.Set("Accept", "application/json")
					return req
				},
			},
			want: "https://example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RestClient{
				httpClient:      tt.fields.httpClient,
				endpoint:        tt.fields.endpoint,
				requestModifier: tt.fields.requestModifier,
			}
			if got := r.Endpoint(); got != tt.want {
				t.Errorf("RestClient.Endpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
