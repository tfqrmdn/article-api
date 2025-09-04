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

# Generate Postman collection
postman-docs:
	@echo "Postman collection available at: postman_collection.json"
	@echo "Import this file into Postman to test the API"

# Complete local development setup
setup-local: setup-local-db postman-docs
	@echo "Local development setup complete!"
	@echo "Run 'make run-local-no-redis' to start the application"
	@echo "Or 'make run-local' if you have Redis running locally"

# Test API endpoints
test-api:
	@echo "Testing basic API endpoints..."
	@echo "1. Testing GET /articles..."
	@curl -s -X GET 'http://localhost:8080/articles' | head -c 200
	@echo "\n2. Testing POST /articles..."
	@curl -s -X POST 'http://localhost:8080/articles' -H 'Content-Type: application/json' -d '{"author_id":"author-1","title":"Test Article","body":"Test body"}' | head -c 200
	@echo "\nAPI tests completed. Check api_tests.json for comprehensive test cases."

# Show API test examples
show-tests:
	@echo "API test cases available in api_tests.json"
	@echo "Example commands:"
	@echo "  Basic list: curl -X GET 'http://localhost:8080/articles'"
	@echo "  With search: curl -X GET 'http://localhost:8080/articles?search=Go'"
	@echo "  With author: curl -X GET 'http://localhost:8080/articles?author=John'"
	@echo "  Create article: curl -X POST 'http://localhost:8080/articles' -H 'Content-Type: application/json' -d '{\"author_id\":\"author-1\",\"title\":\"Test\",\"body\":\"Test body\"}'"

