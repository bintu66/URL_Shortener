# URL Shortener

A simple URL shortener built with Go — paste a long URL and get a short one back!

## What Does This Project Do?

You know how some URLs are super long and ugly? Like:

```
https://www.example.com/products/category/electronics/items?id=12345&ref=homepage&utm_source=google
```

This app turns it into something short like:

```
http://localhost:8080/aB3xZ7
```

When someone clicks that short link, they get **redirected** (sent) to the original long URL automatically.

---

## How To Run It

**Step 1:** Make sure you have [Go installed](https://go.dev/dl/).

**Step 2:** Open your terminal, go to the project folder, and run:

```bash
go run cmd/server/main.go
```

**Step 3:** You'll see this message:

```
Server is running on http://localhost:8080
```

**Step 4:** Open your browser and go to **http://localhost:8080** — you'll see the URL shortener form!

---

## How To Use It

###  From the Browser (the easy way)

1. Open **http://localhost:8080** in your browser
2. You'll see a page with a text box
3. Paste any long URL (like `https://www.google.com`) and click **"Shorten it!"**
4. You'll get a short URL — click it, and it will take you to the original site

### 🛠 From Postman / curl (the API way)

Send a POST request with JSON:

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

You'll get back:

```json
{
  "short_url": "http://localhost:8080/aB3xZ7"
}
```

Visit that short URL in your browser, and it redirects you to Google!

---

## Project Structure (What Each File Does)

```
url_shortener/
├── cmd/
│   └── server/
│       └── main.go            ← Entry point. Starts the server and connects everything.
│
├── internal/
│   ├── generator/
│   │   └── code.go            ← Generates random short codes like "aB3xZ7"
│   │
│   ├── handler/
│   │   └── url_handler.go     ← Handles browser/API requests (the "waiter" at a restaurant)
│   │
│   ├── model/
│   │   └── url.go             ← Defines the shape of data (request and response)
│   │
│   ├── service/
│   │   └── url_service.go     ← Business logic — the "brain" that connects handler and storage
│   │
│   └── storage/
│       └── memory.go          ← Stores URLs in memory (like a dictionary/notebook)
│
├── go.mod                     ← Go module file (declares this project's name)
├── README.md                  ← You're reading this!
└── HOW_THE_CODE_WORKS.md      ← Deep explanation of Go syntax used in this project
```

### Why this folder structure?

- **`cmd/`** = where executable programs live. Our server's `main.go` goes here.
- **`internal/`** = private code that only this project can use. We split it into sub-folders by **responsibility**:
  - **handler** = talks to the outside world (HTTP requests)
  - **service** = contains the business logic
  - **storage** = deals with saving/reading data
  - **model** = defines data shapes (structs)
  - **generator** = utility to create random codes

This pattern is called **"separation of concerns"** — each part has ONE job.

---

## How A Request Flows Through The Code

Let's trace what happens when you shorten a URL:

### Shortening a URL (you submit the form)

```
Browser                 main.go              handler              service            storage
  │                       │                     │                    │                   │
  │──POST /shorten───────>│                     │                    │                   │
  │                       │──ShortenURL()──────>│                    │                   │
  │                       │                     │──CreateShortURL()─>│                   │
  │                       │                     │                    │──Generate(6)      │
  │                       │                     │                    │  (makes "aB3xZ7") │
  │                       │                     │                    │──Save("aB3xZ7",──>│
  │                       │                     │                    │   "https://...")   │
  │                       │                     │                    │<──done────────────│
  │                       │                     │<──"aB3xZ7"────────│                   │
  │<──HTML page with──────│                     │                    │                   │
  │   short URL           │                     │                    │                   │
```

### Redirecting (someone clicks the short URL)

```
Browser                 main.go              handler              service            storage
  │                       │                     │                    │                   │
  │──GET /aB3xZ7─────────>│                     │                    │                   │
  │                       │──RedirectURL()─────>│                    │                   │
  │                       │                     │──GetOriginalURL()─>│                   │
  │                       │                     │                    │──Get("aB3xZ7")───>│
  │                       │                     │                    │<──"https://...",──│
  │                       │                     │                    │    true           │
  │                       │                     │<──"https://..."───│                   │
  │<──302 Redirect to─────│                     │                    │                   │
  │   https://...         │                     │                    │                   │
```

---

## API Routes

| Method | Path       | What It Does                                |
|--------|------------|---------------------------------------------|
| GET    | `/`        | Shows the homepage with the URL form        |
| POST   | `/shorten` | Shortens a URL (form or JSON)               |
| GET    | `/{code}`  | Redirects to the original URL               |
| GET    | `/health`  | Simple check — returns "Server is running"  |

---

## Important Notes

- **Data is stored in memory** — when you stop the server, all shortened URLs are gone. This is fine for learning! In a real app, you'd use a database like PostgreSQL or Redis.
- **This runs locally** — the short URLs only work on your computer (`localhost`). To make it available on the internet, you'd need to deploy it to a server.
- **No duplicate checking** — if you shorten the same URL twice, you get two different short codes. A real app might check for duplicates.

---

## Built With

- **Go (Golang)** — the programming language
- **net/http** — Go's built-in HTTP server (no external frameworks needed!)
- **encoding/json** — Go's built-in JSON parser
