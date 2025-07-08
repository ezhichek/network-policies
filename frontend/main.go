package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var backendURL string

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(fmt.Sprintf("%s/list", backendURL))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error connecting to backend: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var names []string
	if err := json.NewDecoder(resp.Body).Decode(&names); err != nil {
		http.Error(w, fmt.Sprintf("Decode error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

func main() {
	backendURL = os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://backend:8080"
	}

	http.HandleFunc("/", handler)

	log.Println("Frontend is listening on :80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("Failed to start frontend: %v", err)
	}
}
