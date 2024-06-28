package routes

import (
	"encoding/json"
	"mugen_scrapper/src/scrapers"
	"mugen_scrapper/src/shared"
	"net/http"
)

func GetSongs(w http.ResponseWriter, r *http.Request) {
	anime := r.URL.Query().Get("anime")
	typ := r.URL.Query().Get("type")

	// Função de scrapping (a ser implementada)
	songs, err := scrapers.ScrapeSongs(anime, typ)
	if err != nil {
		shared.HandleScrapingError(w, err)
		return
	}

	response := make([]map[string]interface{}, len(songs))
	for i, s := range songs {
		response[i] = s.ToMap()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
