package util

import "strings"

// StringPop takes returns the first space-delimited token and the rest of the string
//
// The resulting strings are stripped of whitespace
func StringPop(m string) (first string, rest string) {
	parts := strings.SplitN(m, " ", 2)

	if len(parts) < 1 {
		return
	}
	first = parts[0]

	rest = ""
	if len(parts) > 1 {
		rest = strings.Trim(parts[1], " \t\r\n")
	}

	return
}
