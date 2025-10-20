# Events API

A simple RESTful API for managing events built with Go, PostgreSQL, and Chi router.

## Features

- Create, read, and list events
- JSON validation and error handling
- PostgreSQL integration with pgx
- Chi router with middleware
- Docker Compose for local development

## Quick Start

1. **Start PostgreSQL with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

2. **Set environment variable:**
   ```bash
   export DATABASE_URL="postgres://events_user:events_password@localhost:5432/events_db?sslmode=disable"
   ```

3. **Run the application:**
   ```bash
   go run cmd/server/main.go
   ```

## API Endpoints

### Create Event
```bash
POST /events
Content-Type: application/json

{
  "title": "My Event",
  "description": "Event description",
  "start_date": "2024-01-01T10:00:00Z",
  "end_date": "2024-01-01T12:00:00Z"
}
```

### Get Event
```bash
GET /events/{id}
```

### List Events
```bash
GET /events?page=1&limit=10
```

## Database Schema

The events table includes:
- `id`: UUID primary key
- `title`: Event title (required)
- `description`: Event description (optional)
- `start_date`: Event start time (required)
- `end_date`: Event end time (required)
- `created_at`: Record creation timestamp

## Development

The project structure:
```
cmd/server/
├── main.go                 # Application entry point
├── internal/
│   ├── db/
│   │   └── db.go          # Database connection
│   ├── handlers/
│   │   └── event_handler.go # HTTP handlers
│   ├── models/
│   │   └── event.go       # Data models
│   └── repository/
│       └── event_repo.go  # Database operations
└── schema.sql             # Database schema
```
