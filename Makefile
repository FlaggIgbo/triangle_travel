.PHONY: seed build dev server

# Seed SQLite from JSON files (run from project root)
seed:
	go run ./cmd/seed

# Build frontend + Go binary
build: build-frontend build-server

build-frontend:
	npm run build

build-server:
	go build -o triangle_travel .

# Dev: run Vite (proxies /api to backend) and Go server separately
dev:
	npm run dev

server:
	go run .

# Full production run (requires: make build first)
serve: build
	./triangle_travel
