package db

import (
	"database/sql"
	"log"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps sql.DB for triangle travel data
type DB struct {
	*sql.DB
}

// New opens the SQLite database at dataDir/db/data.sqlite3
func New(dataDir string) (*DB, error) {
	path := filepath.Join(dataDir, "db", "data.sqlite3")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return &DB{db}, nil
}

// GetAirportsForCity returns airport IATA codes for a city code
func (d *DB) GetAirportsForCity(cityCode string) ([]string, error) {
	rows, err := d.Query("SELECT airport_code FROM iata_cities WHERE city_code = ? ORDER BY airport_code", cityCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var airports []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			log.Println(err)
			continue
		}
		airports = append(airports, code)
	}
	if len(airports) == 0 {
		airports = []string{cityCode}
	}
	return airports, nil
}

// GetDistancesFrom returns map of to_iata -> distance_miles for an airport
func (d *DB) GetDistancesFrom(fromIata string) (map[string]float64, error) {
	rows, err := d.Query("SELECT to_iata, distance_miles FROM distances WHERE from_iata = ?", fromIata)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]float64)
	for rows.Next() {
		var to string
		var dist float64
		if err := rows.Scan(&to, &dist); err != nil {
			log.Println(err)
			continue
		}
		result[to] = dist
	}
	return result, nil
}

// GetRoutesFrom returns list of destination IATA codes for city+alliance
func (d *DB) GetRoutesFrom(cityIata, alliance string) ([]string, error) {
	rows, err := d.Query("SELECT route_to FROM city_routes WHERE city_iata = ? AND alliance = ? ORDER BY route_to", cityIata, alliance)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var routes []string
	for rows.Next() {
		var to string
		if err := rows.Scan(&to); err != nil {
			log.Println(err)
			continue
		}
		routes = append(routes, to)
	}
	return routes, nil
}

// GetRoutesFromWithFallback checks city first, then each airport in city
func (d *DB) GetRoutesFromWithFallback(cityIata, alliance string) ([]string, error) {
	routes, err := d.GetRoutesFrom(cityIata, alliance)
	if err != nil || len(routes) > 0 {
		return routes, err
	}
	airports, err := d.GetAirportsForCity(cityIata)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool)
	for _, apt := range airports {
		r, _ := d.GetRoutesFrom(apt, alliance)
		for _, dest := range r {
			if !seen[dest] {
				seen[dest] = true
				routes = append(routes, dest)
			}
		}
	}
	return routes, nil
}

// HasDirectRoute checks if there's a route from fromIata to any of toIatas for the alliance
func (d *DB) HasDirectRoute(fromIata, alliance string, toIatas map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for k := range toIatas {
		result[k] = false
	}
	routes, err := d.GetRoutesFrom(fromIata, alliance)
	if err != nil {
		return result
	}
	for _, r := range routes {
		if toIatas[r] {
			result[r] = true
		}
	}
	airports, _ := d.GetAirportsForCity(fromIata)
	for _, apt := range airports {
		if apt == fromIata {
			continue
		}
		r, _ := d.GetRoutesFrom(apt, alliance)
		for _, dest := range r {
			if toIatas[dest] {
				result[dest] = true
			}
		}
	}
	return result
}
