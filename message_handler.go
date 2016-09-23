package main

import (
	"strings"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func stringPop(m string) (first string, rest string) {
	parts := strings.SplitN(m, " ", 2)

	if len(parts) < 1 {
		return
	}
	first = parts[0]

	rest = ""
	if len(parts) > 1 {
		rest = strings.Trim(parts[1], " \t\r\n")
	}

	return
}

// NewMessageHandler creates a new MessageHandler
func createMessageHandler(id string, name string, q queue.Queue, commands command.Command,
	notify Notifier, persist Persister) MessageHandler {

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
			cmd, args := stringPop(text)
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
			_, text = stringPop(text)
			cmd, args := stringPop(text)
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
				id, reason := stringPop(args)
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
			err = notify(Notification{
				Channel: channel,
				Message: response,
			})
			if err != nil {
				return
			}
		}

		err = persist(q)

		return
	}
}
