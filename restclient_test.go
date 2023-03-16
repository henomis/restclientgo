package restclientgo

import (
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

func (r *todoRequest) Path() (string, error)   { return "/todos/" + r.ID, nil }
func (r *todoRequest) Encode() (string, error) { return "", nil }
func (r *todoRequest) ContentType() string     { return "" }
func (r *TodoResponse) Decode(body io.ReadCloser) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(r)
}
func (r *TodoResponse) SetBody(body io.ReadCloser) { defer body.Close() }
func (r *TodoResponse) AcceptContentType() string  { return "application/json" }
func (r *TodoResponse) SetStatusCode(code int)     { r.HTTPStatusCode = code }

//---------------------------------------------

// DELETE

type deletePostRequest struct {
	ID int `json:"-"`
}
type DeletePostResponse struct {
	HTTPStatusCode int `json:"-"`
}

func (r *deletePostRequest) Path() (string, error) { return "/posts/" + fmt.Sprintf("%d", r.ID), nil }

func (r *deletePostRequest) Encode() (string, error) { return "", nil }
func (r *deletePostRequest) ContentType() string     { return "" }
func (r *DeletePostResponse) Decode(body io.ReadCloser) error {
	defer body.Close()
	return nil
}
func (r *DeletePostResponse) SetBody(body io.ReadCloser) { defer body.Close() }
func (r *DeletePostResponse) AcceptContentType() string  { return "" }
func (r *DeletePostResponse) SetStatusCode(code int)     { r.HTTPStatusCode = code }

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
func (r *updatePostRequest) Encode() (string, error) {
	jsonBytes, err := json.Marshal(r)
	return string(jsonBytes), err
}
func (r *updatePostRequest) ContentType() string { return "application/json; charset=UTF-8" }
func (r *UpdatePostResponse) Decode(body io.ReadCloser) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(r)
}
func (r *UpdatePostResponse) SetBody(body io.ReadCloser) { defer body.Close() }
func (r *UpdatePostResponse) AcceptContentType() string  { return "application/json" }
func (r *UpdatePostResponse) SetStatusCode(code int)     { r.HTTPStatusCode = code }

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
func (r *createPostRequest) Encode() (string, error) {
	jsonBytes, err := json.Marshal(r)
	return string(jsonBytes), err
}
func (r *createPostRequest) ContentType() string { return "application/json" }
func (r *CreatePostResponse) Decode(body io.ReadCloser) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(r)
}
func (r *CreatePostResponse) SetBody(body io.ReadCloser) { defer body.Close() }
func (r *CreatePostResponse) AcceptContentType() string  { return "application/json" }
func (r *CreatePostResponse) SetStatusCode(code int)     { r.HTTPStatusCode = code }

// ---------------------------------------------

func TestRestClient_Get(t *testing.T) {
	type fields struct {
		httpClient      *http.Client
		endpoint        string
		requestModifier func(*http.Request)
	}
	type args struct {
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
				requestModifier: func(req *http.Request) {
					req.Header.Set("Accept", "application/json")
				},
			},
			args: args{
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
				requestModifier: func(req *http.Request) {
					req.Header.Set("Accept", "application/json")
				},
			},
			args: args{
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
			if err := r.Get(tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RestClient.Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			response := tt.args.response.(*TodoResponse)
			wantResponse := tt.wantResponse.(*TodoResponse)

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
		requestModifier func(*http.Request)
	}
	type args struct {
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
				requestModifier: func(req *http.Request) {
					req.Header.Set("Accept", "application/json")
				},
			},
			args: args{
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
			if err := r.Delete(tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RestClient.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRestClient_Patch(t *testing.T) {
	type fields struct {
		httpClient      *http.Client
		endpoint        string
		requestModifier func(*http.Request)
	}
	type args struct {
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
				requestModifier: func(req *http.Request) {
					req.Header.Set("Accept", "application/json")
				},
			},
			args: args{
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
			if err := r.Patch(tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RestClient.Patch() error = %v, wantErr %v", err, tt.wantErr)
			}

			response := tt.args.response.(*UpdatePostResponse)
			wantResponse := tt.wantResponse.(*UpdatePostResponse)

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
		requestModifier func(*http.Request)
	}
	type args struct {
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
				requestModifier: func(req *http.Request) {
					req.Header.Set("Accept", "application/json")
				},
			},
			args: args{
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
			if err := r.Post(tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RestClient.Post() error = %v, wantErr %v", err, tt.wantErr)
			}

			response := tt.args.response.(*CreatePostResponse)
			wantResponse := tt.wantResponse.(*CreatePostResponse)

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
		requestModifier func(*http.Request)
	}
	type args struct {
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
				requestModifier: func(req *http.Request) {
					req.Header.Set("Accept", "application/json")
				},
			},
			args: args{
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
			if err := r.Put(tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("RestClient.Put() error = %v, wantErr %v", err, tt.wantErr)
			}

			response := tt.args.response.(*UpdatePostResponse)
			wantResponse := tt.wantResponse.(*UpdatePostResponse)

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
