package cumi

import (
	"net/url"
	"strings"
)

// defaultResultChecker checks the state of the response based on status code
func defaultResultChecker(resp *Response) ResultState {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return SuccessState
	}
	if resp.StatusCode >= 400 {
		return ErrorState
	}
	return UnknownState
}

// buildURL builds the final URL with base URL, path params, and query params
func (c *Client) buildURL(rawURL string, pathParams map[string]string, queryParams url.Values) (*url.URL, error) {
	finalURL := rawURL

	// Add base URL if relative
	if !strings.HasPrefix(rawURL, "http") && c.baseURL != "" {
		finalURL = c.baseURL + "/" + strings.TrimLeft(rawURL, "/")
	}

	// Replace path parameters
	allPathParams := make(map[string]string)
	for k, v := range c.pathParams {
		allPathParams[k] = v
	}
	for k, v := range pathParams {
		allPathParams[k] = v
	}

	for key, value := range allPathParams {
		placeholder := "{" + key + "}"
		finalURL = strings.ReplaceAll(finalURL, placeholder, value)
	}

	u, err := url.Parse(finalURL)
	if err != nil {
		return nil, err
	}

	// Merge query parameters
	q := u.Query()
	for k, values := range c.queryParams {
		for _, v := range values {
			q.Add(k, v)
		}
	}
	for k, values := range queryParams {
		for _, v := range values {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()

	return u, nil
}

// shouldRetry determines if a request should be retried based on response and error
func (c *Client) shouldRetry(resp *Response, err error) bool {
	if c.retryCondition != nil {
		return c.retryCondition(resp, err)
	}

	// Default retry logic
	if err != nil {
		return true // Retry on network errors
	}

	if resp != nil && (resp.StatusCode >= 500 || resp.StatusCode == 429) {
		return true // Retry on server errors and rate limiting
	}

	return false
}

// unmarshalResponse unmarshals the response body into the given interface
func (c *Client) unmarshalResponse(resp *Response, v interface{}) error {
	if len(resp.body) == 0 {
		return nil
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		return c.jsonUnmarshal(resp.body, v)
	} else if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		return c.xmlUnmarshal(resp.body, v)
	}

	// Default to JSON
	return c.jsonUnmarshal(resp.body, v)
}
