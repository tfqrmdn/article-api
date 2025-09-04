# Article API Service

A secure REST API service for managing articles and authors, built with Go, PostgreSQL, and Redis.

## Features

- **List Articles**: GET `/articles` - Retrieve all articles with author information
- **Create Article**: POST `/articles` - Create a new article
- **API Key Authentication**: X-API-Key header required for all requests
- **Redis Caching**: 10-minute cache for article listings
- **Request Tracking**: Unique request ID for each request
- **Configurable Server**: Customizable host, port, and timeout settings
- **Graceful Shutdown**: Proper server shutdown with signal handling
- **Database**: PostgreSQL with raw SQL queries (no ORM)
- **Docker**: Containerized application with Docker Compose
- **Testing**: Comprehensive unit tests

## Database Schema

The service uses two main tables:

### Authors Table
- `id` (TEXT, Primary Key)
- `name` (TEXT)

### Articles Table
- `id` (TEXT, Primary Key)
- `author_id` (TEXT, Foreign Key to authors.id)
- `title` (TEXT)
- `body` (TEXT)
- `created_at` (TIMESTAMP)

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Redis (included in Docker setup)
- Make (optional, for using Makefile commands)

## Quick Start

1. **Clone and navigate to the project**:
   ```bash
   cd article-api
   ```

2. **Start the services with Docker Compose**:
   ```bash
   make docker-up
   # or
   docker-compose up -d
   ```

3. **Install Go dependencies**:
   ```bash
   make deps
   # or
   go mod tidy
   ```

4. **Run the application**:
   ```bash
   make run
   # or
   go run main.go
   ```

The API will be available at `http://localhost:8080`

## Authentication

All API requests require an `X-API-Key` header with a valid API key. The default API key for development is `default-api-key-123`.

## API Endpoints

### List Articles
```bash
GET /articles
X-API-Key: default-api-key-123
```

**Response:**
```json
[
  {
    "id": "article-1",
    "author_id": "author-1",
    "title": "Getting Started with Go",
    "body": "This is a comprehensive guide to getting started with Go programming language.",
    "created_at": "2024-01-01T12:00:00Z",
    "author": {
      "id": "author-1",
      "name": "John Doe"
    }
  }
]
```

### Create Article
```bash
POST /articles
Content-Type: application/json
X-API-Key: default-api-key-123

{
  "author_id": "author-1",
  "title": "My New Article",
  "body": "This is the content of my new article."
}
```

**Response:**
```json
{
  "id": "article-1234567890",
  "author_id": "author-1",
  "title": "My New Article",
  "body": "This is the content of my new article.",
  "created_at": "2024-01-01T12:00:00Z",
  "author": {
    "id": "author-1",
    "name": "John Doe"
  }
}
```

## Development

### Project Structure
```
article-api/
├── main.go                          # Application entry point
├── go.mod                           # Go module definition
├── Dockerfile                       # Docker image configuration
├── docker-compose.yml              # Docker services configuration
├── Makefile                        # Build and development commands
├── scripts/
│   └── init.sql                    # Database initialization script
└── internal/
    ├── database/
    │   └── connection.go           # Database connection logic
    ├── models/
    │   └── article.go              # Data models
    ├── repository/
    │   ├── article_repository.go   # Database operations
    │   └── article_repository_test.go # Repository tests
    └── handlers/
        ├── article_handler.go      # HTTP request handlers
        └── article_handler_test.go # Handler tests
```

### Running Tests

**Unit tests only (with mocks):**
```bash
make test
# or
go test -v ./...
```

**Integration tests with database:**
```bash
make test-db
```

### Environment Variables

The application supports the following environment variables:

**Server Configuration:**
- `SERVER_HOST` - Server host address (default: 0.0.0.0)
- `SERVER_PORT` - Server port (default: 8080)
- `SERVER_READ_TIMEOUT` - Read timeout duration (default: 30s)
- `SERVER_WRITE_TIMEOUT` - Write timeout duration (default: 30s)
- `SERVER_IDLE_TIMEOUT` - Idle timeout duration (default: 120s)

**Database Configuration:**
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_DATABASE` - Database name (default: article_db)
- `DB_USERNAME` - Database user (default: article_user)
- `DB_PASSWORD` - Database password (default: article_password)

**Redis Configuration:**
- `REDIS_HOST` - Redis host (default: localhost)
- `REDIS_PORT` - Redis port (default: 6379)
- `REDIS_PASSWORD` - Redis password (default: empty)
- `REDIS_DB` - Redis database number (default: 0)

**Security:**
- `API_KEY` - API key for authentication (default: default-api-key-123)

### Configuration File

You can use environment variables to configure the application. Docker Compose supports environment variable substitution with default values.

**Option 1: Create a .env file (recommended for local development):**
```bash
cp docker.env.example .env
# Edit .env with your specific configuration
```

**Option 2: Set environment variables directly:**
```bash
export SERVER_PORT=3000
export DB_PASSWORD=my_secure_password
export API_KEY=my_custom_api_key
docker-compose up -d
```

**Option 3: Use environment variables inline:**
```bash
SERVER_PORT=3000 DB_PASSWORD=my_secure_password docker-compose up -d
```

### Available Make Commands

- `make build` - Build the application
- `make test` - Run unit tests
- `make run` - Run the application with default configuration
- `make run-dev` - Run the application with custom development settings
- `make docker-up` - Start Docker services
- `make docker-down` - Stop Docker services
- `make test-db` - Run tests with database
- `make deps` - Install dependencies
- `make migrate` - Run database migrations
- `make seed` - Run database seeders
- `make setup-db` - Run migrations and seeders
- `make clean` - Clean up build artifacts and Docker volumes

## Sample Data

The database is initialized with sample data:

**Authors:**
- John Doe (author-1)
- Jane Smith (author-2)

**Articles:**
- "Getting Started with Go" by John Doe
- "Database Design Best Practices" by Jane Smith

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful GET request
- `201 Created` - Successful POST request
- `400 Bad Request` - Invalid request data or missing fields
- `401 Unauthorized` - Missing or invalid API key
- `500 Internal Server Error` - Server-side errors

## Request Tracking

Each request gets a unique ID for tracking (returned in `X-Request-ID` header)

## Caching

The API uses Redis for caching with the following features:

- **Article List**: Cached for 10 minutes
- **Cache Invalidation**: Automatically invalidated when new articles are created
- **Fallback**: If Redis is unavailable, requests fall back to database queries

## Technology Stack

- **Language**: Go 1.21
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Database Driver**: lib/pq
- **Redis Client**: go-redis/v9
- **UUID Generation**: google/uuid
- **Containerization**: Docker & Docker Compose
- **Testing**: Go's built-in testing package

## Example Usage

```bash
# List articles
curl -H "X-API-Key: default-api-key-123" http://localhost:8080/articles

# Create a new article
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -H "X-API-Key: default-api-key-123" \
  -d '{"author_id":"author-1","title":"New Article","body":"Content here"}'
```

## License

This project is open source and available under the MIT License.
