package cumi

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// Client represents an HTTP client with chainable methods
type Client struct {
	httpClient        *http.Client
	baseURL           string
	timeout           time.Duration
	headers           http.Header
	queryParams       url.Values
	pathParams        map[string]string
	formData          url.Values
	cookies           []*http.Cookie
	userAgent         string
	beforeRequest     []RequestMiddleware
	afterResponse     []ResponseMiddleware
	jsonMarshal       func(v interface{}) ([]byte, error)
	jsonUnmarshal     func(data []byte, v interface{}) error
	xmlMarshal        func(v interface{}) ([]byte, error)
	xmlUnmarshal      func(data []byte, v interface{}) error
	debug             bool
	allowGetPayload   bool
	retryCount        int
	retryInterval     time.Duration
	retryCondition    RetryConditionFunc
	errorHandler      ErrorHook
	onError           ErrorHook
	commonErrorResult interface{}
	resultChecker     func(*Response) ResultState
}

// RequestMiddleware defines a function that can modify a request before it's sent
type RequestMiddleware func(*Client, *Request) error

// ResponseMiddleware defines a function that can modify a response after it's received
type ResponseMiddleware func(*Client, *Response) error

// RetryConditionFunc defines when a request should be retried
type RetryConditionFunc func(*Response, error) bool

// ErrorHook is called when an error occurs
type ErrorHook func(*Client, *Request, *Response, error)

// ResultState represents the state of the response
type ResultState int

const (
	SuccessState ResultState = iota
	ErrorState
	UnknownState
)

// C creates a new client (alias for NewClient)
func C() *Client {
	return NewClient()
}

// NewClient creates a new HTTP client with default settings
func NewClient() *Client {
	return NewClientWithConfig(DefaultConfig())
}

// NewClientWithConfig creates a new HTTP client with provided configuration
func NewClientWithConfig(config *Config) *Client {
	jar, _ := cookiejar.New(nil)

	// Use config's transport or create default
	var transport http.RoundTripper
	if config.Transport != nil {
		transport = config.Transport
	} else {
		tlsConfig := config.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{}
		}
		transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	httpClient := &http.Client{
		Timeout:   config.Timeout,
		Jar:       jar,
		Transport: transport,
	}

	// Convert headers map to http.Header
	headers := make(http.Header)
	for k, v := range config.Headers {
		headers.Set(k, v)
	}

	// Convert query params map to url.Values
	queryParams := make(url.Values)
	for k, v := range config.QueryParams {
		queryParams.Set(k, v)
	}

	// Convert path params map
	pathParams := make(map[string]string)
	for k, v := range config.PathParams {
		pathParams[k] = v
	}

	// Set default User-Agent if empty
	userAgent := config.UserAgent
	if userAgent == "" {
		userAgent = "Go-http-client/1.1" // Default Go HTTP client User-Agent
	}

	// Set default resultChecker if nil
	resultChecker := config.ResultChecker
	if resultChecker == nil {
		resultChecker = defaultResultChecker
	}

	// Set default timeout if zero
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	// Ensure headers, queryParams, pathParams are not nil
	if config.Headers == nil {
		config.Headers = make(map[string]string)
	}
	if config.QueryParams == nil {
		config.QueryParams = make(map[string]string)
	}
	if config.PathParams == nil {
		config.PathParams = make(map[string]string)
	}
	if config.BeforeRequest == nil {
		config.BeforeRequest = []RequestMiddleware{}
	}
	if config.AfterResponse == nil {
		config.AfterResponse = []ResponseMiddleware{}
	}

	c := &Client{
		httpClient:        httpClient,
		baseURL:           config.BaseURL,
		timeout:           timeout,
		headers:           headers,
		queryParams:       queryParams,
		pathParams:        pathParams,
		formData:          make(url.Values),
		userAgent:         userAgent,
		debug:             config.Debug,
		allowGetPayload:   config.AllowGetPayload,
		retryCount:        config.RetryCount,
		retryInterval:     config.RetryInterval,
		retryCondition:    config.RetryCondition,
		errorHandler:      config.ErrorHandler,
		onError:           config.OnError,
		commonErrorResult: config.CommonErrorResult,
		resultChecker:     resultChecker,
		jsonMarshal:       json.Marshal,
		jsonUnmarshal:     json.Unmarshal,
		xmlMarshal:        xml.Marshal,
		xmlUnmarshal:      xml.Unmarshal,
		beforeRequest:     append([]RequestMiddleware{}, config.BeforeRequest...),
		afterResponse:     append([]ResponseMiddleware{}, config.AfterResponse...),
	}

	return c
}

