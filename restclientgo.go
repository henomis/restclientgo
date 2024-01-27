package restclientgo

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const maxStreamBufferSize = 512 * 1024

type StreamCallback func([]byte) error

type RestClient struct {
	httpClient         *http.Client
	endpoint           string
	requestModifier    func(*http.Request) *http.Request
	forceDecodeOnError bool
}

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNoContentType  = Error("no content-type found in response")
	ErrRequestPath    = Error("invalid request path")
	ErrRequestEncode  = Error("invalid request encode")
	ErrHTTPRequest    = Error("invalid http request")
	ErrResponseDecode = Error("invalid response decode")
)

type httpMethod string

const (
	methodGet    httpMethod = "GET"
	methodPost   httpMethod = "POST"
	methodDelete httpMethod = "DELETE"
	methodPut    httpMethod = "PUT"
	methodPatch  httpMethod = "PATCH"
)

type Headers map[string][]string

type Request interface {
	// Path returns the path of the request including the query string if any.
	Path() (string, error)
	// Encode return the encoded representation of the request.
	Encode() (io.Reader, error)
	ContentType() string
}

type Response interface {
	// Decode decodes the response body into the given interface if the
	// response matches the AcceptContentType.
	Decode(body io.Reader) error
	// SetBody sets the response raw body if the response can't be decoded.
	SetBody(body io.Reader) error
	// AcceptContentType returns the content type that the response should be decoded to.
	AcceptContentType() string
	// SetStatusCode sets the HTTP response status code.
	SetStatusCode(code int) error
	// SetHeaders sets the HTTP response headers.
	SetHeaders(headers Headers) error
}

type Streamable interface {
	// StreamCallback get the stream callback if any.
	StreamCallback() StreamCallback
}

// New creates a new RestClient.
func New(endpoint string) *RestClient {
	return &RestClient{
		endpoint:   endpoint,
		httpClient: &http.Client{},
	}
}

// SetHTTPClient overrides the default http client.
func (r *RestClient) SetHTTPClient(client *http.Client) {
	r.httpClient = client
}

// SetRequestModifier adds a function that will modify each request
func (r *RestClient) SetRequestModifier(requestModifier func(*http.Request) *http.Request) {
	r.requestModifier = requestModifier
}

// WithRequestModifier adds a function that will modify each request
func (r *RestClient) WithRequestModifier(requestModifier func(*http.Request) *http.Request) *RestClient {
	r.requestModifier = requestModifier
	return r
}

// WithHTTPClient overrides the default http client.
func (r *RestClient) WithHTTPClient(client *http.Client) *RestClient {
	r.httpClient = client
	return r
}

// WithDecodeOnError forces the response to be decoded even if the status code is >= 400.
func (r *RestClient) WithDecodeOnError(decodeOnError bool) *RestClient {
	r.forceDecodeOnError = decodeOnError
	return r
}

func (r *RestClient) SetEndpoint(endpoint string) {
	r.endpoint = endpoint
}

func (r *RestClient) Endpoint() string {
	return r.endpoint
}

// Get performs a GET request.
func (r *RestClient) Get(ctx context.Context, request Request, response Response) error {
	return r.do(ctx, methodGet, request, response)
}

// Post performs a POST request.
func (r *RestClient) Post(ctx context.Context, request Request, response Response) error {
	return r.do(ctx, methodPost, request, response)
}

// Delete performs a DELETE request.
func (r *RestClient) Delete(ctx context.Context, request Request, response Response) error {
	return r.do(ctx, methodDelete, request, response)
}

// Put performs a PUT request.
func (r *RestClient) Put(ctx context.Context, request Request, response Response) error {
	return r.do(ctx, methodPut, request, response)
}

// Patch performs a PATCH request.
func (r *RestClient) Patch(ctx context.Context, request Request, response Response) error {
	return r.do(ctx, methodPatch, request, response)
}

//nolint:gocognit
func (r *RestClient) do(ctx context.Context, method httpMethod, request Request, response Response) error {
	requestPath, err := request.Path()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRequestPath, err)
	}

	requestEncodedBody, err := request.Encode()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRequestEncode, err)
	}

	requestURL := r.endpoint + requestPath
	httpRequest, err := http.NewRequest(string(method), requestURL, requestEncodedBody)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrHTTPRequest, err)
	}

	if request.ContentType() != "" {
		httpRequest.Header.Set("Content-Type", request.ContentType())
	}

	if r.requestModifier != nil {
		httpRequest = r.requestModifier(httpRequest)
	}

	httpRequest = httpRequest.WithContext(ctx)

	httpResponse, err := r.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrHTTPRequest, err)
	}
	defer httpResponse.Body.Close()

	var headers = make(Headers)
	for k, v := range httpResponse.Header {
		headers[k] = v
	}

	err = response.SetHeaders(headers)
	if err != nil {
		return err
	}

	err = response.SetStatusCode(httpResponse.StatusCode)
	if err != nil {
		return err
	}

	if httpResponse.StatusCode >= 400 && !r.forceDecodeOnError {
		err = response.SetBody(httpResponse.Body)
		if err != nil {
			return err
		}
		return nil
	}

	if response.AcceptContentType() == "" {
		err = response.SetBody(httpResponse.Body)
		if err != nil {
			return err
		}
		return nil
	}

	err = matchContentType(httpResponse, response)
	if err != nil {
		return err
	}

	if streamable, isStreamable := response.(Streamable); isStreamable && streamable.StreamCallback() != nil {
		err = stream(streamable.StreamCallback(), httpResponse.Body)
	} else {
		err = response.Decode(httpResponse.Body)
	}

	if err != nil {
		return fmt.Errorf("%w: %w", ErrResponseDecode, err)
	}

	return nil
}

func matchContentType(httpResponse *http.Response, response Response) error {
	contentTypeToMatch := response.AcceptContentType()
	contentType := httpResponse.Header.Get("Content-Type")

	if contentType == "" {
		return ErrNoContentType
	}

	for _, v := range strings.Split(contentType, ";") {
		if strings.TrimSpace(v) == contentTypeToMatch {
			return nil
		}
	}

	return ErrNoContentType
}

func stream(streamCallback StreamCallback, body io.Reader) error {
	scanner := bufio.NewScanner(body)

	scanBuf := make([]byte, 0, maxStreamBufferSize)
	scanner.Buffer(scanBuf, maxStreamBufferSize)

	for scanner.Scan() {
		err := streamCallback(scanner.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}
