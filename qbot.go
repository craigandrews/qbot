package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/slack"
	"reflect"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: qbot <name> <token> <data file>")
		os.Exit(1)
	}

	name := os.Args[1]

	// start a websocket-based Real Time API session
	slackConn, err := slack.New(os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dumpfile := os.Args[3]
	q, err := queue.Load(dumpfile)

	fmt.Println("mybot ready, ^C exits")
	for {
		// read each incoming message
		m, err := slackConn.GetMessage()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}


		// see if we're mentioned
		n := ""
		fmt.Println(m.Text)
		if m.Type == "message" && strings.HasPrefix(m.Text, name) {
			// if so try to parse if
			parts := strings.SplitN(m.Text, " ", 3)

			if len(parts) < 2 {
				continue
			}
			cmd := parts[1]

			rest := ""
			if len(parts) > 2 {
				rest = parts[2]
			}

			user := slackConn.GetUsername(m.User)
			oq := q

			fmt.Printf("User: %s Command: %s Args: %v\n", user, cmd, rest)
			switch cmd {
			case "join":
				q, n = command.Join(q, user, rest)
			case "leave":
				q, n = command.Leave(q, user, rest)
			case "done":
				q, n = command.Done(q, user)
			case "yield":
				q, n = command.Yield(q, user)
			case "barge":
				q, n = command.Barge(q, user, rest)
			case "boot":
				args := strings.SplitN(rest, " ", 2)
				if len(args) == 2 {
					q, n = command.Boot(q, user, args[0], args[1])
				} else {
					q, n = command.Boot(q, user, args[0], "")
				}
			case "oust":
				args := strings.SplitN(rest, " ", 2)
				if len(args) == 2 {
					q, n = command.Oust(q, user, args[0], args[1])
				} else {
					q, n = command.Oust(q, user, args[0], "")
				}
			case "list":
				n = command.List(q)
			case "help":
				n = command.Help(name)
			}
			if n != "" {
				if !reflect.DeepEqual(oq, q) {
					err = q.Save(dumpfile)
					if err != nil {
						fmt.Println("Error saving file to %s: %s", dumpfile, err)
					}
				}
				err = slackConn.PostMessage(m.Channel, n)
				if err != nil {
					fmt.Println("Error when sending: %s", err)
				}
			}
		}
	}
}