// R creates a new request
func (c *Client) Http() *Request {
	return &Request{
		client:      c,
		headers:     make(http.Header),
		queryParams: make(url.Values),
		pathParams:  make(map[string]string),
		formData:    make(url.Values),
		cookies:     []*http.Cookie{},
		ctx:         context.Background(),
	}
}

// Get creates a new GET request
func (c *Client) Get(url ...string) *Request {
	r := c.Http()
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodGet
	return r
}

// Post creates a new POST request
func (c *Client) Post(url ...string) *Request {
	r := c.Http()
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodPost
	return r
}

// Put creates a new PUT request
func (c *Client) Put(url ...string) *Request {
	r := c.Http()
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodPut
	return r
}

// Patch creates a new PATCH request
func (c *Client) Patch(url ...string) *Request {
	r := c.Http()
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodPatch
	return r
}

// Delete creates a new DELETE request
func (c *Client) Delete(url ...string) *Request {
	r := c.Http()
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodDelete
	return r
}

// Head creates a new HEAD request
func (c *Client) Head(url ...string) *Request {
	r := c.Http()
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodHead
	return r
}

// Options creates a new OPTIONS request
func (c *Client) Options(url ...string) *Request {
	r := c.Http()
	if len(url) > 0 {
		r.url = url[0]
	}
	r.method = http.MethodOptions
	return r
}

// Clone creates a copy of the client
func (c *Client) Clone() *Client {
	jar, _ := cookiejar.New(nil)

	transport := &http.Transport{}
	if t, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport = t.Clone()
	}

	httpClient := &http.Client{
		Timeout:   c.httpClient.Timeout,
		Jar:       jar,
		Transport: transport,
	}

	headers := make(http.Header)
	for k, v := range c.headers {
		headers[k] = append([]string(nil), v...)
	}

	queryParams := make(url.Values)
	for k, v := range c.queryParams {
		queryParams[k] = append([]string(nil), v...)
	}

	pathParams := make(map[string]string)
	for k, v := range c.pathParams {
		pathParams[k] = v
	}

	formData := make(url.Values)
	for k, v := range c.formData {
		formData[k] = append([]string(nil), v...)
	}

	cookies := make([]*http.Cookie, len(c.cookies))
	copy(cookies, c.cookies)

	return &Client{
		httpClient:        httpClient,
		baseURL:           c.baseURL,
		timeout:           c.timeout,
		headers:           headers,
		queryParams:       queryParams,
		pathParams:        pathParams,
		formData:          formData,
		cookies:           cookies,
		userAgent:         c.userAgent,
		beforeRequest:     append([]RequestMiddleware(nil), c.beforeRequest...),
		afterResponse:     append([]ResponseMiddleware(nil), c.afterResponse...),
		jsonMarshal:       c.jsonMarshal,
		jsonUnmarshal:     c.jsonUnmarshal,
		xmlMarshal:        c.xmlMarshal,
		xmlUnmarshal:      c.xmlUnmarshal,
		debug:             c.debug,
		allowGetPayload:   c.allowGetPayload,
		retryCount:        c.retryCount,
		retryInterval:     c.retryInterval,
		retryCondition:    c.retryCondition,
		errorHandler:      c.errorHandler,
		onError:           c.onError,
		commonErrorResult: c.commonErrorResult,
		resultChecker:     c.resultChecker,
	}
}

// SetBaseURL sets the base URL for the client
func (c *Client) SetBaseURL(baseURL string) *Client {
	c.baseURL = strings.TrimRight(baseURL, "/")
	return c
}

// SetTimeout sets the request timeout
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	c.httpClient.Timeout = timeout
	return c
}

// SetUserAgent sets the User-Agent header for all requests
func (c *Client) SetUserAgent(userAgent string) *Client {
	c.userAgent = userAgent
	return c
}

// SetCommonHeader sets a header that will be added to all requests
func (c *Client) SetCommonHeader(key, value string) *Client {
	c.headers.Set(key, value)
	return c
}

// SetCommonHeaders sets multiple headers from a map
func (c *Client) SetCommonHeaders(headers map[string]string) *Client {
	for k, v := range headers {
		c.headers.Set(k, v)
	}
	return c
}

// SetCommonQueryParam sets a query parameter that will be added to all requests
func (c *Client) SetCommonQueryParam(key, value string) *Client {
	c.queryParams.Set(key, value)
	return c
}

