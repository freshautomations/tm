package utils

// Contains looks for a string in a string slice and returns true if it finds it.
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
