.PHONY: dev dev-backend dev-frontend docker-up docker-down build test seed

# Run everything with Docker
docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

# Local development
dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev

# Build
build-backend:
	cd backend && go build -o server ./cmd/server

build-frontend:
	cd frontend && npm run build

# Test
test-backend:
	cd backend && go test ./...

# Seed
seed:
	bash scripts/seed_tenants.sh
