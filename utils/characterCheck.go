package utils

import "regexp"

func CharacterCheck(str string) bool {
	hasAlphanumeric := regexp.MustCompile(`[a-zA-Z0-9]`).MatchString(str)
	return hasAlphanumeric && len(str) > 0
}
