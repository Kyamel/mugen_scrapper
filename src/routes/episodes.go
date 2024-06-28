package routes

import (
	"encoding/json"
	"mugen_scrapper/src/scrapers"
	"mugen_scrapper/src/shared"
	"net/http"
)

func GetEpisodes(w http.ResponseWriter, r *http.Request) {
	anime := r.URL.Query().Get("anime")
	episode := r.URL.Query().Get("episode")
	platform := r.URL.Query().Get("platform")
	dub := r.URL.Query().Get("dub")

	anime = replaceSpacesWithHyphens(anime)

	// Função de scrapping (a ser implementada)
	episodes, err := scrapers.ScrapeEpisodes(anime, episode, platform, dub)
	if err != nil {
		shared.HandleScrapingError(w, err)
		return
	}
	// Create the response map
	response := map[string]interface{}{
		"episodes": episodes,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
