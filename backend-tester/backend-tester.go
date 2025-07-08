package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://backend:8080"
	}

	resp, err := http.Get(fmt.Sprintf("%s/list", backendURL))
	if err != nil {
		log.Fatalf("Failed to call backend: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("Response from backend: %s\n", body)
}
