package shared

type Anime struct {
	Title         string
	TitleEnglish  string
	TitleJapanese string
	Year          int
	Season        string
	Synopsis      string
	Source        string
	Episodes      int
	Status        string
	Duration      string
	Rating        string
	Airing        bool
	Type          string
	Background    string
	Aired         string
	MalID         int
	WatchLink     string
}

type Episode struct {
	MalID           int
	Localization    string
	Platform        string
	Number          int
	WatchLink       string
	DownloadLinkSD  string
	DownloadLinkHD  string
	DownloadLinkFHD string
	DownloadLink4K  string
	Temp            bool
}

type Song struct {
	SongName        string
	RelatedAnime    string
	MalID           int
	Type            string
	AniPlaylistLink string
	SpotifyLink     string
	YoutubeLink     string
	AppleMusicLink  string
}

func (a Anime) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":         a.Title,
		"titleEnglish":  a.TitleEnglish,
		"titleJapanese": a.TitleJapanese,
		"year":          a.Year,
		"season":        a.Season,
		"synopsis":      a.Synopsis,
		"source":        a.Source,
		"episodes":      a.Episodes,
		"status":        a.Status,
		"duration":      a.Duration,
		"rating":        a.Rating,
		"airing":        a.Airing,
		"type":          a.Type,
		"background":    a.Background,
		"aired":         a.Aired,
		"malID":         a.MalID,
	}
}

func (e Episode) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"malID":           e.MalID,
		"localization":    e.Localization,
		"platform":        e.Platform,
		"number":          e.Number,
		"watchLink":       e.WatchLink,
		"downloadLinkSD":  e.DownloadLinkSD,
		"downloadLinkHD":  e.DownloadLinkHD,
		"downloadLinkFHD": e.DownloadLinkFHD,
		"downloadLink4K":  e.DownloadLink4K,
		"temp":            e.Temp,
	}
}

func (s Song) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"songName":        s.SongName,
		"relatedAnime":    s.RelatedAnime,
		"malID":           s.MalID,
		"type":            s.Type,
		"aniPlaylistLink": s.AniPlaylistLink,
		"spotifyLink":     s.SpotifyLink,
		"youtubeLink":     s.YoutubeLink,
		"appleMusicLink":  s.AppleMusicLink,
	}
}
