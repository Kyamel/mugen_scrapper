package animefire

import (
	"errors"
	"fmt"
	"log"
	"mugen_scrapper/src/shared"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Extracts text and hyperlinks from the Anime Fire web page.
func getAnimeFireReleases(url string, extractAmount int, startPage int, extractDub bool, opts ...shared.Option) ([][2]string, error) {
	opt := shared.Options{
		PrintLog: false,
	}
	for _, o := range opts {
		o(&opt)
	}
	var results [][2]string
	page := startPage
	hardLimit := 100
	extractAmount = min(extractAmount, hardLimit) // Aplicar limite máximo
	for len(results) < extractAmount {
		currentURL := strings.Replace(url, "/1", fmt.Sprintf("/%d", page), 1)
		resp, err := http.Get(currentURL)
		if err != nil {
			return nil, fmt.Errorf("erro ao acessar a página: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status de resposta não esperado: %s", resp.Status)
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar o documento HTML: %v", err)
		}
		baseURL := URL_AF_RELEASES // Substitua com o URL base correto
		found := false
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			text := s.Text()
			if strings.Contains(href, baseURL) {
				if !extractDub && strings.Contains(text, "(Dublado)") {
					return
				}
				results = append(results, [2]string{text, href})
				if opt.PrintLog {
					fmt.Printf("Texto: %s, Link: %s\n", text, href)
				}
				if len(results) >= extractAmount {
					found = true
				}
			}
		})
		if found {
			break
		}
		page++
	}
	return results, nil
}

// Helper function to check if any of the substrings are in the string.
func containsAny(str string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.Contains(str, substring) {
			return true
		}
	}
	return false
}

// GetAnimeTitleFromAFURL  Extract the anime name from an Anime Fire url
func getAnimeTitleFromAnimeFireURL(afurl string, opts ...shared.Option) (string, error) {
	opt := shared.Options{
		PrintLog: false,
	}
	for _, o := range opts {
		o(&opt)
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("An error occurred while extracting the anime name: %v", r)
		}
	}()
	// Remove the base URL and split by '/'
	parts := strings.Split(afurl, "/")
	// Get the part containing the anime name and remove '-todos-os-episodios'
	animeTitle := strings.Replace(parts[len(parts)-1], "-todos-os-episodios", "", 1)
	if opt.PrintLog {
		fmt.Printf("Anime name: %s", animeTitle)
	}
	return animeTitle, nil
}

// extractReleasingAnimesFromAnimeFireF collects watch and download links from Anime Fire and anime metadata from MyAnimeList
func extractReleasingAnimesFromAnimeFireF(extractAmount int, startPage int, opts ...shared.Option) ([][]map[string]interface{}, [][]map[string]interface{}, []map[string]interface{}, error) {
	opt := shared.Options{
		PrintLog: false,
	}
	for _, o := range opts {
		o(&opt)
	}
	BAD_ID := 0
	urlAF := URL_AF_RELEASES

	// Step 1: Get titles and hyperlinks from Anime Fire
	animeReleases, err := getAnimeFireReleases(urlAF, extractAmount, startPage, false, opts...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not get anime from anime fire: %v", err)
	}

	var animeIDs []int
	var watchLinks [][]map[string]interface{}
	var downloadLinks [][]map[string]interface{}

	// Step 2: Search each anime in the Jikan API and get the mal_id
	for _, anime := range animeReleases {
		MalID, err := shared.SearchAnimeMalIdOnJikanV4(anime[0], opts...)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error searching anime on jikanv4: %v", err)
		}
		if MalID != BAD_ID {
			animeIDs = append(animeIDs, MalID)
			watchLink, err := getEpisodesWatchLinksFromAnimeFire(anime[1], MalID)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("could not get episodes watch links from anime fire: %v", err)
			}
			watchLinks = append(watchLinks, watchLink)
			AnimeTitle, err := getAnimeTitleFromAnimeFireURL(anime[1], opts...)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("could not get anime title from anime fire url: %v", err)
			}
			downloadLink, err := getEpisodesDownloadLinksFromAnimeFire(AnimeTitle, MalID, opts...)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("could not get download links from anime fire: %v", err)
			}
			downloadLinks = append(downloadLinks, downloadLink)
		}
	}

	// Step 3: Get anime details
	var animesMetadata []map[string]interface{}
	for _, AnimeID := range animeIDs {
		data, err := shared.GetAnimeMetadataFromJikanV4(AnimeID, opts...)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("error getting anime resources from jikanv4: %v", err)
		}
		animesMetadata = append(animesMetadata, data)
	}

	return watchLinks, downloadLinks, animesMetadata, nil
}

