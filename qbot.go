package main

import (
	"fmt"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/slack"
	"log"
	"os"
	"reflect"
	"strings"
	"github.com/doozr/qbot/usercache"
	"github.com/doozr/qbot/notification"
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

	log.Printf("Connecting to Slack as %s", name)
	slackConn, err := slack.New(name, os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Getting user list")
	users, err := slackConn.GetUserList()
	if err != nil {
		log.Fatal(err)
	}
	userCache := usercache.New(users)

	dumpfile := os.Args[3]
	log.Printf("Attempting to load queue from %s", dumpfile)
	q, err := queue.Load(dumpfile)

	notifications := notification.New(userCache)
	commands := command.New(notifications)

	messageChan := make(chan slack.RtmMessage, 100)
	saveChan := make(chan queue.Queue, 5)
	notifyChan := make(chan Notification, 5)

	go MessageDispatch(name, q, commands, messageChan, saveChan, notifyChan)
	go Save(dumpfile, saveChan)
	go Notify(slackConn, notifyChan)

	log.Println("Ready to receive messages")
	for {
		// read each incoming message
		m, err := slackConn.GetEvent()
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" {
			msg := slack.ConvertEventToMessage(m)
			if strings.HasPrefix(msg.Text, name) || strings.HasPrefix(msg.Text, "<@"+slackConn.Id+">") {
				messageChan <- msg
			}
		}

		if m.Type == "user_change" {
			uc := slack.ConvertEventToUserChange(m)
			userCache.UpdateUserName(uc.User)
		}
	}
}

func splitUser(u string) (username string, reason string) {
	args := strings.SplitN(u, " ", 2)
	username = args[0]
	reason = ""
	if len(args) > 1 {
		reason = args[1]
	}
	return
}

func MessageDispatch(name string, q queue.Queue, commands command.Command,
	messageChan chan slack.RtmMessage, saveChan chan queue.Queue, notifyChan chan Notification) {

	logNotification := func(n string) {
		for _, l := range strings.Split(n, "\n") {
			if l != "" {
				log.Println(l)
			}
		}
	}

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

		oq := q
		n := ""

		switch cmd {
		case "join":
			q, n = commands.Join(q, m.User, rest)
		case "leave":
			q, n = commands.Leave(q, m.User, rest)
		case "done":
			q, n = commands.Done(q, m.User)
		case "yield":
			q, n = commands.Yield(q, m.User)
		case "barge":
			q, n = commands.Barge(q, m.User, rest)
		case "boot":
			id, reason := splitUser(rest)
			q, n = commands.Boot(q, m.User, id, reason)
		case "oust":
			id, reason := splitUser(rest)
			q, n = commands.Oust(q, m.User, id, reason)
		case "list":
			n = commands.List(q)
		case "help":
			n = commands.Help(name)
		}

		if n != "" {
			if !reflect.DeepEqual(oq, q) {
				logNotification(n)
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
			log.Printf("Error saving file to %s: %s", filename, err)
		}
	}
}

func Notify(slackConn *slack.Slack, notifyChan chan Notification) {
	for n := range notifyChan {
		err := slackConn.PostMessage(n.Channel, n.Message)
		if err != nil {
			log.Printf("Error when sending: %s", err)
		}
	}
}
