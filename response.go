package cumi

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
	"time"
)

// Response represents an HTTP response
type Response struct {
	Request    *Request
	Response   *http.Response
	body       []byte
	size       int64
	receivedAt time.Time
	duration   time.Duration
	state      ResultState
	Err        error

	// Embedded from http.Response for direct access
	Status     string
	StatusCode int
	Proto      string
	ProtoMajor int
	ProtoMinor int
	Header     http.Header
}

// Body returns the response body as bytes
func (r *Response) Body() []byte {
	return r.body
}

// String returns the response body as a string
func (r *Response) String() string {
	return string(r.body)
}

// JSON unmarshals the response body into the provided interface using JSON
func (r *Response) JSON(v interface{}) error {
	if len(r.body) == 0 {
		return nil
	}
	return json.Unmarshal(r.body, v)
}

// XML unmarshals the response body into the provided interface using XML
func (r *Response) XML(v interface{}) error {
	if len(r.body) == 0 {
		return nil
	}
	return xml.Unmarshal(r.body, v)
}

// IsSuccess returns true if the response is successful (2xx status code)
func (r *Response) IsSuccess() bool {
	return r.state == SuccessState
}

// IsError returns true if the response is an error (4xx or 5xx status code)
func (r *Response) IsError() bool {
	return r.state == ErrorState
}

// Time returns the time when the response was received
func (r *Response) Time() time.Time {
	return r.receivedAt
}

// Duration returns the time taken for the request
func (r *Response) Duration() time.Duration {
	return r.duration
}

// Size returns the size of the response body in bytes
func (r *Response) Size() int64 {
	return r.size
}

// ResultState returns the state of the response
func (r *Response) ResultState() ResultState {
	return r.state
}

// Error returns the error if any occurred during the request
func (r *Response) Error() error {
	return r.Err
}

// ContentType returns the Content-Type header value
func (r *Response) ContentType() string {
	return r.Header.Get("Content-Type")
}

// IsJSON returns true if the response content type is JSON
func (r *Response) IsJSON() bool {
	contentType := r.ContentType()
	return strings.Contains(contentType, "application/json")
}

// IsXML returns true if the response content type is XML
func (r *Response) IsXML() bool {
	contentType := r.ContentType()
	return strings.Contains(contentType, "application/xml") ||
		strings.Contains(contentType, "text/xml")
}

// IsHTML returns true if the response content type is HTML
func (r *Response) IsHTML() bool {
	contentType := r.ContentType()
	return strings.Contains(contentType, "text/html")
}

// IsText returns true if the response content type is plain text
func (r *Response) IsText() bool {
	contentType := r.ContentType()
	return strings.Contains(contentType, "text/plain")
}

// Cookies returns the cookies set by the server
func (r *Response) Cookies() []*http.Cookie {
	if r.Response == nil {
		return nil
	}
	return r.Response.Cookies()
}

// Location returns the Location header value (useful for redirects)
func (r *Response) Location() string {
	return r.Header.Get("Location")
}
