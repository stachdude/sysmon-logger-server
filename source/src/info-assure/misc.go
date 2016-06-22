package main

// Removes any non numeric characters. To be used with the strings.Map function
func RemoveNonNumericChars(r rune) rune {
	switch {
	case r < '0' || r > '9':
		return -1
	}

	return r
}

// Removes any non numeric characters. To be used with the strings.Map function
func RemoveNonAlphaNumericChars(r rune) rune {
	switch {
	case r >= 'A' && r <= 'Z':
		return r
	case r >= 'a' && r <= 'z':
		return r
	default:
		return -1
	}
}