// SetCommonQueryParams sets multiple query parameters from a map
func (c *Client) SetCommonQueryParams(params map[string]string) *Client {
	for k, v := range params {
		c.queryParams.Set(k, v)
	}
	return c
}

// SetCommonPathParam sets a path parameter that will be used for URL replacement
func (c *Client) SetCommonPathParam(key, value string) *Client {
	if c.pathParams == nil {
		c.pathParams = make(map[string]string)
	}
	c.pathParams[key] = value
	return c
}

// SetCommonPathParams sets multiple path parameters from a map
func (c *Client) SetCommonPathParams(params map[string]string) *Client {
	if c.pathParams == nil {
		c.pathParams = make(map[string]string)
	}
	for k, v := range params {
		c.pathParams[k] = v
	}
	return c
}

// SetCommonFormData sets form data that will be added to all requests
func (c *Client) SetCommonFormData(data map[string]string) *Client {
	for k, v := range data {
		c.formData.Set(k, v)
	}
	return c
}

// SetCommonCookies sets cookies that will be added to all requests
func (c *Client) SetCommonCookies(cookies ...*http.Cookie) *Client {
	c.cookies = append(c.cookies, cookies...)
	return c
}

// EnableDebug enables debug mode
func (c *Client) EnableDebug() *Client {
	c.debug = true
	return c
}

// DisableDebug disables debug mode
func (c *Client) DisableDebug() *Client {
	c.debug = false
	return c
}

// DevMode enables debug mode (alias for EnableDebug)
func (c *Client) DevMode() *Client {
	return c.EnableDebug()
}

// EnableAllowGetMethodPayload allows GET requests to have a body
func (c *Client) EnableAllowGetMethodPayload() *Client {
	c.allowGetPayload = true
	return c
}

// DisableAllowGetMethodPayload disallows GET requests to have a body
func (c *Client) DisableAllowGetMethodPayload() *Client {
	c.allowGetPayload = false
	return c
}

// SetTLSClientConfig sets the TLS configuration
func (c *Client) SetTLSClientConfig(config *tls.Config) *Client {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport.TLSClientConfig = config
	}
	return c
}

// EnableInsecureSkipVerify enables skipping TLS certificate verification
func (c *Client) EnableInsecureSkipVerify() *Client {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.InsecureSkipVerify = true
	}
	return c
}

// DisableInsecureSkipVerify disables skipping TLS certificate verification
func (c *Client) DisableInsecureSkipVerify() *Client {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.InsecureSkipVerify = false
	}
	return c
}

// SetProxy sets the proxy function
func (c *Client) SetProxy(proxy func(*http.Request) (*url.URL, error)) *Client {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport.Proxy = proxy
	}
	return c
}

// SetRetryCount sets the number of retry attempts
func (c *Client) SetRetryCount(count int) *Client {
	c.retryCount = count
	return c
}

// SetRetryInterval sets the interval between retries
func (c *Client) SetRetryInterval(interval time.Duration) *Client {
	c.retryInterval = interval
	return c
}

// SetRetryCondition sets the condition for when to retry
func (c *Client) SetRetryCondition(condition RetryConditionFunc) *Client {
	c.retryCondition = condition
	return c
}

// SetCommonErrorResult sets the common error result type
func (c *Client) SetCommonErrorResult(err interface{}) *Client {
	c.commonErrorResult = err
	return c
}

// SetResultStateCheckFunc sets the function to check result state
func (c *Client) SetResultStateCheckFunc(fn func(*Response) ResultState) *Client {
	c.resultChecker = fn
	return c
}

// OnError sets the error handler
func (c *Client) OnError(handler ErrorHook) *Client {
	c.onError = handler
	return c
}

// OnBeforeRequest adds a middleware that runs before sending the request
func (c *Client) OnBeforeRequest(middleware RequestMiddleware) *Client {
	c.beforeRequest = append(c.beforeRequest, middleware)
	return c
}

// OnAfterResponse adds a middleware that runs after receiving the response
func (c *Client) OnAfterResponse(middleware ResponseMiddleware) *Client {
	c.afterResponse = append(c.afterResponse, middleware)
	return c
}

// SetJSONMarshal sets the JSON marshal function
func (c *Client) SetJSONMarshal(fn func(v interface{}) ([]byte, error)) *Client {
	c.jsonMarshal = fn
	return c
}

// SetJSONUnmarshal sets the JSON unmarshal function
func (c *Client) SetJSONUnmarshal(fn func(data []byte, v interface{}) error) *Client {
	c.jsonUnmarshal = fn
	return c
}

