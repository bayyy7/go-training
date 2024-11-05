package utils

import "regexp"

func MagicSum(a, b int) int {
	return a + b
}

func MagicSub(a, b int) int {
	return a / b
}

func MagicPow(n int) int {
	return n * n
}

func MagicOdd(n int) bool {
	return n%2 == 0
}

func MagicGrade(n int) string {
	switch n {
	case 0:
		return "Zero"
	case 1:
		return "Bad"
	case 2:
		return "Ok"
	case 3:
		return "Nice"
	case 4:
		return "Awesome"
	case 5:
		return "Exceptinal"
	default:
		return "Unknown"
	}
}

func MagicName(n int) []string {
	var name []string

	for i := 0; i < n; i++ {
		name = append(name, "Bayu")
	}
	return name
}

func MagicTria(n int) int {
	var x int
	for i := 0; i < n; i++ {
		x += i
	}

	return x
}

func MagicChange(n *int) {
	*n = *n * 2
}

func AlphanumericCheck(email string) bool {
	return regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`).MatchString(email)
}

func NumericCheck(password string) bool {
	return regexp.MustCompile("^[0-9_]*$").MatchString(password)
}
