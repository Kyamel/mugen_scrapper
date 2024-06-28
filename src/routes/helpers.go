package routes

import (
	"regexp"
	"strings"
)

func replaceSpacesWithHyphens(input string) string {
	// Remove múltiplos espaços
	reducedSpaces := reduceSpaces(input)
	// Substitui espaços por hífens
	result := strings.ReplaceAll(reducedSpaces, " ", "-")
	return result
}

func reduceSpaces(input string) string {
	// Remove múltiplos espaços usando expressão regular
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(input, " ")
}
