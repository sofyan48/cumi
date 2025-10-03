package cumi

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Config holds default configuration for Client
type Config struct {
	BaseURL           string
	Timeout           time.Duration
	Headers           map[string]string
	QueryParams       map[string]string
	PathParams        map[string]string
	UserAgent         string
	Debug             bool
	AllowGetPayload   bool
	RetryCount        int
	RetryInterval     time.Duration
	TLSConfig         *tls.Config
	Transport         http.RoundTripper
	BeforeRequest     []RequestMiddleware
	AfterResponse     []ResponseMiddleware
	RetryCondition    RetryConditionFunc
	ErrorHandler      ErrorHook
	OnError           ErrorHook
	CommonErrorResult interface{}
	ResultChecker     func(*Response) ResultState
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Timeout:         30 * time.Second,
		UserAgent:       "Go-http-client/1.1",
		Debug:           false,
		AllowGetPayload: false,
		RetryCount:      0,
		RetryInterval:   time.Second,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		QueryParams:   make(map[string]string),
		PathParams:    make(map[string]string),
		BeforeRequest: []RequestMiddleware{},
		AfterResponse: []ResponseMiddleware{},
		ResultChecker: defaultResultChecker,
	}
}
