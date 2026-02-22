-- Triangle Travel SQLite Schema
-- Replaces iata_cities.json, distances.json, city_routes.json

-- City codes to airport IATA codes (from iata_cities.json)
CREATE TABLE IF NOT EXISTS iata_cities (
    city_code TEXT NOT NULL,
    airport_code TEXT NOT NULL,
    PRIMARY KEY (city_code, airport_code)
);

CREATE INDEX IF NOT EXISTS idx_iata_cities_city ON iata_cities(city_code);

-- Distances between airports in miles (from distances.json)
CREATE TABLE IF NOT EXISTS distances (
    from_iata TEXT NOT NULL,
    to_iata TEXT NOT NULL,
    distance_miles REAL NOT NULL,
    PRIMARY KEY (from_iata, to_iata)
);

CREATE INDEX IF NOT EXISTS idx_distances_from ON distances(from_iata);
CREATE INDEX IF NOT EXISTS idx_distances_to ON distances(to_iata);

-- City routes by alliance (from city_routes.json)
CREATE TABLE IF NOT EXISTS city_routes (
    city_iata TEXT NOT NULL,
    alliance TEXT NOT NULL,
    route_to TEXT NOT NULL,
    PRIMARY KEY (city_iata, alliance, route_to)
);

CREATE INDEX IF NOT EXISTS idx_city_routes_city ON city_routes(city_iata);
CREATE INDEX IF NOT EXISTS idx_city_routes_alliance ON city_routes(city_iata, alliance);