// SetXMLMarshal sets the XML marshal function
func (c *Client) SetXMLMarshal(fn func(v interface{}) ([]byte, error)) *Client {
	c.xmlMarshal = fn
	return c
}

// SetXMLUnmarshal sets the XML unmarshal function
func (c *Client) SetXMLUnmarshal(fn func(data []byte, v interface{}) error) *Client {
	c.xmlUnmarshal = fn
	return c
}

// GetClient returns the underlying http.Client
func (c *Client) GetClient() *http.Client {
	return c.httpClient
}

// GetTLSClientConfig returns the TLS configuration
func (c *Client) GetTLSClientConfig() *tls.Config {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		return transport.TLSClientConfig
	}
	return nil
}

// prepareRequest prepares the HTTP request
func (c *Client) prepareRequest(req *Request) (*http.Request, error) {
	// Build URL
	u, err := c.buildURL(req.url, req.pathParams, req.queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// Prepare body
	var body io.Reader
	var contentType string

	if req.body != nil {
		if req.bodyType == "json" {
			jsonData, err := c.jsonMarshal(req.body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal JSON: %w", err)
			}
			body = bytes.NewReader(jsonData)
			contentType = "application/json"
		} else if req.bodyType == "xml" {
			xmlData, err := c.xmlMarshal(req.body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal XML: %w", err)
			}
			body = bytes.NewReader(xmlData)
			contentType = "application/xml"
		} else if data, ok := req.body.([]byte); ok {
			body = bytes.NewReader(data)
		} else if s, ok := req.body.(string); ok {
			body = strings.NewReader(s)
		} else if r, ok := req.body.(io.Reader); ok {
			body = r
		}
	} else if len(req.formData) > 0 || len(c.formData) > 0 {
		// Merge form data
		formData := make(url.Values)
		for k, values := range c.formData {
			for _, v := range values {
				formData.Add(k, v)
			}
		}
		for k, values := range req.formData {
			for _, v := range values {
				formData.Add(k, v)
			}
		}
		body = strings.NewReader(formData.Encode())
		contentType = "application/x-www-form-urlencoded"
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(req.ctx, req.method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	for k, values := range c.headers {
		for _, v := range values {
			httpReq.Header.Add(k, v)
		}
	}
	for k, values := range req.headers {
		for _, v := range values {
			httpReq.Header.Add(k, v)
		}
	}

	// Set User-Agent with priority: Request > Client Config > Default Go
	if httpReq.Header.Get("User-Agent") == "" {
		var userAgent string
		if req.userAgent != "" {
			// Priority 1: Request-specific User-Agent
			userAgent = req.userAgent
		} else if c.userAgent != "" {
			// Priority 2: Client config User-Agent
			userAgent = c.userAgent
		} else {
			// Priority 3: Default Go HTTP client User-Agent
			userAgent = "Go-http-client/1.1"
		}
		httpReq.Header.Set("User-Agent", userAgent)
	}

	// Set content type if not already set
	if httpReq.Header.Get("Content-Type") == "" {
		// Use content type determined by body type (JSON, XML, form data)
		httpReq.Header.Set("Content-Type", contentType)
	}

	// Set basic auth
	if req.basicAuth.username != "" {
		httpReq.SetBasicAuth(req.basicAuth.username, req.basicAuth.password)
	}

	// Set bearer token
	if req.bearerToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+req.bearerToken)
	}

	// Add cookies
	for _, cookie := range c.cookies {
		httpReq.AddCookie(cookie)
	}
	for _, cookie := range req.cookies {
		httpReq.AddCookie(cookie)
	}

	return httpReq, nil
}

// execute performs the actual HTTP request with retry logic
func (c *Client) execute(req *Request) (*Response, error) {
	var lastErr error
	var resp *Response

	maxAttempts := c.retryCount + 1
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Prepare the HTTP request
		httpReq, err := c.prepareRequest(req)
		if err != nil {
			return nil, err
		}

		// Debug: Print request details
		if c.debug {
			c.debugRequest(httpReq, attempt+1, maxAttempts)
		}

		// Run before request middlewares
		for _, middleware := range c.beforeRequest {
			if err := middleware(c, req); err != nil {
				return nil, fmt.Errorf("before request middleware error: %w", err)
			}
		}

		// Execute the request
		startTime := time.Now()
		httpResp, err := c.httpClient.Do(httpReq)
		duration := time.Since(startTime)

		// Create response
		resp = &Response{
			Request:    req,
			Response:   httpResp,
			receivedAt: time.Now(),
			duration:   duration,
		}

		if err != nil {
			lastErr = err
			resp.Err = err

			// Check if we should retry
			if attempt < maxAttempts-1 && c.shouldRetry(resp, err) {
				time.Sleep(c.retryInterval)
				continue
			}
			break
		}

		// Read response body
		if httpResp.Body != nil {
			defer httpResp.Body.Close()
			bodyBytes, err := io.ReadAll(httpResp.Body)
			if err != nil {
				resp.Err = fmt.Errorf("failed to read response body: %w", err)
				lastErr = resp.Err
				if attempt < maxAttempts-1 && c.shouldRetry(resp, resp.Err) {
					time.Sleep(c.retryInterval)
					continue
				}
				break
			}
			resp.body = bodyBytes
			resp.size = int64(len(bodyBytes))
		}

		// Copy status information
		if httpResp != nil {
			resp.StatusCode = httpResp.StatusCode
			resp.Status = httpResp.Status
			resp.Proto = httpResp.Proto
			resp.ProtoMajor = httpResp.ProtoMajor
			resp.ProtoMinor = httpResp.ProtoMinor
			resp.Header = httpResp.Header
		}

		// Run after response middlewares
		for _, middleware := range c.afterResponse {
			if err := middleware(c, resp); err != nil {
				resp.Err = fmt.Errorf("after response middleware error: %w", err)
				lastErr = resp.Err
				if attempt < maxAttempts-1 && c.shouldRetry(resp, resp.Err) {
					time.Sleep(c.retryInterval)
					continue
				}
				break
			}
		}

		// Unmarshal success/error results
		if resp.Err == nil {
			resp.state = c.resultChecker(resp)

			if resp.state == SuccessState && req.successResult != nil {
				if err := c.unmarshalResponse(resp, req.successResult); err != nil {
					resp.Err = fmt.Errorf("failed to unmarshal success result: %w", err)
				}
			} else if resp.state == ErrorState {
				if req.errorResult != nil {
					c.unmarshalResponse(resp, req.errorResult)
				} else if c.commonErrorResult != nil {
					c.unmarshalResponse(resp, c.commonErrorResult)
				}
			}
		}

		// Debug: Print response details
		if c.debug {
			c.debugResponse(resp)
		}

		// Check if we should retry
		if attempt < maxAttempts-1 && c.shouldRetry(resp, resp.Err) {
			if c.debug {
				log.Printf("[DEBUG] RETRY - Retrying in %v...", c.retryInterval)
			}
			time.Sleep(c.retryInterval)
			continue
		}

		break
	}

	// Call error handler if there's an error
	if resp != nil && resp.Err != nil && c.onError != nil {
		c.onError(c, req, resp, resp.Err)
	}

	if resp == nil && lastErr != nil {
		return nil, lastErr
	}

	return resp, resp.Err
}

