package main

import (
	"fmt"
	"net/http"
	"os"
)

// Test backend
// To run print PORT=9001 go run ./main.go
// Or $env:PORT=9001
// go run ./main.go on Windows
func main() {
	// Get port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000" // default port
	}

	// Simple test handler
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("Hello from test backend on port %s\n", port)
		_, err := rw.Write([]byte(msg))
		if err != nil {
			fmt.Println("Failed to write message:", err)
		}
	})

	fmt.Printf("Test backend listening on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Failed to start test backend:", err)
	}
}
