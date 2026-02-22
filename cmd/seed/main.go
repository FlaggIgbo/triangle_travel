// Seed script: go run ./cmd/seed
// Populates SQLite from JSON files. Run from project root.

package main

import (
	"database/sql"
	"encoding/json"
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

	schemaPath := filepath.Join(wd, "db", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatal(err)
	}
	log.Println("Schema created")

	// Seed iata_cities
	iataPath := filepath.Join(wd, "iata_cities.json")
	iataData, err := os.ReadFile(iataPath)
	if err != nil {
		log.Fatal(err)
	}
	var iataCities map[string][]string
	if err := json.Unmarshal(iataData, &iataCities); err != nil {
		log.Fatal(err)
	}
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("INSERT INTO iata_cities (city_code, airport_code) VALUES (?, ?)")
	for city, airports := range iataCities {
		for _, airport := range airports {
			stmt.Exec(city, airport)
		}
	}
	stmt.Close()
	tx.Commit()
	log.Printf("Loaded %d iata_cities entries", len(iataCities))

	// Seed distances
	distPath := filepath.Join(wd, "distances.json")
	distData, err := os.ReadFile(distPath)
	if err != nil {
		log.Fatal(err)
	}
	var distances map[string]map[string]float64
	if err := json.Unmarshal(distData, &distances); err != nil {
		log.Fatal(err)
	}
	tx, _ = db.Begin()
	stmt, _ = tx.Prepare("INSERT INTO distances (from_iata, to_iata, distance_miles) VALUES (?, ?, ?)")
	count := 0
	for from, toMap := range distances {
		for to, dist := range toMap {
			stmt.Exec(from, to, dist)
			count++
		}
	}
	stmt.Close()
	tx.Commit()
	log.Printf("Loaded %d distance entries", count)

	// Seed city_routes
	routesPath := filepath.Join(wd, "city_routes.json")
	routesData, err := os.ReadFile(routesPath)
	if err != nil {
		log.Fatal(err)
	}
	var routes map[string]map[string][]string
	if err := json.Unmarshal(routesData, &routes); err != nil {
		log.Fatal(err)
	}
	tx, _ = db.Begin()
	stmt, _ = tx.Prepare("INSERT INTO city_routes (city_iata, alliance, route_to) VALUES (?, ?, ?)")
	count = 0
	for city, alliances := range routes {
		for alliance, dests := range alliances {
			for _, dest := range dests {
				stmt.Exec(city, alliance, dest)
				count++
			}
		}
	}
	stmt.Close()
	tx.Commit()
	log.Printf("Loaded %d city_routes entries", count)

	fmt.Println("Database seeded successfully at", dbPath)
}
