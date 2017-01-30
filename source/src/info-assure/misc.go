package main

// Removes any non numeric characters. To be used with the strings.Map function
func RemoveNonNumericChars(r rune) rune {
	switch {
	case r < '0' || r > '9':
		return -1
	}

	return r
}
