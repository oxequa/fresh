package fresh

import "strings"

// Contain check if a string is inserted into a strings array
func contain(s string, arr []string) bool {
	s = strings.ToLower(s)
	for _, val := range arr {
		if val == s {
			return true
		}
	}
	return false
}
