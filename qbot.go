package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/command"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: mybot slack-bot-token\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	ws, id := slackConnect(os.Args[1])
	fmt.Println("mybot ready, ^C exits")
	q := queue.Queue{}

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}


		// see if we're mentioned
		n := ""
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.SplitN(m.Text, " ", 2)
			fmt.Println(parts)
			if len(parts) < 2 {
				continue
			}
			cmd := parts[1]
			fmt.Println(cmd)
			rest := ""
			if len(parts) > 2 {
				rest = parts[2]
			}
			switch cmd {
			case "join":
				q, n = command.Join(q, m.User, rest)
			case "leave":
				q, n = command.Leave(q, m.User, rest)
			case "done":
				q, n = command.Done(q, m.User)
			case "yield":
				q, n = command.Yield(q, m.User)
			case "barge":
				q, n = command.Barge(q, m.User, rest)
			case "boot":
				args := strings.SplitN(rest, " ", 0)
				q, n = command.Boot(q, m.User, args[0], args[1])
			case "oust":
				args := strings.SplitN(rest, " ", 0)
				q, n = command.Oust(q, m.User, args[0], args[1])
			case "list":
				n = command.List(q)
			}
			if n != "" {
				fmt.Println(n)
				go func(m Message) {
					m.Text = n
					postMessage(ws, m)
				}(m)
			}
		}
	}
}

