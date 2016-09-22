package dispatch

import (
	"strings"
	"sync"

	"github.com/doozr/jot"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/util"
)

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
			jot.Printf("message dispatch: private message %s with cmd %s and args %v", m.Text, cmd, args)
			switch cmd {
			case "list":
				response = commands.List(q)
			case "help":
				response = commands.Help(name)
			case "morehelp":
				response = commands.MoreHelp(name)
			}

		} else {
			jot.Printf("message dispatch: public message %s with cmd %s and args %v", m.Text, cmd, args)
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
