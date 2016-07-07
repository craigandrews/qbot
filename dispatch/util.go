package dispatch

import (
	"log"
	"strings"
)

func splitUser(u string) (username string, reason string) {
	args := strings.SplitN(u, " ", 2)
	username = args[0]
	reason = ""
	if len(args) > 1 {
		reason = args[1]
	}
	return
}

func logResponse(n string) {
	for _, l := range strings.Split(n, "\n") {
		if l != "" {
			log.Println(l)
		}
	}
}

func splitCommand(m string) (cmd string, args string) {
	parts := strings.SplitN(m, " ", 3)

	if len(parts) < 2 {
		return
	}
	cmd = parts[1]

	args = ""
	if len(parts) > 2 {
		args = parts[2]
	}

	return
}
