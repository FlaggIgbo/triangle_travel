// Seed script: go run ./cmd/seed
// Populates SQLite from db/schema.sql and db/seed_data.sql. Run from project root.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func main() {
	wd, _ := os.Getwd()
	dbPath := filepath.Join(wd, "db", "data.sqlite3")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatal(err)
	}
	os.Remove(dbPath) // Start fresh

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Schema
	schemaPath := filepath.Join(wd, "db", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatal(err)
	}
	log.Println("Schema created")

	// Auth schema
	authSchemaPath := filepath.Join(wd, "db", "schema_auth.sql")
	authSchema, err := os.ReadFile(authSchemaPath)
	if err == nil {
		if _, err := db.Exec(string(authSchema)); err != nil {
			log.Fatal(err)
		}
		log.Println("Auth schema created")
	}

	// Seed data (embedded in SQL)
	seedPath := filepath.Join(wd, "db", "seed_data.sql")
	seedData, err := os.ReadFile(seedPath)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(string(seedData)); err != nil {
		log.Fatal(err)
	}
	log.Println("Seed data loaded")

	fmt.Println("Database seeded successfully at", dbPath)
}
