# Article API Service

A REST API service for managing articles and authors, built with Go, PostgreSQL, and Redis.

## Features

- **List Articles**: GET `/articles` - Retrieve articles with search, filtering, and pagination
- **Create Article**: POST `/articles` - Create a new article
- **Redis Caching**: 10-minute cache for article listings (with fallback to mock cache)
- **Search & Filtering**: Search by title/body content and filter by author name
- **Pagination**: Configurable page size and page navigation
- **Configurable Server**: Customizable host, port, and timeout settings
- **Graceful Shutdown**: Proper server shutdown with signal handling
- **Database**: PostgreSQL with raw SQL queries (no ORM)
- **Docker**: Containerized application with Docker Compose
- **Local Development**: Run without Docker using local PostgreSQL and Redis
- **Testing**: Comprehensive unit tests with coverage reports

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

### Option 1: Docker (Recommended)

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

### Option 2: Local Development

1. **Setup local PostgreSQL and Redis**:
   ```bash
   # Create database and user
   sudo -u postgres psql -c "CREATE DATABASE article_db;"
   sudo -u postgres psql -c "CREATE USER default WITH PASSWORD 'secret';"
   sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE article_db TO default;"
   ```

2. **Setup local database**:
   ```bash
   make setup-local-db
   ```

3. **Run the application locally**:
   ```bash
   # With Redis (if available)
   make run-local
   
   # Without Redis (uses mock cache)
   make run-local-no-redis
   ```

The API will be available at `http://localhost:8080`

## API Endpoints

### List Articles
```bash
GET /articles?search=go&author=John&page=1&limit=10
```

**Query Parameters:**
- `search` (optional): Search term for title and body content
- `author` (optional): Filter by author name
- `page` (optional): Page number for pagination (default: 1)
- `limit` (optional): Number of items per page (default: 10)

**Response Headers:**
- `X-Total-Count`: Total number of articles
- `X-Page`: Current page number
- `X-Limit`: Items per page
- `X-Total-Pages`: Total number of pages

**Response:**
```json
[
  {
    "id": "article-1",
    "author_id": "author-1",
    "title": "Getting Started with Go",
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
├── docker.env.example              # Environment variables example
├── Makefile                        # Build and development commands
├── postman_collection.json         # Postman API collection
├── scripts/
│   ├── migrations/                 # Database migration files
│   │   ├── 001_create_authors_table.sql
│   │   ├── 002_create_articles_table.sql
│   │   └── 003_create_migrations_table.sql
│   ├── seeders/                    # Database seeder files
│   │   ├── 001_seed_authors.sql
│   │   ├── 002_seed_articles.sql
│   │   └── 003_seed_comprehensive_articles.sql
│   ├── migrate/                    # Migration runner
│   │   └── migrate.go
│   └── seed/                       # Seeder runner
│       └── seed.go
├── tests/                          # Test coverage reports (git ignored)
└── internal/
    ├── config/
    │   ├── config.go               # Configuration management
    │   └── config_test.go          # Config tests
    ├── database/
    │   └── connection.go           # Database connection logic
    ├── models/
    │   └── article.go              # Data models
    ├── repository/
    │   ├── interfaces.go           # Repository interfaces
    │   ├── article_repository.go   # Database operations
    │   └── article_repository_test.go # Repository tests
    ├── handlers/
    │   ├── article_handler.go      # HTTP request handlers
    │   └── article_handler_test.go # Handler tests
    ├── cache/
    │   ├── interface.go            # Cache interface
    │   ├── redis.go                # Redis implementation
    │   └── mock.go                 # Mock cache for testing
    └── migration/
        └── migrate.go              # Migration runner for app startup
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

**Generate test coverage:**
```bash
go test -v -coverprofile=tests/coverage.out ./...
go tool cover -html=tests/coverage.out -o tests/coverage.html
```

### Environment Variables

The application supports the following environment variables:

**Server Configuration:**
- `HTTP_SERVER_HOST` - Server host address (default: 0.0.0.0)
- `HTTP_SERVER_PORT` - Server port (default: 8080)
- `HTTP_SERVER_READ_TIMEOUT` - Read timeout duration (default: 30s)
- `HTTP_SERVER_WRITE_TIMEOUT` - Write timeout duration (default: 30s)
- `HTTP_SERVER_IDLE_TIMEOUT` - Idle timeout duration (default: 120s)

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

### Configuration File

You can use environment variables to configure the application. Docker Compose supports environment variable substitution with default values.

**Option 1: Create a .env file (recommended for local development):**
```bash
cp docker.env.example .env
# Edit .env with your specific configuration
```

**Option 2: Set environment variables directly:**
```bash
export HTTP_SERVER_PORT=3000
export DB_PASSWORD=my_secure_password
docker-compose up -d
```

**Option 3: Use environment variables inline:**
```bash
HTTP_SERVER_PORT=3000 DB_PASSWORD=my_secure_password docker-compose up -d
```

### Available Make Commands

**Development:**
- `make build` - Build the application
- `make test` - Run unit tests
- `make run` - Run the application with default configuration
- `make run-dev` - Run the application with custom development settings
- `make deps` - Install dependencies

**Docker:**
- `make docker-up` - Start Docker services
- `make docker-down` - Stop Docker services
- `make test-db` - Run tests with database
- `make clean` - Clean up build artifacts and Docker volumes

**Database:**
- `make migrate` - Run database migrations
- `make seed` - Run database seeders
- `make setup-db` - Run migrations and seeders

**Local Development:**
- `make run-local` - Run locally with Redis
- `make run-local-no-redis` - Run locally without Redis (uses mock cache)
- `make setup-local-db` - Setup local database (migrate + seed)
- `make seed-local` - Seed local database only
- `make setup-local` - Complete local setup (database + docs)

**API Testing:**
- `make test-api` - Test basic API endpoints
- `make show-tests` - Show example API test commands
- `make postman-docs` - Generate Postman collection

## Sample Data

The database is initialized with comprehensive sample data:

**Authors:** 20 authors with diverse names
**Articles:** 100+ articles with realistic content covering various topics including:
- Programming languages (Go, Python, JavaScript, etc.)
- Database technologies (PostgreSQL, Redis, MongoDB, etc.)
- Web development (APIs, REST, GraphQL, etc.)
- DevOps and deployment
- Software architecture and design patterns

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Successful GET request
- `201 Created` - Successful POST request
- `400 Bad Request` - Invalid request data or missing fields
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server-side errors

## Caching

The API uses Redis for caching with the following features:

- **Article List**: Cached for 10 minutes
- **Cache Invalidation**: Automatically invalidated when new articles are created
- **Fallback**: If Redis is unavailable, the application uses a mock cache service
- **Local Development**: Can run without Redis using mock cache for development

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
# List all articles
curl http://localhost:8080/articles

# List articles with search and pagination
curl "http://localhost:8080/articles?search=go&page=1&limit=5"

# List articles by specific author
curl "http://localhost:8080/articles?author=John"

# Create a new article
curl -X POST http://localhost:8080/articles \
  -H "Content-Type: application/json" \
  -d '{"author_id":"author-1","title":"New Article","body":"Content here"}'

# Test API endpoints
make test-api

# Show more test examples
make show-tests
```

## API Testing

The project includes comprehensive API testing resources:

- **Postman Collection**: `postman_collection.json` - Import into Postman for GUI testing
- **Test Commands**: Use `make show-tests` to see example curl commands
- **Automated Testing**: Use `make test-api` for basic endpoint testing

## License

This project is open source and available under the MIT License.