// debugRequest prints debug information for the request
func (c *Client) debugRequest(req *http.Request, attempt, maxAttempts int) {
	log.Printf("[DEBUG] REQUEST - Attempt: %d/%d, Method: %s, URL: %s", attempt, maxAttempts, req.Method, req.URL.String())

	for key, values := range req.Header {
		for _, value := range values {
			log.Printf("[DEBUG] REQUEST Header - %s: %s", key, value)
		}
	}

	if req.Body != nil {
		// Try to read body for debug (this won't consume the original body)
		if req.GetBody != nil {
			if body, err := req.GetBody(); err == nil {
				if bodyBytes, err := io.ReadAll(body); err == nil && len(bodyBytes) > 0 {
					bodyStr := string(bodyBytes)
					if len(bodyStr) > 300 {
						bodyStr = bodyStr[:300] + "...(truncated)"
					}
					log.Printf("[DEBUG] REQUEST Body - %s", bodyStr)
				}
				body.Close()
			}
		}
	}
} // debugResponse prints debug information for the response
func (c *Client) debugResponse(resp *Response) {
	log.Printf("[DEBUG] RESPONSE - Status: %s (%d), Duration: %v, Size: %d bytes",
		resp.Status, resp.StatusCode, resp.Duration(), resp.Size())

	for key, values := range resp.Header {
		for _, value := range values {
			log.Printf("[DEBUG] RESPONSE Header - %s: %s", key, value)
		}
	}

	if len(resp.body) > 0 {
		// Limit body display to first 300 characters
		bodyStr := string(resp.body)
		if len(bodyStr) > 300 {
			bodyStr = bodyStr[:300] + "...(truncated)"
		}
		log.Printf("[DEBUG] RESPONSE Body - %s", bodyStr)
	}
}
