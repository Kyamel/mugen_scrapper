package routes

import (
	"encoding/json"
	"mugen_scrapper/src/scrapers"
	"mugen_scrapper/src/shared"

	"net/http"
)

func GetAnimes(w http.ResponseWriter, r *http.Request) {
	anime := r.URL.Query().Get("anime")
	platform := r.URL.Query().Get("platform")

	animes, err := scrapers.ScrapeAnimes(anime, platform)
	if err != nil {
		shared.HandleScrapingError(w, err)
		return
	}

	response := make([]map[string]interface{}, len(animes))
	for i, a := range animes {
		response[i] = a.ToMap()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
