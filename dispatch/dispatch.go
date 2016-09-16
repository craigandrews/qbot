package dispatch

import (
	"log"
	"sync"

	"strings"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
	"github.com/doozr/qbot/util"
)

const privateOnly = "Please send this request in private"

// Notification represents a message to a channel
type Notification struct {
	Channel string
	Message string
}

// MessageChan is a stream of Slack real-time messages
type MessageChan chan guac.MessageEvent

// SaveChan is a stream of queue instances to persist
type SaveChan chan queue.Queue

// NotifyChan is a stream of notifications
type NotifyChan chan Notification

// UserChan is a stream of user info updates
type UserChan chan guac.UserInfo

// Message handles executing user commands and passing on the results
func Message(name string, q queue.Queue, commands command.Command,
	messageChan MessageChan, saveChan SaveChan, notifyChan NotifyChan,
	waitGroup *sync.WaitGroup) {

	jot.Print("message dispatch started")
	defer func() {
		waitGroup.Done()
		jot.Print("message dispatch done")
	}()

	for m := range messageChan {
		text := strings.Trim(m.Text, " \t\r\n")
		cmd, args := util.StringPop(text)

		channel := m.Channel
		oldQ := q
		response := ""

		if util.IsPrivateChannel(channel) {
			jot.Print("message dispatch: private message ", m)
			switch cmd {
			case "list":
				response = commands.List(q)
			case "help":
				response = commands.Help(name)
			case "morehelp":
				response = commands.MoreHelp(name)
			}

		} else {
			jot.Print("message dispatch: public message ", m)
			switch cmd {
			case "join":
				q, response = commands.Join(q, m.User, args)
			case "leave":
				q, response = commands.Leave(q, m.User, args)
			case "done":
				q, response = commands.Done(q, m.User)
			case "drop":
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
				channel = m.User
			case "morehelp":
				response = commands.MoreHelp(name)
				channel = m.User
			}
		}

		if response != "" {
			if !q.Equal(oldQ) {
				util.LogMultiLine(response)
				saveChan <- q
			}
			notifyChan <- Notification{channel, response}
		}
	}
}

// Save handles serialising the queue to disk
func Save(filename string, saveChan SaveChan, waitGroup *sync.WaitGroup) {

	jot.Print("save dispatch started")
	defer func() {
		waitGroup.Done()
		jot.Print("save dispatch done")
	}()

	for q := range saveChan {
		jot.Print("save dispath: queue to save ", q)
		err := q.Save(filename)
		if err != nil {
			log.Printf("Error saving file to %s: %s", filename, err)
		} else {
			jot.Print("save dispatch: saved to ", filename)
		}
	}
}

// Notify handles sending messages to the Slack channel after a command runs
func Notify(client guac.RealTimeClient, notifyChan NotifyChan, waitGroup *sync.WaitGroup) {

	jot.Print("notify dispatch started")
	defer func() {
		waitGroup.Done()
		jot.Print("notify dispatch done")
	}()

	for n := range notifyChan {
		if util.IsUser(n.Channel) {
			channel, err := client.IMOpen(n.Channel)
			if err != nil {
				log.Printf("Could not get IM channel for user %s: %s", n.Channel, err)
			} else {
				n.Channel = channel
			}
		}

		err := client.PostMessage(n.Channel, n.Message)
		if err != nil {
			log.Printf("Error when sending: %s", err)
		}
	}
}

// User handles user renaming in the user cache
func User(userCache *usercache.UserCache, userUpdateChan UserChan, waitGroup *sync.WaitGroup) {

	jot.Print("user dispatch started")
	defer func() {
		waitGroup.Done()
		jot.Print("user dispatch done")
	}()

	for u := range userUpdateChan {
		oldName := userCache.GetUserName(u.ID)
		userCache.UpdateUserName(u.ID, u.Name)
		if oldName == "" {
			log.Printf("New user %s cached", u.Name)
		} else if oldName != u.Name {
			log.Printf("User %s renamed to %s", oldName, u.Name)
		}
	}
}
