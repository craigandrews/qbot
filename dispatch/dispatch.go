package dispatch

import (
	"log"
	"sync"

	"strings"

	"github.com/doozr/guac"
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
type UserChan chan guac.UserChangeEvent

// Message handles executing user commands and passing on the results
func Message(name string, q queue.Queue, commands command.Command,
	messageChan MessageChan, saveChan SaveChan, notifyChan NotifyChan,
	wg *sync.WaitGroup) {

	log.Print("Starting message dispatch")
	for {
		m, ok := <-messageChan
		if !ok {
			log.Print("Terminating message dispatch")
			wg.Done()
			return
		}

		text := strings.Trim(m.Text, " \t\r\n")
		cmd, args := util.StringPop(text)
		cmd = strings.ToLower(cmd)

		channel := m.Channel
		oldQ := q
		response := ""

		if util.IsPrivateChannel(channel) {
			switch cmd {
			case "list":
				response = commands.List(q)
			case "help":
				response = commands.Help(name)
			case "morehelp":
				response = commands.MoreHelp(name)
			}

		} else {
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
func Save(filename string, saveChan SaveChan, wg *sync.WaitGroup) {
	log.Print("Starting save dispatch")
	for {
		q, ok := <-saveChan
		if !ok {
			log.Print("Terminating save dispatch")
			wg.Done()
			return
		}

		err := q.Save(filename)
		if err != nil {
			log.Printf("Error saving file to %s: %s", filename, err)
		}
	}
}

// Notify handles sending messages to the Slack channel after a command runs
func Notify(client guac.RealTimeClient, notifyChan NotifyChan, wg *sync.WaitGroup) {
	log.Print("Starting notify dispatch")
	for {
		n, ok := <-notifyChan
		if !ok {
			log.Println("Terminating notify dispatch")
			wg.Done()
			return
		}

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
func User(userCache *usercache.UserCache, userUpdateChan UserChan, wg *sync.WaitGroup) {
	log.Print("Starting user dispatch")
	for {
		u, ok := <-userUpdateChan
		if !ok {
			log.Print("Terminating user dispatch")
			wg.Done()
			return
		}

		oldName := userCache.GetUserName(u.ID)
		userCache.UpdateUserName(u)
		if oldName == "" {
			log.Printf("New user %s cached", u.Name)
		} else if oldName != u.Name {
			log.Printf("User %s renamed to %s", oldName, u.Name)
		}
	}
}