// extractCustomAnimeFromAnimeFire extracts a single anime from Anime Fire.
func extractCustomAnimeFromAnimeFire(animeUrlOnAF string, opts ...shared.Option) ([]map[string]interface{}, []map[string]interface{}, map[string]interface{}, error) {
	opt := shared.Options{
		PrintLog: false,
	}
	for _, o := range opts {
		o(&opt)
	}
	BAD_ID := 0
	// Step 1: Extract the anime name from the URL
	urlParts := strings.Split(animeUrlOnAF, "/")
	animeNameRomaji := strings.TrimSuffix(urlParts[len(urlParts)-1], "-todos-os-episodios")
	// Step 2: Search for the anime on Jikan and get the MAL ID
	MalID, err := shared.SearchAnimeMalIdOnJikanV4(animeNameRomaji, opts...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error searching malID on jikanv4: %v", err)
	}
	if MalID == BAD_ID {
		return nil, nil, nil, errors.New("anime not found")
	}
	// Step 3: Get anime data from Jikan
	animeData, err := shared.GetAnimeMetadataFromJikanV4(MalID, opts...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("coud not get anime resources from jikanv4: %v", err)
	}
	// Step 4: Get download links from Anime Fire
	downloadLinks, err := getEpisodesDownloadLinksFromAnimeFire(animeNameRomaji, MalID, opts...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not get download links from anime fire: %v", err)
	}
	// Step 5: Get watch links from Anime Fire
	watchLinks, err := getEpisodesWatchLinksFromAnimeFire(animeUrlOnAF, MalID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not get watch links from anime fire: %v", err)
	}
	return watchLinks, downloadLinks, animeData, nil
}

// ExtractCustomAnimeFromAnimeFireAndSelectedEpisode extracts a single anime from Anime Fire and optionally fetches specific episode data.

func searchAnimeOnAnimeFire(animeName string) ([]string, error) {
	// Convert to lower case
	animeName = strings.ToLower(animeName)
	// Replace spaces with hyphens
	animeName = strings.ReplaceAll(animeName, " ", "-")
	// Remove special characters, only allow letters and numbers
	re := regexp.MustCompile(`[^a-z0-9-]+`)
	animeName = re.ReplaceAllString(animeName, "")
	animeName = animeName + "-todos-os-episodios"

	baseURL := URL_AF_SEARCH
	queryURL := baseURL + url.QueryEscape(animeName)

	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch the page: status code %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	var animeLinks []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, "/animes/") {
			animeLinks = append(animeLinks, href)
		}
	})
	return animeLinks, nil
}

// fetches episode watch links from Anime Fire url
func getEpisodesWatchLinksFromAnimeFire(urlAllEps string, malID int) ([]map[string]interface{}, error) {

	var results []map[string]interface{}
	baseURL := strings.Split(urlAllEps, "-todos-os-episodios")[0]
	for i := 0; i < 1000; i++ { // Assume there won't be more than 1000 episodes
		episodeURL := fmt.Sprintf("%s/%d", baseURL, i)
		resp, err := http.Get(episodeURL)
		if err != nil {
			return nil, fmt.Errorf("error accessing the page: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			if i == 0 {
				continue // Continue checking if there's an error for episode 0
			} else {
				break // Exit the loop when no more episodes are found
			}
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error creating HTML document: %v", err)
		}
		links := doc.Find("a")
		if links.Length() == 0 {
			if i == 0 {
				continue // Continue checking if the page is empty for episode 0
			} else {
				break // Exit the loop when a page without hyperlinks is found
			}
		}
		result := map[string]interface{}{
			"localization":  "pt-br",
			"MalID":         malID,
			"EpisodeNumber": i,
			"WatchLink":     episodeURL,
		}
		results = append(results, result)
	}
	return results, nil
}

