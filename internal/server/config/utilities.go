package config

import "unicode"

// isNumeric - функция проверки строки на наличие в ней только цифр
//
// isNumeric - function of checking a string for the presence of only numbers in it
func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
