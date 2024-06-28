package shared

import "strings"

func getAnimeDataString(animeData map[string]interface{}, key string, defaultValue string) string {
	if val, ok := animeData[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
}

func getAnimeDataInt(animeData map[string]interface{}, key string, defaultValue int) int {
	if val, ok := animeData[key]; ok {
		if floatVal, ok := val.(float64); ok {
			return int(floatVal)
		}
	}
	return defaultValue
}

func getAnimeDataBool(animeData map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := animeData[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return defaultValue
}

func getAnimeAiredDate(animeData map[string]interface{}) string {
	if aired, ok := animeData["aired"].(map[string]interface{}); ok {
		if from, ok := aired["from"].(string); ok {
			return strings.Split(from, "T")[0]
		}
	}
	return "N/A"
}

func getAnimeDuration(animeData map[string]interface{}) string {
	if duration, ok := animeData["duration"].(string); ok {
		return duration
	}
	return "N/A"
}
