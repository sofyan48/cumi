# Cumi Http Requester

A simple HTTP client library for Go
## Installation

```bash
go get github.com/sofyan48/requester
```

## Quick Start
### Basic GET Request

```go
package main

import (
    "fmt"
    "log"
    "github.com/sofyan48/cumi"
)

func main() {
    client := cumi.NewClient()
    resp, err := client.Http().Get("https://httpbin.org/get")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Status:", resp.Status)
    fmt.Println("Response:", resp.String())
}
```

### SetSuccessResult - Auto JSON Parsing

Compared to standard net/http which requires manual steps:

```go
// Traditional way with net/http
resp, err := http.Get("https://api.example.com/users/1")
if err != nil {
    return err
}
defer resp.Body.Close()

body, err := io.ReadAll(resp.Body)
if err != nil {
    return err
}

var user User
err = json.Unmarshal(body, &user)
if err != nil {
    return err
}
```

**With requester SetSuccessResult, it becomes much simpler:**

```go
// With requester - auto parsing!
var user User
client := cumi.NewClient()
resp, err := client.Http().
    SetSuccessResult(&user).
    Get("https://api.example.com/users/1")

if err != nil {
    return err
}

if resp.IsSuccess() {
    // user is already populated with data!
    fmt.Printf("User: %+v\n", user)
}
```

## Features

- **SetSuccessResult/SetErrorResult** - Automatic response parsing
- **Zero external dependencies** - Only uses Go standard library  
- **Built-in retry mechanism** with configurable backoff
- **Authentication support** - Basic Auth, Bearer Token, API Key
- **Request/Response middleware** for logging and preprocessing
- **File upload/download** with progress callbacks
- **Debug mode** for request/response logging
- **TLS configuration** for custom certificates
- **Context support** for timeout and cancellation
- **Form data and JSON** request bodies
- **Query parameters and path parameters**
- **Cookie management**
- **Custom headers** per request or globally
- **Tracing support** with OpenTelemetry integration

## Examples

### GET Request with Auto JSON Parsing

```go
package main

import (
    "fmt"
    "log"
    "github.com/sofyan48/requester"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

func main() {
    client := cumi.NewClient()
    var user User
    resp, err := client.Http().
        SetSuccessResult(&user).
        Get("https://jsonplaceholder.typicode.com/users/1")
    
    if err != nil {
        log.Fatal(err)
    }
    
    if resp.IsSuccess() {
        fmt.Printf("User: %+v\n", user)
    }
}
```

### POST with JSON

```go
func main() {
    client := cumi.NewClient()
    user := User{Name: "John", Email: "john@example.com"}
    
    resp, err := client.Http().
        SetBodyJSON(user).
        Post("https://httpbin.org/post")
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Status:", resp.Status)
    fmt.Println("Response:", resp.String())
}
```

### POST with Auto Result Binding

```go
func main() {
    client := cumi.NewClient()
    user := User{Name: "John", Email: "john@example.com"}
    
    // Auto result binding with SetSuccessResult
    var result map[string]interface{}
    resp, err := client.Http().
        SetBodyJSON(user).
        SetSuccessResult(&result).
        Post("https://httpbin.org/post")
    
    if err != nil {
        panic(err)
    }
    
    if resp.IsSuccess() {
        fmt.Printf("Received JSON: %+v\n", result["json"])
    }
}
```

### Error Handling with SetErrorResult

```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func main() {
    client := cumi.NewClient()
    var user User
    var apiError APIError
    
    resp, err := client.Http().
        SetSuccessResult(&user).
        SetErrorResult(&apiError).
        Get("https://api.example.com/users/999")
    
    if err != nil {
        log.Fatal(err)
    }
    
    if resp.IsSuccess() {
        fmt.Printf("User: %+v\n", user)
    } else {
        fmt.Printf("API Error: %+v\n", apiError)
    }
}
```

### Client Configuration

```go
client := cumi.NewClient().
    SetBaseURL("https://api.example.com").
    SetTimeout(30*time.Second).
    SetCommonHeader("Authorization", "Bearer your-token").
    SetRetryCount(3)

var users []User
resp, err := client.Http().
    SetSuccessResult(&users).
    Get("/users")
```
