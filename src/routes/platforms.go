package routes

import (
	"encoding/json"
	"mugen_scrapper/src/scrapers"
	"mugen_scrapper/src/shared"
	"net/http"
)

// GetPlatforms handles the /get-platforms endpoint.
func GetPlatforms(w http.ResponseWriter, r *http.Request) {
	// Get the query parameters
	platform := r.URL.Query().Get("platform")
	anime := r.URL.Query().Get("anime")
	dub := r.URL.Query().Get("dub")

	anime = replaceSpacesWithHyphens(anime)
	// Scrape the platform for the given anime
	animeData, episodesData, err := scrapers.ScrapePlatform(platform, anime, dub)
	if err != nil {
		shared.HandleScrapingError(w, err)
		return
	}
	// Create the response map
	response := map[string]interface{}{
		"anime":    animeData,
		"episodes": episodesData,
	}
	// Set the content type and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
