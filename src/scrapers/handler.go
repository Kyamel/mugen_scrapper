package scrapers

import (
	"errors"
	"fmt"
	"mugen_scraper/src/scrapers/animefire"
	"mugen_scraper/src/shared"
	"strconv"
)

func ScrapeAnimes(animeName string, platform string) ([]shared.Anime, error) {
	return []shared.Anime{}, nil
}

func ScrapeEpisodes(anime string, episode string, platform string, dub string) ([]shared.Episode, error) {
	// Check if the required "anime" parameter is provided
	if anime == "" {
		return nil, errors.New("anime is a required parameter")
	}
	if episode == "" {
		return nil, errors.New("episode is a required parameter")
	}
	// Convert episode to integer
	episodeInt, err := strconv.Atoi(episode)
	if err != nil {
		return nil, fmt.Errorf("episode must be an integer: %w", shared.ErrInvalidRequest)
	}
	// Set default dub if not provided
	if dub == "" {
		dub = "none"
	}
	// Validate the 'platform' if provided
	if platform != "" && !isValidPlatform(platform) {
		return nil, fmt.Errorf("platform not supported: %w", shared.ErrInvalidRequest)
	}
	if !isValidDub(dub) {
		return nil, fmt.Errorf("dub not supported: %w", shared.ErrInvalidRequest)
	}

	var episodes []shared.Episode

	if platform != "" {
		// Fetch episode(s) for the specified platform
		switch platform {
		case "animefire":
			ep, err := animefire.HandlerEpisode(anime, episodeInt, dub)
			if ep == nil {
				return nil, fmt.Errorf("anime or episodes not found: %w", shared.ErrNotFound)
			}
			if err != nil {
				return nil, fmt.Errorf("error handling episode: %w", err)
			}
			return append(episodes, *ep), nil
		default:
			return nil, fmt.Errorf("platform not supported: %w", shared.ErrInvalidRequest)
		}
	}

	// Fetch episodes from all platforms if no platform is specified
	availablePlatforms := []string{"animefire"} // Add more platforms as needed
	for _, plat := range availablePlatforms {
		switch plat {
		case "animefire":
			ep, err := animefire.HandlerEpisode(anime, episodeInt, dub)
			if err != nil {
				return nil, fmt.Errorf("error handling episode: %w", err)
			}
			if ep != nil {
				episodes = append(episodes, *ep)
			}
		}
	}

	if len(episodes) == 0 {
		return nil, fmt.Errorf("anime or episodes not found: %w", shared.ErrNotFound)
	}
	return episodes, nil
}

func ScrapePlatform(platform string, anime string, dub string) (*shared.Anime, []shared.Episode, error) {
	// Check if the required "anime" parameter is provided
	if anime == "" {
		return nil, nil, errors.New("anime is a required parameter")
	}
	// Set default platform if not provided
	if platform == "" {
		platform = "animefire"
	}
	// Set default platform if not provided
	if dub == "" {
		dub = "none"
	}
	// Validate the 'platform'
	if !isValidPlatform(platform) {
		return nil, nil, fmt.Errorf("platform not supported: %w", shared.ErrInvalidRequest)
	}
	if !isValidDub(dub) {
		return nil, nil, fmt.Errorf("dub not supported: %w", shared.ErrInvalidRequest)
	}

	// Switch case for different platforms
	switch platform {
	case "animefire":
		animeData, episodesData, err := animefire.HandlerPlatform(anime, dub)
		if animeData == nil || episodesData == nil {
			return nil, nil, fmt.Errorf("anime or episodes not found: %w", shared.ErrNotFound)
		}
		if err != nil {
			return nil, nil, fmt.Errorf("error handling platform: %w", err)
		}
		return animeData, episodesData, nil
	default:
		return nil, nil, fmt.Errorf("platform handler not implemented: %w", shared.ErrInvalidRequest)
	}
}

func ScrapeSongs(animeName string, typ string) ([]shared.Song, error) {
	return []shared.Song{}, nil
}
