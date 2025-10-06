package cumi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel/trace"
)

// Request represents an HTTP request
type Request struct {
	client      *Client
	method      string
	url         string
	ctx         context.Context
	headers     http.Header
	queryParams url.Values
	pathParams  map[string]string
	formData    url.Values
	body        interface{}
	bodyType    string
	cookies     []*http.Cookie
	userAgent   string
	basicAuth   struct {
		username string
		password string
	}
	bearerToken    string
	successResult  interface{}
	errorResult    interface{}
	downloadPath   string
	uploadCallback func(written int64, total int64)
	tracer         trace.Tracer
	spanName       string
}

// SetContext sets the context for the request
func (r *Request) SetContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// Context returns the request context
func (r *Request) Context() context.Context {
	if r.ctx == nil {
		return context.Background()
	}
	return r.ctx
}

// SetHeader sets a header for the request
func (r *Request) SetHeader(key, value string) *Request {
	r.headers.Set(key, value)
	return r
}

// SetHeaders sets multiple headers from a map
func (r *Request) SetHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.headers.Set(k, v)
	}
	return r
}

// SetUserAgent sets the User-Agent header for this specific request
func (r *Request) SetUserAgent(userAgent string) *Request {
	r.userAgent = userAgent
	return r
}

// SetHeaderVerbatim sets a header without canonicalizing the key
func (r *Request) SetHeaderVerbatim(key, value string) *Request {
	r.headers[key] = []string{value}
	return r
}

// SetQueryParam sets a query parameter for the request
func (r *Request) SetQueryParam(key, value string) *Request {
	r.queryParams.Set(key, value)
	return r
}

// SetQueryParams sets multiple query parameters from a map
func (r *Request) SetQueryParams(params map[string]string) *Request {
	for k, v := range params {
		r.queryParams.Set(k, v)
	}
	return r
}

// SetQueryParamsFromValues sets query parameters from url.Values
func (r *Request) SetQueryParamsFromValues(params url.Values) *Request {
	for k, values := range params {
		for _, v := range values {
			r.queryParams.Add(k, v)
		}
	}
	return r
}

// SetQueryString sets the query string directly
func (r *Request) SetQueryString(query string) *Request {
	values, err := url.ParseQuery(query)
	if err == nil {
		r.queryParams = values
	}
	return r
}

// SetPathParam sets a path parameter for URL replacement
func (r *Request) SetPathParam(key, value string) *Request {
	if r.pathParams == nil {
		r.pathParams = make(map[string]string)
	}
	r.pathParams[key] = value
	return r
}

// SetPathParams sets multiple path parameters from a map
func (r *Request) SetPathParams(params map[string]string) *Request {
	if r.pathParams == nil {
		r.pathParams = make(map[string]string)
	}
	for k, v := range params {
		r.pathParams[k] = v
	}
	return r
}

// SetFormData sets form data for the request
func (r *Request) SetFormData(data map[string]string) *Request {
	for k, v := range data {
		r.formData.Set(k, v)
	}
	return r
}

// SetFormDataFromValues sets form data from url.Values
func (r *Request) SetFormDataFromValues(data url.Values) *Request {
	for k, values := range data {
		for _, v := range values {
			r.formData.Add(k, v)
		}
	}
	return r
}

// SetBody sets the request body
func (r *Request) SetBody(body interface{}) *Request {
	r.body = body
	return r
}

// SetBodyBytes sets the request body from bytes
func (r *Request) SetBodyBytes(body []byte) *Request {
	r.body = body
	return r
}

// SetBodyString sets the request body from a string
func (r *Request) SetBodyString(body string) *Request {
	r.body = body
	return r
}

// SetBodyReader sets the request body from an io.Reader
func (r *Request) SetBodyReader(body io.Reader) *Request {
	r.body = body
	return r
}

// SetBodyJSON sets the request body as JSON
func (r *Request) SetBodyJSON(body interface{}) *Request {
	r.body = body
	r.bodyType = "json"
	return r
}

// SetBodyXML sets the request body as XML
func (r *Request) SetBodyXML(body interface{}) *Request {
	r.body = body
	r.bodyType = "xml"
	return r
}

// SetBasicAuth sets basic authentication
func (r *Request) SetBasicAuth(username, password string) *Request {
	r.basicAuth.username = username
	r.basicAuth.password = password
	return r
}

// SetBearerToken sets the bearer token for authentication
func (r *Request) SetBearerToken(token string) *Request {
	r.bearerToken = token
	return r
}

// SetAuthToken is an alias for SetBearerToken
func (r *Request) SetAuthToken(token string) *Request {
	return r.SetBearerToken(token)
}

// SetCookies sets cookies for the request
func (r *Request) SetCookies(cookies ...*http.Cookie) *Request {
	r.cookies = append(r.cookies, cookies...)
	return r
}

// SetCookie sets a single cookie for the request
func (r *Request) SetCookie(cookie *http.Cookie) *Request {
	r.cookies = append(r.cookies, cookie)
	return r
}

// SetSuccessResult sets the struct to unmarshal successful response into
func (r *Request) SetSuccessResult(result interface{}) *Request {
	r.successResult = result
	return r
}

// SetResult is an alias for SetSuccessResult
func (r *Request) SetResult(result interface{}) *Request {
	return r.SetSuccessResult(result)
}

