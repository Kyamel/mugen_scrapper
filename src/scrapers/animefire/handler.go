package animefire

import (
	"fmt"
	"mugen_scrapper/src/shared"
	"strings"
)

func HandlerEpisode(anime string, episode int, dub string) (*shared.Episode, error) {
	// Search for the anime on Jikan and get the MAL ID
	MalID, err := shared.SearchAnimeMalIdOnJikanV4(anime)
	if err != nil {
		return nil, fmt.Errorf("error searching malID on jikanv4: %v", err)
	}
	// Get anime data from Jikan
	animeData, err := shared.GetAnimeMetadataFromJikanV4(MalID)
	if err != nil {
		return nil, fmt.Errorf("could not get anime resources from jikanv4: %v", err)
	}
	// Search for the anime on Anime Fire
	links, err := searchAnimeOnAnimeFire(animeData["Title"].(string))
	if err != nil {
		return nil, fmt.Errorf("error searching anime: %v", err)
	}
	if len(links) == 0 {
		return nil, fmt.Errorf("anime not found")
	}
	// Select the appropriate link based on the 'dub' parameter
	var selectedLink string
	for _, link := range links {
		if (dub == "pt-br" && strings.Contains(link, "dublado")) || (dub != "pt-br" && !strings.Contains(link, "dublado")) {
			selectedLink = link
			break
		}
	}
	if selectedLink == "" {
		return nil, fmt.Errorf("no suitable link found")
	}
	// Extract the anime name from the URL
	urlParts := strings.Split(selectedLink, "/")
	if len(urlParts) == 0 {
		return nil, fmt.Errorf("invalid URL format")
	}
	animeNameOnAF := strings.TrimSuffix(urlParts[len(urlParts)-1], "-todos-os-episodios")
	// Get download links from Anime Fire
	downloadLink, err := getSpecificEpisodeDownloadLinksFromAnimeFire(animeNameOnAF, MalID, episode)
	if err != nil {
		return nil, fmt.Errorf("could not get download links from anime fire: %v", err)
	}
	// Get watch links from Anime Fire
	watchLink, err := getSpecificEpisodeWatchLinksFromAnimeFire(selectedLink, MalID, episode)
	if err != nil {
		return nil, fmt.Errorf("could not get watch links from anime fire: %v", err)
	}
	episodeData := shared.Episode{
		MalID:           MalID,
		Localization:    "pt-br",
		Platform:        "AnimeFire",
		Number:          getInt(watchLink, "EpisodeNumber"),
		WatchLink:       getString(watchLink, "WatchLink"),
		DownloadLinkSD:  getString(downloadLink, "DownloadLinkSD"),
		DownloadLinkHD:  getString(downloadLink, "DownloadLinkHD"),
		DownloadLinkFHD: getString(downloadLink, "DownloadLinkFHD"),
		DownloadLink4K:  getString(downloadLink, "DownloadLink4K"),
		Temp:            getBool(downloadLink, "Temp"),
	}
	return &episodeData, nil

}

func HandlerPlatform(animeName string, dub string) (*shared.Anime, []shared.Episode, error) {
	// Search for the anime on Jikan and get the MAL ID
	MalID, err := shared.SearchAnimeMalIdOnJikanV4(animeName)
	if err != nil {
		return nil, nil, fmt.Errorf("error searching malID on jikanv4: %v", err)
	}
	// Get anime data from Jikan
	animeData, err := shared.GetAnimeMetadataFromJikanV4(MalID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get anime resources from jikanv4: %v", err)
	}
	// Search for the anime on Anime Fire
	links, err := searchAnimeOnAnimeFire(animeData["Title"].(string))
	if err != nil {
		return nil, nil, fmt.Errorf("error searching anime: %v", err)
	}
	if len(links) == 0 {
		return nil, nil, fmt.Errorf("anime not found")
	}
	// Select the appropriate link based on the 'dub' parameter
	var selectedLink string
	for _, link := range links {
		if (dub == "pt-br" && strings.Contains(link, "dublado")) || (dub != "pt-br" && !strings.Contains(link, "dublado")) {
			selectedLink = link
			break
		}
	}
	if selectedLink == "" {
		return nil, nil, fmt.Errorf("no suitable link found")
	}
	// Extract the anime name from the URL
	urlParts := strings.Split(selectedLink, "/")
	if len(urlParts) == 0 {
		return nil, nil, fmt.Errorf("invalid URL format")
	}
	animeNameOnAF := strings.TrimSuffix(urlParts[len(urlParts)-1], "-todos-os-episodios")

	// Get download links from Anime Fire
	downloadLinks, err := getEpisodesDownloadLinksFromAnimeFire(animeNameOnAF, MalID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get download links from anime fire: %v", err)
	}
	// Get watch links from Anime Fire
	watchLinks, err := getEpisodesWatchLinksFromAnimeFire(selectedLink, MalID)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get watch links from anime fire: %v", err)
	}

	// Create the anime
	anime := shared.Anime{
		Title:         getString(animeData, "Title"),
		TitleEnglish:  getString(animeData, "TitleEnglish"),
		TitleJapanese: getString(animeData, "TitleJapanese"),
		Year:          getInt(animeData, "Year"),
		Season:        getString(animeData, "Season"),
		Synopsis:      getString(animeData, "Synopsis"),
		Source:        getString(animeData, "Source"),
		Episodes:      getInt(animeData, "Episodes"),
		Status:        getString(animeData, "Status"),
		Duration:      getString(animeData, "Duration"),
		Rating:        getString(animeData, "Rating"),
		Airing:        getBool(animeData, "Airing"),
		Type:          getString(animeData, "Type"),
		Background:    getString(animeData, "Background"),
		Aired:         getString(animeData, "Aired"),
		MalID:         MalID,
		WatchLink:     selectedLink,
	}

	// Create the list of episodes
	var episodes []shared.Episode
	for i := range watchLinks {
		if i >= len(downloadLinks) {
			break // Ensure we don't access out-of-bounds indices
		}
		episodes = append(episodes, shared.Episode{
			MalID:           MalID,
			Localization:    "pt-br",
			Platform:        "AnimeFire",
			Number:          getInt(watchLinks[i], "EpisodeNumber"),
			WatchLink:       getString(watchLinks[i], "WatchLink"),
			DownloadLinkSD:  getString(downloadLinks[i], "DownloadLinkSD"),
			DownloadLinkHD:  getString(downloadLinks[i], "DownloadLinkHD"),
			DownloadLinkFHD: getString(downloadLinks[i], "DownloadLinkFHD"),
			DownloadLink4K:  getString(downloadLinks[i], "DownloadLink4K"),
			Temp:            getBool(downloadLinks[i], "Temp"),
		})
	}
	return &anime, episodes, nil
}
