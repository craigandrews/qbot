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

	messageChan := make(chan slack.RtmMessage, 100)
	saveChan := make(chan queue.Queue, 5)
	notifyChan := make(chan Notification, 5)

	go MessageDispatch(name, q, userCache, messageChan, saveChan, notifyChan)
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

func MessageDispatch(name string, q queue.Queue, userCache *usercache.UserCache,
	messageChan chan slack.RtmMessage, saveChan chan queue.Queue, notifyChan chan Notification) {

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

		user := userCache.GetUserName(m.User)
		oq := q
		n := ""

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
			username, reason := splitUser(rest)
			q, n = command.Boot(q, user, username, reason)
		case "oust":
			username, reason := splitUser(rest)
			q, n = command.Oust(q, user, username, reason)
		case "list":
			n = command.List(q)
		case "help":
			n = command.Help(name)
		}

		if n != "" {
			if !reflect.DeepEqual(oq, q) {
				for _, l := range strings.Split(n, "\n") {
					if l != "" {
						log.Println(l)
					}
				}
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
