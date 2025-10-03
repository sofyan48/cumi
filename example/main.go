package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sofyan48/cumi"
)

// User represents a user object for JSON parsing examples
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// Post represents a post object for JSON parsing examples
type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// APIError represents an error response from API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func main() {
	fmt.Println("=== Cumi HTTP Client Examples ===")

	// Example 1: Basic GET Request
	basicGetExample()

	// Example 2: GET with Auto JSON Parsing
	getWithJSONParsingExample()

	// Example 3: POST with JSON Body
	postWithJSONExample()

	// Example 4: POST with Auto Result Binding
	postWithResultBindingExample()

	// Example 5: Error Handling with SetErrorResult
	errorHandlingExample()

	// Example 6: Client Configuration
	clientConfigurationExample()

	// Example 7: Authentication Examples
	authenticationExamples()

	// Example 8: Retry and Timeout Examples
	retryAndTimeoutExample()

	// Example 9: Context and Cancellation
	contextAndCancellationExample()

	// Example 10: Middleware Examples
	middlewareExample()

	// Example 11: Query Parameters
	queryParamsExample()

	// Example 12: Form Data Example
	formDataExample()

	// Example 13: Cookie Management
	cookieExample()

	// Example 14: Debug Mode
	debugModeExample()

	fmt.Println("\n=== All Examples Completed ===")
}

// Example 1: Basic GET Request
func basicGetExample() {
	fmt.Println("1. Basic GET Request:")

	client := cumi.NewClient()
	resp, err := client.Http().Get("https://httpbin.org/get")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("   Status: %s\n", resp.Status)
	fmt.Printf("   Response length: %d bytes\n", len(resp.String()))
	fmt.Println()
}

// Example 2: GET with Auto JSON Parsing
func getWithJSONParsingExample() {
	fmt.Println("2. GET with Auto JSON Parsing:")

	client := cumi.NewClient()
	var user User
	resp, err := client.Http().
		SetSuccessResult(&user).
		Get("https://jsonplaceholder.typicode.com/users/1")

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp.IsSuccess() {
		fmt.Printf("   User: %+v\n", user)
	} else {
		fmt.Printf("   Request failed with status: %s\n", resp.Status)
	}
	fmt.Println()
}

// Example 3: POST with JSON Body
func postWithJSONExample() {
	fmt.Println("3. POST with JSON Body:")

	client := cumi.NewClient()
	user := User{Name: "John Doe", Email: "john@example.com", Phone: "123-456-7890"}

	resp, err := client.Http().
		SetBodyJSON(user).
		Post("https://httpbin.org/post")
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("   Status: %s\n", resp.Status)
	fmt.Printf("   Response length: %d bytes\n", len(resp.String()))
	fmt.Println()
}

// Example 4: POST with Auto Result Binding
func postWithResultBindingExample() {
	fmt.Println("4. POST with Auto Result Binding:")

	client := cumi.NewClient()
	user := User{Name: "Jane Doe", Email: "jane@example.com", Phone: "098-765-4321"}

	// Auto result binding with SetSuccessResult
	var result map[string]interface{}
	resp, err := client.Http().
		SetBodyJSON(user).
		SetSuccessResult(&result).
		Post("https://httpbin.org/post")

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp.IsSuccess() {
		if jsonData, ok := result["json"].(map[string]interface{}); ok {
			fmt.Printf("   Received JSON: %+v\n", jsonData)
		}
	}
	fmt.Println()
}

// Example 5: Error Handling with SetErrorResult
func errorHandlingExample() {
	fmt.Println("5. Error Handling with SetErrorResult:")

	client := cumi.NewClient()
	var user User
	var apiError APIError

	resp, err := client.Http().
		SetSuccessResult(&user).
		SetErrorResult(&apiError).
		Get("https://jsonplaceholder.typicode.com/users/999") // Non-existent user

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp.IsSuccess() {
		fmt.Printf("   User: %+v\n", user)
	} else {
		fmt.Printf("   Error Status: %s\n", resp.Status)
		fmt.Printf("   Error Response: %s\n", resp.String())
	}
	fmt.Println()
}

// Example 6: Client Configuration
func clientConfigurationExample() {
	fmt.Println("6. Client Configuration:")

	client := cumi.NewClient().
		SetBaseURL("https://jsonplaceholder.typicode.com").
		SetTimeout(30*time.Second).
		SetCommonHeader("User-Agent", "Cumi-Example/1.0").
		SetCommonHeader("Accept", "application/json").
		SetRetryCount(3)

	var posts []Post
	resp, err := client.Http().
		SetSuccessResult(&posts).
		Get("/posts?_limit=5")

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp.IsSuccess() {
		fmt.Printf("   Retrieved %d posts\n", len(posts))
		if len(posts) > 0 {
			fmt.Printf("   First post: %+v\n", posts[0])
		}
	}
	fmt.Println()
}

