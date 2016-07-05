package dispatch

import (
	"log"
	"reflect"
	"strings"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/slack"
	"github.com/doozr/qbot/usercache"
)

type Notification struct {
	Channel string
	Message string
}

type MessageChan chan slack.RtmMessage
type SaveChan chan queue.Queue
type NotifyChan chan Notification
type UserChan chan slack.UserInfo

// Message handles executing user commands and passing on the results
func Message(name string, q queue.Queue, commands command.Command,
	messageChan MessageChan, saveChan SaveChan, notifyChan NotifyChan) {

	splitUser := func(u string) (username string, reason string) {
		args := strings.SplitN(u, " ", 2)
		username = args[0]
		reason = ""
		if len(args) > 1 {
			reason = args[1]
		}
		return
	}

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
			q, n = commands.Oust(q, m.User, rest)
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

// Save handles serialising the queue to disk
func Save(filename string, saveChan SaveChan) {
	for q := range saveChan {
		err := q.Save(filename)
		if err != nil {
			log.Printf("Error saving file to %s: %s", filename, err)
		}
	}
}

// Notify handles sending messages to the Slack channel after a command runs
func Notify(slackConn *slack.Slack, notifyChan NotifyChan) {
	for n := range notifyChan {
		err := slackConn.PostMessage(n.Channel, n.Message)
		if err != nil {
			log.Printf("Error when sending: %s", err)
		}
	}
}

// User handles user renaming in the user cache
func User(userCache *usercache.UserCache, userUpdateChan UserChan) {
	for u := range userUpdateChan {
		old_name := userCache.GetUserName(u.Id)
		userCache.UpdateUserName(u)
		if old_name == "" {
			log.Printf("New user %s cached", u.Name)
		} else {
			log.Printf("User %s renamed to %s", old_name, u.Name)
		}
	}
}