// SetErrorResult sets the struct to unmarshal error response into
func (r *Request) SetErrorResult(result interface{}) *Request {
	r.errorResult = result
	return r
}

// SetError is an alias for SetErrorResult
func (r *Request) SetError(result interface{}) *Request {
	return r.SetErrorResult(result)
}

// SetTracer sets the tracer and span name for tracing HTTP request
func (r *Request) SetTracer(tracer trace.Tracer, spanName string) *Request {
	r.tracer = tracer
	r.spanName = spanName
	return r
}

// SetOutput sets the file path to save the response body
func (r *Request) SetOutput(filePath string) *Request {
	r.downloadPath = filePath
	return r
}

// SetUploadCallback sets a callback function for upload progress
func (r *Request) SetUploadCallback(callback func(written int64, total int64)) *Request {
	r.uploadCallback = callback
	return r
}

// Get executes a GET request
func (r *Request) Get(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodGet
	return r.Execute()
}

// Post executes a POST request
func (r *Request) Post(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodPost
	return r.Execute()
}

// Put executes a PUT request
func (r *Request) Put(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodPut
	return r.Execute()
}

// Patch executes a PATCH request
func (r *Request) Patch(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodPatch
	return r.Execute()
}

// Delete executes a DELETE request
func (r *Request) Delete(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodDelete
	return r.Execute()
}

// Head executes a HEAD request
func (r *Request) Head(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodHead
	return r.Execute()
}

// Options executes an OPTIONS request
func (r *Request) Options(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodOptions
	return r.Execute()
}

// Execute executes the request
func (r *Request) Execute() (*Response, error) {
	return r.client.execute(r)
}

// Do is an alias for Execute
func (r *Request) Do() (*Response, error) {
	return r.Execute()
}

// Send is an alias for Execute
func (r *Request) Send() (*Response, error) {
	return r.Execute()
}

// MustGet executes a GET request and panics on error
func (r *Request) MustGet(url ...string) *Response {
	resp, err := r.Get(url...)
	if err != nil {
		panic(err)
	}
	return resp
}

// MustPost executes a POST request and panics on error
func (r *Request) MustPost(url ...string) *Response {
	resp, err := r.Post(url...)
	if err != nil {
		panic(err)
	}
	return resp
}

// MustPut executes a PUT request and panics on error
func (r *Request) MustPut(url ...string) *Response {
	resp, err := r.Put(url...)
	if err != nil {
		panic(err)
	}
	return resp
}

// MustPatch executes a PATCH request and panics on error
func (r *Request) MustPatch(url ...string) *Response {
	resp, err := r.Patch(url...)
	if err != nil {
		panic(err)
	}
	return resp
}

// MustDelete executes a DELETE request and panics on error
func (r *Request) MustDelete(url ...string) *Response {
	resp, err := r.Delete(url...)
	if err != nil {
		panic(err)
	}
	return resp
}

// MustHead executes a HEAD request and panics on error
func (r *Request) MustHead(url ...string) *Response {
	resp, err := r.Head(url...)
	if err != nil {
		panic(err)
	}
	return resp
}

// MustOptions executes an OPTIONS request and panics on error
func (r *Request) MustOptions(url ...string) *Response {
	resp, err := r.Options(url...)
	if err != nil {
		panic(err)
	}
	return resp
}

// MustExecute executes the request and panics on error
func (r *Request) MustExecute() *Response {
	resp, err := r.Execute()
	if err != nil {
		panic(err)
	}
	return resp
}

// Clone creates a copy of the request
func (r *Request) Clone() *Request {
	headers := make(http.Header)
	for k, v := range r.headers {
		headers[k] = append([]string(nil), v...)
	}

	queryParams := make(url.Values)
	for k, v := range r.queryParams {
		queryParams[k] = append([]string(nil), v...)
	}

	pathParams := make(map[string]string)
	for k, v := range r.pathParams {
		pathParams[k] = v
	}

	formData := make(url.Values)
	for k, v := range r.formData {
		formData[k] = append([]string(nil), v...)
	}

	cookies := make([]*http.Cookie, len(r.cookies))
	copy(cookies, r.cookies)

	return &Request{
		client:         r.client,
		method:         r.method,
		url:            r.url,
		ctx:            r.ctx,
		headers:        headers,
		queryParams:    queryParams,
		pathParams:     pathParams,
		formData:       formData,
		body:           r.body,
		bodyType:       r.bodyType,
		cookies:        cookies,
		basicAuth:      r.basicAuth,
		bearerToken:    r.bearerToken,
		successResult:  r.successResult,
		errorResult:    r.errorResult,
		downloadPath:   r.downloadPath,
		uploadCallback: r.uploadCallback,
	}
}

// URL returns the final request URL (after path parameter replacement)
func (r *Request) URL() string {
	u, err := r.client.buildURL(r.url, r.pathParams, r.queryParams)
	if err != nil {
		return r.url
	}
	return u.String()
}

// Method returns the HTTP method
func (r *Request) Method() string {
	return r.method
}

// Header returns the request headers
func (r *Request) Header() http.Header {
	return r.headers
}

// Validate validates the request
func (r *Request) Validate() error {
	if r.method == "" {
		return fmt.Errorf("HTTP method is required")
	}
	if r.url == "" {
		return fmt.Errorf("URL is required")
	}
	return nil
}

// String returns a string representation of the request
func (r *Request) String() string {
	return fmt.Sprintf("%s %s", r.method, r.URL())
}
