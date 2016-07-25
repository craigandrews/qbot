package util

import (
	"log"
	"strings"
)

func LogMultiLine(n string) {
	for _, l := range strings.Split(n, "\n") {
		if l != "" {
			log.Println(l)
		}
	}
}

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

func IsPrivateChannel(channel string) bool {
	return strings.HasPrefix(channel, "D")
}
