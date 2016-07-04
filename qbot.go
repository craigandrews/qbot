package main

import (
	"fmt"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/slack"
	"os"
	"reflect"
	"strings"
)

type Notification struct {
	Channel string
	Message string
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: qbot <name> <token> <data file>")
		os.Exit(1)
	}

	name := os.Args[1]

	slackConn, err := slack.New(name, os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dumpfile := os.Args[3]
	q, err := queue.Load(dumpfile)

	messageChan := make(chan slack.RtmMessage, 100)
	saveChan := make(chan queue.Queue, 5)
	notifyChan := make(chan Notification, 5)

	go MessageDispatch(q, slackConn, messageChan, saveChan, notifyChan)
	go Save(dumpfile, saveChan)
	go Notify(slackConn, notifyChan)

	fmt.Println("mybot ready, ^C exits")
	for {
		// read each incoming message
		m, err := slackConn.GetMessage()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// see if we're mentioned
		hasPrefix := strings.HasPrefix(m.Text, name) || strings.HasPrefix(m.Text, "<@"+slackConn.Id+">")
		if m.Type == "message" && hasPrefix {
			messageChan <- m
		}
	}
}

func MessageDispatch(q queue.Queue,
	slackConn *slack.Slack,
	messageChan chan slack.RtmMessage,
	saveChan chan queue.Queue,
	notifyChan chan Notification) {
	n := ""
	for m := range messageChan {
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
			n = command.Help(slackConn.Name)
		}
		if n != "" {
			if !reflect.DeepEqual(oq, q) {
				fmt.Println(n)
				saveChan <- q
			}
			notifyChan <- Notification{m.Channel, n}
		}
	}
}

func Save(filename string, saveChan chan queue.Queue) {
	for q := range saveChan {
		err := q.Save(filename)
		if err != nil {
			fmt.Println("Error saving file to %s: %s", filename, err)
		}
	}
}

func Notify(slackConn *slack.Slack, notifyChan chan Notification) {
	for n := range notifyChan {
		err := slackConn.PostMessage(n.Channel, n.Message)
		if err != nil {
			fmt.Println("Error when sending: %s", err)
		}
	}
}