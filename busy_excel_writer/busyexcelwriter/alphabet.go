package busyexcelwriter

// GetAlphabetByNumber converts an integer column number to its alphabet representation (e.g. 1 -> "A").
func GetAlphabetByNumber(number int) string {
	result := ""
	for number > 0 {
		remainder := (number - 1) % 26
		result = string('A'+rune(remainder)) + result
		number = (number - 1) / 26
	}
	return result
}