// Example 7: Authentication Examples
func authenticationExamples() {
	fmt.Println("7. Authentication Examples:")

	// Basic Auth Example
	fmt.Println("   a) Basic Auth:")
	client := cumi.NewClient()
	resp, err := client.Http().
		SetBasicAuth("user", "pass").
		Get("https://httpbin.org/basic-auth/user/pass")

	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   Basic Auth Status: %s\n", resp.Status)
	}

	// Bearer Token Example
	fmt.Println("   b) Bearer Token:")
	resp2, err2 := client.Http().
		SetBearerToken("your-jwt-token-here").
		Get("https://httpbin.org/bearer")

	if err2 != nil {
		log.Printf("   Error: %v", err2)
	} else {
		fmt.Printf("   Bearer Token Status: %s\n", resp2.Status)
	}

	// API Key Example
	fmt.Println("   c) API Key:")
	resp3, err3 := client.Http().
		SetHeader("X-API-Key", "your-api-key").
		Get("https://httpbin.org/headers")

	if err3 != nil {
		log.Printf("   Error: %v", err3)
	} else {
		fmt.Printf("   API Key Status: %s\n", resp3.Status)
	}
	fmt.Println()
}

// Example 8: Retry and Timeout Examples
func retryAndTimeoutExample() {
	fmt.Println("8. Retry and Timeout Examples:")

	client := cumi.NewClient().
		SetTimeout(5 * time.Second).
		SetRetryCount(3).
		SetRetryInterval(1 * time.Second)

	// This will likely timeout and retry
	resp, err := client.Http().Get("https://httpbin.org/delay/10")

	if err != nil {
		fmt.Printf("   Expected timeout error: %v\n", err)
	} else {
		fmt.Printf("   Unexpected success: %s\n", resp.Status)
	}
	fmt.Println()
}

// Example 9: Context and Cancellation
func contextAndCancellationExample() {
	fmt.Println("9. Context and Cancellation:")

	client := cumi.NewClient()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.Http().
		SetContext(ctx).
		Get("https://httpbin.org/delay/5")

	if err != nil {
		fmt.Printf("   Context timeout error: %v\n", err)
	} else {
		fmt.Printf("   Unexpected success: %s\n", resp.Status)
	}
	fmt.Println()
}

// Example 10: Middleware Examples
func middlewareExample() {
	fmt.Println("10. Middleware Examples:")

	client := cumi.NewClient()

	// Add request middleware for logging
	client.OnBeforeRequest(func(c *cumi.Client, req *cumi.Request) error {
		fmt.Printf("   Making request to API\n")
		return nil
	})

	// Add response middleware for logging
	client.OnAfterResponse(func(c *cumi.Client, resp *cumi.Response) error {
		fmt.Printf("   Response: %s\n", resp.Status)
		return nil
	})

	resp, err := client.Http().Get("https://httpbin.org/get")

	if err != nil {
		log.Printf("   Error: %v", err)
	} else {
		fmt.Printf("   Final Status: %s\n", resp.Status)
	}
	fmt.Println()
}

// Example 11: Query Parameters
func queryParamsExample() {
	fmt.Println("11. Query Parameters:")

	client := cumi.NewClient().
		SetBaseURL("https://jsonplaceholder.typicode.com")

	var posts []Post
	resp, err := client.Http().
		SetQueryParam("userId", "1").
		SetQueryParam("_limit", "3").
		SetSuccessResult(&posts).
		Get("/posts")

	if err != nil {
		log.Printf("   Error: %v", err)
		return
	}

	if resp.IsSuccess() {
		fmt.Printf("   Retrieved %d posts for user 1\n", len(posts))
	}
	fmt.Println()
}

// Example 12: Form Data Example
func formDataExample() {
	fmt.Println("12. Form Data Example:")

	client := cumi.NewClient()
	formData := map[string]string{
		"name":    "John Doe",
		"email":   "john@example.com",
		"message": "Hello from Cumi!",
	}

	resp, err := client.Http().
		SetFormData(formData).
		Post("https://httpbin.org/post")

	if err != nil {
		log.Printf("   Error: %v", err)
		return
	}

	fmt.Printf("   Form submission status: %s\n", resp.Status)
	fmt.Println()
}

// Example 13: Cookie Management
func cookieExample() {
	fmt.Println("13. Cookie Management:")

	client := cumi.NewClient()

	// Create cookies
	sessionCookie := &http.Cookie{
		Name:  "session_id",
		Value: "abc123",
	}
	prefCookie := &http.Cookie{
		Name:  "user_pref",
		Value: "dark_mode",
	}

	// Set cookies
	resp, err := client.Http().
		SetCookie(sessionCookie).
		SetCookie(prefCookie).
		Get("https://httpbin.org/cookies")

	if err != nil {
		log.Printf("   Error: %v", err)
		return
	}

	fmt.Printf("   Cookie request status: %s\n", resp.Status)
	fmt.Println()
}

// Example 14: Debug Mode
func debugModeExample() {
	fmt.Println("14. Debug Mode:")

	client := cumi.NewClient().
		EnableDebug() // Enable debug mode for detailed logging

	resp, err := client.Http().
		SetHeader("Custom-Header", "Debug-Value").
		Get("https://httpbin.org/headers")

	if err != nil {
		log.Printf("   Error: %v", err)
		return
	}

	fmt.Printf("   Debug request status: %s\n", resp.Status)
	fmt.Println()
}
