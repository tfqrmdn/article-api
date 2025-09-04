.PHONY: build test run docker-up docker-down clean

# Build the application
build:
	go build -o bin/article-api main.go

# Run tests
test:
	go test -v ./...

# Run the application
run:
	go run main.go

# Run with custom configuration
run-dev:
	SERVER_PORT=3000 SERVER_READ_TIMEOUT=60s go run main.go

# Run locally with local database and Redis
run-local:
	DB_HOST=localhost DB_USERNAME=default DB_PASSWORD=secret DB_DATABASE=article_db REDIS_HOST=localhost HTTP_SERVER_PORT=8080 go run main.go

# Run locally without Redis (uses mock cache)
run-local-no-redis:
	DB_HOST=localhost DB_USERNAME=default DB_PASSWORD=secret DB_DATABASE=article_db HTTP_SERVER_PORT=8080 go run main.go

# Start Docker services
docker-up:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# Clean up
clean:
	rm -rf bin/
	docker-compose down -v

# Run tests with database
test-db: docker-up
	sleep 5
	go test -v ./...
	docker-compose down

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run migrations
migrate:
	go run scripts/migrate/migrate.go

# Run seeders
seed:
	go run scripts/seed/seed.go

# Run migrations and seeders
setup-db: migrate seed

# Setup local database (migrate + seed)
setup-local-db:
	DB_HOST=localhost DB_USERNAME=default DB_PASSWORD=secret DB_DATABASE=article_db go run scripts/migrate/migrate.go
	DB_HOST=localhost DB_USERNAME=default DB_PASSWORD=secret DB_DATABASE=article_db go run scripts/seed/seed.go

# Seed local database only
seed-local:
	DB_HOST=localhost DB_USERNAME=default DB_PASSWORD=secret DB_DATABASE=article_db go run scripts/seed/seed.go

