package dispatch

import (
	"log"

	"strings"

	"github.com/doozr/goslack"
	"github.com/doozr/goslack/apitypes"
	"github.com/doozr/goslack/rtmtypes"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
	"github.com/doozr/qbot/util"
)

// Notification represents a message to a channel
type Notification struct {
	Channel string
	Message string
}

// MessageChan is a stream of Slack real-time messages
type MessageChan chan rtmtypes.RtmMessage

// SaveChan is a stream of queue instances to persist
type SaveChan chan queue.Queue

// NotifyChan is a stream of notifications
type NotifyChan chan Notification

// UserChan is a stream of user info updates
type UserChan chan apitypes.UserInfo

// Message handles executing user commands and passing on the results
func Message(name string, q queue.Queue, commands command.Command,
	messageChan MessageChan, saveChan SaveChan, notifyChan NotifyChan) {

	for m := range messageChan {
		text := strings.Trim(m.Text, " \t\r\n")
		cmd, args := util.StringPop(text)

		oldQ := q
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
			id, reason := util.StringPop(args)
			q, response = commands.Boot(q, m.User, id, reason)
		case "oust":
			q, response = commands.Oust(q, m.User, args)
		case "list":
			response = commands.List(q)
		case "help":
			response = commands.Help(name)
		case "morehelp":
			response = commands.MoreHelp(name)
		}

		if response != "" {
			if !q.Equal(oldQ) {
				util.LogMultiLine(response)
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
func Notify(slackConn *goslack.Slack, notifyChan NotifyChan) {
	for n := range notifyChan {
		err := slackConn.PostRealTimeMessage(n.Channel, n.Message)
		if err != nil {
			log.Printf("Error when sending: %s", err)
		}
	}
}

// User handles user renaming in the user cache
func User(userCache *usercache.UserCache, userUpdateChan UserChan) {
	for u := range userUpdateChan {
		oldName := userCache.GetUserName(u.ID)
		userCache.UpdateUserName(u)
		if oldName == "" {
			log.Printf("New user %s cached", u.Name)
		} else {
			log.Printf("User %s renamed to %s", oldName, u.Name)
		}
	}
}
