package config

import "unicode"

// isNumeric is a function that checks whether a string is a number.
func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
