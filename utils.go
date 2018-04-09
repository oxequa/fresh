package fresh

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	B = 1 << (10 * iota)
	K
	M
	G
	T
)

// Size convert a string like 10K or 5MB in relative int64 number size
func size(s string) (r int64) {
	format := regexp.MustCompile("[A-Z]+")
	num, err := strconv.ParseInt(regexp.MustCompile("[0-9]+").String(), 10, 64)
	if err != nil {
		return
	}
	switch format.String() {
	case "B", "b":
		return num * B
	case "KB", "K", "kb", "k":
		return num * K
	case "MB", "M", "mb", "m":
		return num * M
	case "GB", "G", "gb", "g":
		return num * G
	case "TB", "T", "tb", "t":
		return num * G
	}
	return
}

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
