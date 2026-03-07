# go_server

A Go Fiber application created with create-fiber-app.

## Getting Started

### Prerequisites

- Go 1.24.7 or higher
- Docker (optional)

### Installation

```bash
go mod download
```

### Running the application

```bash
# Run normally
go run main.go

# Run with hot reload (recommended for development)
make dev
```

### Running with Docker

```bash
docker-compose up
```

### Building for production

```bash
docker build -t go_server .
```

## API Endpoints

- `GET /` - Welcome message
- `GET /health` - Health check
- `GET /api/v1/users` - Get users

## Environment Variables

- `PORT` - Server port (default: 3000)

## License

MIT