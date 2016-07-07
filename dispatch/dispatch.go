package dispatch

import (
	"log"

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

	for m := range messageChan {
		cmd, args := splitCommand(m.Text)

		old_q := q
		response := ""

		switch cmd {
		case "join":
			q, response = commands.Join(q, m.User, args)
		case "leave":
			q, response = commands.Leave(q, m.User, args)
		case "done":
			q, response = commands.Done(q, m.User)
		case "yield":
			q, response = commands.Yield(q, m.User)
		case "barge":
			q, response = commands.Barge(q, m.User, args)
		case "boot":
			id, reason := splitUser(args)
			q, response = commands.Boot(q, m.User, id, reason)
		case "oust":
			q, response = commands.Oust(q, m.User, args)
		case "list":
			response = commands.List(q)
		case "help":
			response = commands.Help(name)
		}

		if response != "" {
			if !q.Equal(old_q) {
				logResponse(response)
				saveChan <- q
			}
			notifyChan <- Notification{m.Channel, response}
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
