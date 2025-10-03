package cumi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGetRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Http().Get(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if !resp.IsSuccess() {
		t.Errorf("Expected success response")
	}
}

func TestPostRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "created"})
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Http().
		SetBodyJSON(map[string]string{"name": "John"}).
		Post(server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestPutRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Http().
		SetBodyJSON(map[string]string{"name": "John"}).
		Put(server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestPatchRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Http().
		SetBodyJSON(map[string]string{"name": "John"}).
		Patch(server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestDeleteRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Http().Delete(server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestClientConfiguration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient().SetBaseURL(server.URL)
	resp, err := client.Http().Get("/test")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestSetSuccessResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := User{Name: "John", Age: 30}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	var user User
	client := NewClient()
	resp, err := client.Http().
		SetSuccessResult(&user).
		Get(server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if user.Name != "John" {
		t.Errorf("Expected name=John, got %s", user.Name)
	}

	if user.Age != 30 {
		t.Errorf("Expected age=30, got %d", user.Age)
	}
}

func TestErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Http().Get(server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}

	if resp.IsError() {
		t.Logf("Response correctly identified as error: %s", resp.Status)
	}
}

func TestJSONUnmarshaling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := User{Name: "Jane", Age: 25}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Http().Get(server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var user User
	err = resp.JSON(&user)
	if err != nil {
		t.Fatalf("Expected no JSON error, got %v", err)
	}

	if user.Name != "Jane" {
		t.Errorf("Expected name=Jane, got %s", user.Name)
	}

	if user.Age != 25 {
		t.Errorf("Expected age=25, got %d", user.Age)
	}
}

func TestUserAgentDefault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"user_agent": userAgent})
	}))
	defer server.Close()

	// Test 1: Config with empty User-Agent should use Go default
	config1 := &Config{
		BaseURL: server.URL,
		// UserAgent is empty
	}
	client1 := NewClientWithConfig(config1)
	resp1, err := client1.Get("/").Execute()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result1 map[string]string
	if err := json.Unmarshal(resp1.Body(), &result1); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result1["user_agent"] != "Go-http-client/1.1" {
		t.Errorf("Expected User-Agent 'Go-http-client/1.1', got '%s'", result1["user_agent"])
	}

	// Test 2: Config with custom User-Agent
	config2 := &Config{
		BaseURL:   server.URL,
		UserAgent: "CustomApp/1.0",
	}
	client2 := NewClientWithConfig(config2)
	resp2, err := client2.Get("/").Execute()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result2 map[string]string
	if err := json.Unmarshal(resp2.Body(), &result2); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result2["user_agent"] != "CustomApp/1.0" {
		t.Errorf("Expected User-Agent 'CustomApp/1.0', got '%s'", result2["user_agent"])
	}

	// Test 3: DefaultConfig should have "Go-http-client/1.1"
	config3 := DefaultConfig()
	config3.BaseURL = server.URL
	client3 := NewClientWithConfig(config3)
	resp3, err := client3.Get("/").Execute()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result3 map[string]string
	if err := json.Unmarshal(resp3.Body(), &result3); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result3["user_agent"] != "Go-http-client/1.1" {
		t.Errorf("Expected User-Agent 'Go-http-client/1.1', got '%s'", result3["user_agent"])
	}
}

func TestUserAgentPriority(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"user_agent": userAgent})
	}))
	defer server.Close()

	// Test 1: Request SetUserAgent should override config UserAgent
	config := &Config{
		BaseURL:   server.URL,
		UserAgent: "ConfigAgent/1.0",
	}
	client := NewClientWithConfig(config)
	resp, err := client.Http().SetUserAgent("RequestAgent/1.0").Get("/")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["user_agent"] != "RequestAgent/1.0" {
		t.Errorf("Expected User-Agent 'RequestAgent/1.0', got '%s'", result["user_agent"])
	}

	// Test 2: No SetUserAgent, should use config UserAgent
	resp2, err := client.Http().Get("/")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result2 map[string]string
	if err := json.Unmarshal(resp2.Body(), &result2); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result2["user_agent"] != "ConfigAgent/1.0" {
		t.Errorf("Expected User-Agent 'ConfigAgent/1.0', got '%s'", result2["user_agent"])
	}

	// Test 3: No config UserAgent, no request UserAgent, should use Go default
	emptyConfig := &Config{
		BaseURL: server.URL,
		// No UserAgent
	}
	client3 := NewClientWithConfig(emptyConfig)
	resp3, err := client3.Http().Get("/")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result3 map[string]string
	if err := json.Unmarshal(resp3.Body(), &result3); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result3["user_agent"] != "Go-http-client/1.1" {
		t.Errorf("Expected User-Agent 'Go-http-client/1.1', got '%s'", result3["user_agent"])
	}
}

func TestDefaultContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"content_type": contentType})
	}))
	defer server.Close()

	// Test 1: No explicit Content-Type should default to application/json
	client := NewClient()
	resp, err := client.Http().Get(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["content_type"] != "application/json" {
		t.Errorf("Expected default Content-Type 'application/json', got '%s'", result["content_type"])
	}

	// Test 2: Explicit Content-Type should override default
	resp2, err := client.Http().
		SetHeader("Content-Type", "text/plain").
		Get(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result2 map[string]string
	if err := json.Unmarshal(resp2.Body(), &result2); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result2["content_type"] != "text/plain" {
		t.Errorf("Expected Content-Type 'text/plain', got '%s'", result2["content_type"])
	}
}
