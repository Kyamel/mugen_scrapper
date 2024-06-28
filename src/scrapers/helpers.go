package scrapers

var validPlatforms = map[string]bool{
	"animefire": true,
	// Adicione outras plataformas válidas aqui
	// "anotherPlatform": true,
}

var validDubs = map[string]bool{
	"pt-br": true,
	"none":  true,
	// Adicione outras dubções aqui
	// "anotherDub": true,
}

func isValidPlatform(platform string) bool {
	_, isValid := validPlatforms[platform]
	return isValid
}

func isValidDub(dub string) bool {
	_, isValid := validDubs[dub]
	return isValid
}
