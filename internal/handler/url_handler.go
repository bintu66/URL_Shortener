package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"url_shortener/internal/model"
	"url_shortener/internal/service"
)

// URLHandler handles all incoming HTTP requests related to URLs.
// It talks to the URLService to do the actual work (shorten, lookup).
type URLHandler struct {
	service *service.URLService
}

// NewURLHandler creates a new handler connected to the given service.
func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) HomePage(w http.ResponseWriter, r *http.Request) {

	html := `<!DOCTYPE html>
<!DOCTYPE html>
<html>
<head>
    <title>URL Shortener</title>
    <style>
        body {
            margin: 0;
            height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            background: #000;
            color: white;
            font-family: Arial, sans-serif;
        }

        form {
            text-align: center;
        }

        input {
            width: 300px;
            padding: 12px;
            border: 1px solid #444;
            border-radius: 5px;
            background: #111;
            color: white;
        }

        button {
            padding: 12px 20px;
            margin-left: 8px;
            border: none;
            border-radius: 5px;
            background: #fff;
            color: #000;
            cursor: pointer;
        }
    </style>
</head>
<body>
    <form action="/shorten" method="POST">
        <h1>URL Shortener</h1>
        <input
            type="text"
            name="url"
            placeholder="Paste URL here..."
            required
        >
        <button type="submit">Shorten</button>
    </form>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var longURL string

	// Check what type of request this is by looking at the Content-Type header
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		// --- JSON request (from Postman / curl / API tools) ---
		var req model.URLRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		longURL = req.URL
	} else {
		longURL = r.FormValue("url")
	}

	if longURL == "" {
		http.Error(w, "Please provide a URL", http.StatusBadRequest)
		return
	}

	code := h.service.CreateShortURL(longURL)
	shortURL := "http://localhost:8080/" + code

	//  If the request came from the browser, shows a  HTML result page 
	if !strings.Contains(contentType, "application/json") {
		resultHTML := fmt.Sprintf(`<!DOCTYPE html>
        <html>
        <head>
            <title>URL Shortened</title>
            <style>
                body {
                    margin: 0;
                    height: 100vh;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    background: #000;
                    color: #fff;
                    font-family: Arial, sans-serif;
                }

                .container {
                    text-align: center;
                    max-width: 700px;
                    padding: 20px;
                }

                .short-url a {
                    color: #4da6ff;
                    font-size: 24px;
                    text-decoration: none;
                    word-break: break-all;
                }

                .original {
                    color: #888;
                    margin-top: 15px;
                    word-break: break-all;
                }

                .back {
                    display: inline-block;
                    margin-top: 25px;
                    color: #fff;
                    text-decoration: none;
                    border: 1px solid #444;
                    padding: 10px 16px;
                    border-radius: 5px;
                }

                .back:hover {
                    background: #111;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>URL Shortened</h1>

                <p>Your short URL:</p>

                <p class="short-url">
                    <a href="%s">%s</a>
                </p>

                <p class="original">
                    Original: %s
                </p>

                <a class="back" href="/">Shorten Another URL</a>
            </div>
        </body>
</html>`, shortURL, shortURL, longURL)

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, resultHTML)
		return
	}

	//  If the request was JSON (API), send back JSON 
	resp := model.URLResponse{
		ShortURL: shortURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}


func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	longURL, found := h.service.GetOriginalURL(code)

	if !found {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}