// Fetches episode download links from Anime Fire URL
func getEpisodesDownloadLinksFromAnimeFire(AnimeTitle string, MalID int, opts ...shared.Option) ([]map[string]interface{}, error) {
	opt := shared.Options{
		PrintLog: false,
	}
	for _, o := range opts {
		o(&opt)
	}
	urlAfDown := URL_AF_DOWNLOADS
	var results []map[string]interface{}
	baseURL := fmt.Sprintf("%s%s/", urlAfDown, AnimeTitle)
	if opt.PrintLog {
		fmt.Printf("Base URL: %s\n", baseURL)
	}
	// Assume there are no more than 1000 episodes
	for i := 0; i < 1000; i++ {
		episodeURL := fmt.Sprintf("%s%d", baseURL, i)
		resp, err := http.Get(episodeURL)
		if err != nil {
			return nil, fmt.Errorf("error accessing the page: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNotFound {
			if i == 0 {
				continue // If error on episode 0, continue checking
			} else {
				break // Exit loop when no more episodes are found
			}
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error creating HTML document: %v", err)
		}
		links := doc.Find("a")
		filterStrings := []string{
			"https://s2.lightspeedst.net/s2/mp4/",
			"https://s2.lightspeedst.net/s2/mp4_temp/",
			"https://rr",
		}
		var downloadLinks []string
		links.Each(func(index int, item *goquery.Selection) {
			href, exists := item.Attr("href")
			if exists && containsAny(href, filterStrings) {
				downloadLinks = append(downloadLinks, href)
			}
		})
		if len(downloadLinks) == 0 && i != 0 {
			break // Exit loop if no links found after episode 0
		}
		if len(downloadLinks) > 0 {
			downloadLink := map[string]interface{}{
				"Localization":     "portuguese-br",
				"MalID":            MalID,
				"EpisodeNumber":    i,
				"DownloadLinkSD":   "", // Inicializa como uma string vazia
				"DownloadLinkHD":   "",
				"DownloaodLinkFHD": "",
				"DownloadLink4k":   "", // Inicializa como uma string vazia
				"Temp":             false,
			}
			for _, link := range downloadLinks {
				if strings.Contains(link, "(SD)") {
					downloadLink["DownloadLinkSD"] = link
				} else if strings.Contains(link, "(HD)") {
					downloadLink["DownloadLinkHD"] = link
				}
				if strings.Contains(link, "mp4_temp") {
					downloadLink["Temp"] = true
				}
			}
			results = append(results, downloadLink)
			if opt.PrintLog {
				fmt.Printf("Download links for episode %d: %v\n", i, downloadLink)
			}
		}
	}
	return results, nil
}

// GetSpecificEpisodeWatchLinksFromAnimeFire fetches watch links for a specific episode from Anime Fire URL
func getSpecificEpisodeWatchLinksFromAnimeFire(url string, malID int, episode int, opts ...shared.Option) (map[string]interface{}, error) {
	opt := shared.Options{
		PrintLog: false,
	}
	for _, o := range opts {
		o(&opt)
	}
	if url == "" {
		return nil, fmt.Errorf("url is empty")
	}
	if opt.PrintLog {
		fmt.Printf("Base URL: %s\n", url)
	}
	baseURL := strings.Split(url, "-todos-os-episodios")[0]
	episodeURL := fmt.Sprintf("%s/%d", baseURL, episode)
	resp, err := http.Get(episodeURL)
	if err != nil {
		fmt.Printf("Error accessing the page: %v\n", err)
		return nil, fmt.Errorf("error accessing the page: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("Episode not found: %d\n", episode)
		return nil, fmt.Errorf("episode not found")
	}
	if opt.PrintLog {
		fmt.Printf("Watch link for episode %d: %s\n", episode, episodeURL)
	}
	whachLink := map[string]interface{}{
		"localization":  "pt-br",
		"MalID":         malID,
		"EpisodeNumber": episode,
		"WatchLink":     episodeURL,
	}
	return whachLink, nil
}

// GetSpecificEpisodeDownloadLinksFromAnimeFire fetches download links for a specific episode from Anime Fire URL
func getSpecificEpisodeDownloadLinksFromAnimeFire(AnimeTitle string, MalID int, episode int, opts ...shared.Option) (map[string]interface{}, error) {
	opt := shared.Options{
		PrintLog: false,
	}
	for _, o := range opts {
		o(&opt)
	}
	urlAfDown := URL_AF_DOWNLOADS
	baseURL := fmt.Sprintf("%s%s/", urlAfDown, AnimeTitle)
	episodeURL := fmt.Sprintf("%s/%d", baseURL, episode)
	resp, err := http.Get(episodeURL)
	if err != nil {
		return nil, fmt.Errorf("error accessing the page: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("episode not found")
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error creating HTML document: %v", err)
	}
	links := doc.Find("a")
	var downloadLinks []map[string]interface{}
	filterStrings := []string{
		"https://s2.lightspeedst.net/s2/mp4/",
		"https://s2.lightspeedst.net/s2/mp4_temp/",
		"https://rr",
	}
	links.Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if exists && containsAny(href, filterStrings) {
			downloadLink := map[string]interface{}{
				"Localization":     "pt-br",
				"MalID":            MalID,
				"EpisodeNumber":    episode,
				"DownloadLinkSD":   "", // Inicializa como uma string vazia
				"DownloadLinkHD":   "",
				"DownloaodLinkFHD": "",
				"DownloadLink4k":   "", // Inicializa como uma string vazia
				"Temp":             false,
			}
			if strings.Contains(href, "(SD)") {
				downloadLink["DownloadLinkSD"] = href
			} else if strings.Contains(href, "(HD)") {
				downloadLink["DownloadLinkHD"] = href
			}
			if strings.Contains(href, "mp4_temp") {
				downloadLink["Temp"] = true
			}
			downloadLinks = append(downloadLinks, downloadLink)
		}
	})
	return downloadLinks[0], nil
}
