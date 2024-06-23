package utils

import "errors"

var (
	timezoneShortcodeFullMap = map[string]string{
		"GMT": "Europe/London",
		"EET": "Europe/Helsinki",
	}
	timezoneFullShortcodeMap = ReverseMap(timezoneShortcodeFullMap)
)

func GetFullLocation(shortcode string) (string, error) {
	if full, exists := timezoneShortcodeFullMap[shortcode]; exists {
		return full, nil
	}
	return "", errors.New("unsupported timezone shortcode")
}

func GetShortLocation(full string) (string, error) {
	if short, exists := timezoneFullShortcodeMap[full]; exists {
		return short, nil
	}
	return "", errors.New("unsupported timezone full location")
}
