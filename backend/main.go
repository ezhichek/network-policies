package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Person struct {
	Name string `json:"name"`
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT name FROM people")
	if err != nil {
		log.Printf("DB query failed: %v", err)
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("Scan failed: %v", err)
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		names = append(names, name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(names)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	var p Person
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("INSERT INTO people(name) VALUES($1)", p.Name)
	if err != nil {
		http.Error(w, "DB insert error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Inserted")
}

func main() {
	connStr := os.Getenv("DB_CONN")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS people (
		id SERIAL PRIMARY KEY,
		name TEXT
	)`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/list", listHandler)

	log.Println("Backend is listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
