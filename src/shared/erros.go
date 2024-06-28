package shared

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrNotFound           = errors.New("not found")
	ErrInvalidRequest     = errors.New("invalid request")
)

// handles the various scraping errors and sends appropriate HTTP responses.
func HandleScrapingError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrServiceUnavailable):
		http.Error(w, "scraping service temporarily unavailable", http.StatusServiceUnavailable)
	case errors.Is(err, ErrNotFound):
		http.Error(w, "not found", http.StatusNotFound)
	case errors.Is(err, ErrInvalidRequest):
		http.Error(w, "invalid request parameters", http.StatusBadRequest)
	default:
		http.Error(w, fmt.Sprintf("platforms scraper failed: %v", err), http.StatusInternalServerError)
	}
}
