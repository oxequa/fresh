package fresh

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
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

func PrintBanner() {
	color.Set(color.FgHiGreen)
	println(banner)
	color.Unset()
}

// Print the list of routes
func PrintRouter(r *router) {
	println()
	var tree func(routes []*route, parentPath string) error
	tree = func(routes []*route, parentPath string) error {
		for _, route := range routes {
			separator := ""
			if strings.HasSuffix(parentPath, "/") == false {
				separator = "/"
			}
			currentPath := parentPath + separator + route.path
			for _, handler := range route.handlers {
				print(time.Now().Format("(2006-01-02 03:04:05)---"))
				print("[")
				color.Set(color.FgHiGreen)
				print(handler.method)
				color.Unset()
				print("]")
				for i := len(handler.method); i < 8; i++ {
					print("-")
				}
				print(">")
				println(currentPath)
			}
			tree(route.children, currentPath)
		}
		return nil
	}
	tree([]*route{r.route}, "")
}
