package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/slack"
	"encoding/json"
	"io/ioutil"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: qbot <token> <data file>")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	slackConn, err := slack.New(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dumpfile := os.Args[2]
	q := queue.Queue{}
	if _, err := os.Stat(dumpfile); err == nil {
		dat, err := ioutil.ReadFile(dumpfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		json.Unmarshal(dat, &q)
	}

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
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+slackConn.Id+">") {
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

			user := slackConn.GetUsername(m.User)

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
				args := strings.SplitN(rest, " ", 0)
				q, n = command.Boot(q, user, args[0], args[1])
			case "oust":
				args := strings.SplitN(rest, " ", 0)
				q, n = command.Oust(q, user, args[0], args[1])
			case "list":
				n = command.List(q)
			}
			if n != "" {
				j, err := json.Marshal(q)
				ioutil.WriteFile(dumpfile, j, 0644)
				fmt.Println(string(j))
				fmt.Println(n)
				err = slackConn.PostMessage(m.Channel, n)
				if err != nil {
					fmt.Println("Error when sending: %s", err)
				}
			}
		}
	}
}

