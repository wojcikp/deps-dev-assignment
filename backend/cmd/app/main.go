package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	cwd, _ := os.Getwd()

	p := path.Join(cwd, "..", "..", "..", "data", "app.db")
	log.Println("db path:", p)

	db, err := sql.Open("sqlite3", p)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insertTestData(db)

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Content-Type", "application/json")
		response := Response{Message: "Hello from Go Backend!"}
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Server started at :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func insertTestData(db *sql.DB) {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS dependencies (
        name TEXT NOT NULL
    );`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	testData := []string{"ABC", "DEF", "QWERTY", "TESTTEST"}
	for _, name := range testData {
		_, err := db.Exec("INSERT INTO dependencies (name) VALUES (?)", name)
		if err != nil {
			log.Printf("db insert err '%s': %v", name, err)
		}
	}
	log.Println("Test data inserted successfully.")
}
