package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Opções que podem ser passadas para a função
type Options struct {
	PrintLog bool
}

// Função para aplicar opções personalizadas
type Option func(*Options)

// Função para configurar o log de impressão
func WithPrintLog(enable bool) Option {
	return func(opt *Options) {
		opt.PrintLog = enable
	}
}

// SearchAnimeMalIdOnJikanV4 pesquisa o ID de um anime na API Jikan v4.
func SearchAnimeMalIdOnJikanV4(animeTitle string, opts ...Option) (int, error) {
	// Passo 1: Realizar uma pesquisa pelo nome do anime
	opt := Options{
		PrintLog: false,
	}
	for _, o := range opts {
		o(&opt)
	}
	BAD_ID := 0
	searchURL := fmt.Sprintf("%s%s", URL_JIKAN_SEARCH, strings.ReplaceAll(animeTitle, " ", "%20"))
	resp, err := http.Get(searchURL)
	if err != nil {
		return BAD_ID, fmt.Errorf("http error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return BAD_ID, fmt.Errorf("bad response: %d %s", resp.StatusCode, resp.Status)

	}
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return BAD_ID, fmt.Errorf("could not decode JSON: %v", err)
	}
	// Verificar se a pesquisa retornou resultados
	if animeData, ok := data["data"].([]interface{}); ok && len(animeData) > 0 {
		if firstAnime, ok := animeData[0].(map[string]interface{}); ok {
			if malID, ok := firstAnime["mal_id"].(float64); ok {
				if opt.PrintLog {
					fmt.Printf("Anime ID: %d", int(malID))
				}
				return int(malID), nil
			}
		}
	} else {
		return BAD_ID, errors.New("error: Anime not found")
	}
	return BAD_ID, errors.New("error: Anime not found")
}

// GetAnimeMetadataFromJikanV4 fetches anime data from MAL through the Jikan v4 API
// GetAnimeMetadataFromJikanV4 fetches anime data from MAL through the Jikan v4 API
func GetAnimeMetadataFromJikanV4(AnimeID int, opts ...Option) (map[string]interface{}, error) {
	opt := Options{
		PrintLog: false,
	}
	for _, o := range opts {
		o(&opt)
	}
	BAD_ID := 0
	if AnimeID == BAD_ID {
		return nil, errors.New("bad AnimeID")
	}
	searchURL := fmt.Sprintf("%s%d", URL_JIKAN_SEARCH_BY_MALID, AnimeID)
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %d %s", resp.StatusCode, resp.Status)
	}
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		if opt.PrintLog {
			fmt.Printf("Error decoding JSON response: %v\n", err)
		}
		return nil, fmt.Errorf("could not decode JSON: %v", err)
	}
	if animeData, ok := data["data"].(map[string]interface{}); ok {
		animeDict := map[string]interface{}{
			"MalID":         AnimeID,
			"Title":         getAnimeDataString(animeData, "title", "N/A"),
			"TitleEnglish":  getAnimeDataString(animeData, "title_english", "N/A"),
			"TitleJapanese": getAnimeDataString(animeData, "title_japanese", "N/A"),
			"Type":          getAnimeDataString(animeData, "type", "N/A"),
			"Episodes":      getAnimeDataInt(animeData, "episodes", 0),
			"Status":        getAnimeDataString(animeData, "status", "N/A"),
			"Airing":        getAnimeDataBool(animeData, "airing", false),
			"Aired":         getAnimeAiredDate(animeData),
			"Rating":        getAnimeDataString(animeData, "rating", "N/A"),
			"Duration":      getAnimeDuration(animeData),
			"Season":        getAnimeDataString(animeData, "season", "N/A"),
			"Year":          getAnimeDataInt(animeData, "year", 0),
			"Synopsis":      getAnimeDataString(animeData, "synopsis", "N/A"),
			"Source":        getAnimeDataString(animeData, "source", "N/A"),
			"Background":    getAnimeDataString(animeData, "background", "N/A"),
		}
		if opt.PrintLog {
			fmt.Println(animeDict)
		}
		return animeDict, nil
	} else {
		return nil, fmt.Errorf("expected a dictionary for anime_data, got: %T", animeData)
	}
}
