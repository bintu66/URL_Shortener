# URL Shortener

A minimal URL shortener written in Go. No dependencies outside the standard library.

For example:
Turns this:

```
https://www.example.com/products/category/electronics/items?id=12345&ref=homepage&utm_source=google
```

into this:

```
http://localhost:8080/aB3xZ7
```

## Requirements

- Go 1.21+

## Running

```bash
go run cmd/server/main.go
```

Server starts on `http://localhost:8080`.

## Usage

### Web UI

Go to `http://localhost:8080`, paste a URL, submit. You get back a short link that redirects to the original.

### API

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

```json
{
  "short_url": "http://localhost:8080/aB3xZ7"
}
```

Hitting the short URL issues a 302 redirect to the original.

## API Routes

| Method | Path       | Description                          |
|--------|------------|---------------------------------------|
| GET    | `/`        | HTML form                             |
| POST   | `/shorten` | Create a short URL (form or JSON)     |
| GET    | `/{code}`  | Redirect to the original URL          |
| GET    | `/health`  | Health check                          |

## Project Layout

```
url_shortener/
├── cmd/server/main.go          # entry point, wiring
├── internal/
│   ├── generator/code.go       # random short-code generation
│   ├── handler/url_handler.go  # HTTP layer
│   ├── model/url.go            # request/response types
│   ├── service/url_service.go  # business logic
│   └── storage/memory.go       # in-memory store
├── go.mod
└── HOW_THE_CODE_WORKS.md       # walkthrough of the Go syntax used here
```

Standard handler → service → storage split. Handlers parse requests and write responses, the service owns the logic, storage is swappable behind an interface.

## Limitations

- **In-memory storage.** Everything is lost on restart. Swap `internal/storage` for a Postgres/Redis-backed implementation if you need persistence.
- **No dedupe.** Shortening the same URL twice produces two different codes.
- **Not collision-checked against a real keyspace long-term** — fine for a learning project, not for scale.
- **localhost only.** Deploy behind a real domain + reverse proxy to use it outside your machine.

## Stack

Go, `net/http`, `encoding/json`. That's it.