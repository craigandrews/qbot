package dispatch

import (
	"strings"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/util"
)

// MessageHandler handles an incoming message event
type MessageHandler func(guac.MessageEvent) error

// NewMessageHandler creates a new MessageHandler
func NewMessageHandler(id string, name string, q queue.Queue, commands command.Command, notify Notifier, persist Persister) MessageHandler {
	isDirectedAtUs := func(id string, name string, message guac.MessageEvent) bool {
		return strings.HasPrefix(message.Text, name) || strings.HasPrefix(message.Text, "<@"+id+">")
	}

	isPrivateChannel := func(channel string) bool {
		return strings.HasPrefix(channel, "D")
	}

	return func(m guac.MessageEvent) (err error) {
		text := strings.Trim(m.Text, " \t\r\n")

		channel := m.Channel
		response := ""

		if isPrivateChannel(channel) {
			cmd, args := util.StringPop(text)
			cmd = strings.ToLower(cmd)

			jot.Printf("message dispatch: private message %s with cmd %s and args %v", m.Text, cmd, args)
			switch cmd {
			case "list":
				response = commands.List(q)
			case "help":
				response = commands.Help(name)
			case "morehelp":
				response = commands.MoreHelp(name)
			}

		} else if isDirectedAtUs(id, name, m) {
			_, text = util.StringPop(text)
			cmd, args := util.StringPop(text)
			cmd = strings.ToLower(cmd)

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
			err = notify(Notification{channel, response})
			if err != nil {
				return
			}
		}

		err = persist(q)

		return
	}
}
