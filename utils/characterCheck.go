package utils

import "regexp"

func AlphanumericCheck(email string) bool {
	return regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`).MatchString(email)
}

func NumericCheck(password string) bool {
	return regexp.MustCompile("^[0-9_]*$").MatchString(password)
}
