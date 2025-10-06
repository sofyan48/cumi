package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/sofyan48/cumi"
)

// Example demonstrating safe and unsafe concurrent usage
// Run this file separately: go run concurrency_safety_test.go

func main() {
	fmt.Println("=== Concurrency Safety Examples ===\n")

	// Example 1: UNSAFE - Race condition
	unsafeExample()

	// Example 2: SAFE - Using Clone()
	safeWithCloneExample()

	// Example 3: SAFE - Configure before concurrent use
	safeConfigBeforeConcurrentExample()

	// Example 4: SAFE - Request-level configuration
	safeRequestLevelExample()

	fmt.Println("\n=== All Examples Completed ===")
}

// ❌ UNSAFE: Modifying shared client concurrently
func unsafeExample() {
	fmt.Println("1. ❌ UNSAFE Example (Race Condition):")
	fmt.Println("   This will cause race condition if run with -race flag")

	client := cumi.NewClient()

	var wg sync.WaitGroup

	// ❌ DON'T DO THIS: Concurrent modification of shared client
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// ⚠️ RACE CONDITION: Multiple goroutines modifying same client
			client.SetCommonHeader("X-Request-ID", fmt.Sprintf("request-%d", id))
			client.SetCommonQueryParam("user_id", fmt.Sprintf("%d", id))

			// This might use wrong headers from other goroutines!
			resp, err := client.Http().Get("https://httpbin.org/get")
			if err != nil {
				fmt.Printf("   Goroutine %d error: %v\n", id, err)
			} else {
				fmt.Printf("   Goroutine %d status: %s\n", id, resp.Status)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println()
}

// ✅ SAFE: Using Clone() for each goroutine
func safeWithCloneExample() {
	fmt.Println("2. ✅ SAFE Example (Using Clone):")

	baseClient := cumi.NewClient().
		SetBaseURL("https://httpbin.org").
		SetTimeout(10 * time.Second)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// ✅ SAFE: Each goroutine gets its own client copy
			client := baseClient.Clone()
			client.SetCommonHeader("X-Request-ID", fmt.Sprintf("request-%d", id))
			client.SetCommonQueryParam("user_id", fmt.Sprintf("%d", id))

			resp, err := client.Http().Get("/get")
			if err != nil {
				fmt.Printf("   Goroutine %d error: %v\n", id, err)
			} else {
				fmt.Printf("   Goroutine %d status: %s\n", id, resp.Status)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println()
}

// ✅ SAFE: Configure client BEFORE concurrent use
func safeConfigBeforeConcurrentExample() {
	fmt.Println("3. ✅ SAFE Example (Configure Before Concurrent Use):")

	// ✅ Configure client ONCE before spawning goroutines
	client := cumi.NewClient().
		SetBaseURL("https://httpbin.org").
		SetTimeout(10 * time.Second).
		SetCommonHeader("User-Agent", "Cumi-Client/1.0").
		SetCommonHeader("Accept", "application/json")

	var wg sync.WaitGroup

	// Now safe to use concurrently (read-only access)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// ✅ SAFE: Only creating new Request, not modifying client
			resp, err := client.Http().
				SetHeader("X-Request-ID", fmt.Sprintf("request-%d", id)).
				SetQueryParam("user_id", fmt.Sprintf("%d", id)).
				Get("/get")

			if err != nil {
				fmt.Printf("   Goroutine %d error: %v\n", id, err)
			} else {
				fmt.Printf("   Goroutine %d status: %s\n", id, resp.Status)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println()
}

// ✅ SAFE: Using request-level configuration
func safeRequestLevelExample() {
	fmt.Println("4. ✅ SAFE Example (Request-Level Configuration):")

	client := cumi.NewClient()

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// ✅ SAFE: Each request has its own configuration
			resp, err := client.Http().
				SetHeader("X-Request-ID", fmt.Sprintf("request-%d", id)).
				SetHeader("User-Agent", fmt.Sprintf("Client-%d", id)).
				SetQueryParam("user_id", fmt.Sprintf("%d", id)).
				SetQueryParam("timestamp", fmt.Sprintf("%d", time.Now().Unix())).
				Get("https://httpbin.org/get")

			if err != nil {
				fmt.Printf("   Goroutine %d error: %v\n", id, err)
			} else {
				fmt.Printf("   Goroutine %d status: %s\n", id, resp.Status)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println()
}

// Best Practices Summary:
//
// ❌ DON'T:
// - Call client.SetCommonHeader(), SetCommonQueryParam(), etc. from multiple goroutines
// - Modify shared client state after spawning goroutines
//
// ✅ DO:
// - Configure client BEFORE spawning goroutines (for shared config)
// - Use client.Clone() to create independent copies for each goroutine
// - Use request-level SetHeader(), SetQueryParam() for per-request config
// - Use sync.Mutex if you really need to modify shared client concurrently (advanced)
