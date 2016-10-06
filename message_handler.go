package qbot

import (
	"strings"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/util"
)

// MessageHandler handles an incoming message event.
type MessageHandler func(queue.Queue, guac.MessageEvent) (queue.Queue, error)

// CommandMap is a dictionary of command strings to functions.
type CommandMap map[string]command.Command

// CreateMessageHandler creates a message handler that calls a command function.
func CreateMessageHandler(commands CommandMap,
	notify Notifier, persist Persister) MessageHandler {
	return func(oq queue.Queue, m guac.MessageEvent) (q queue.Queue, err error) {
		text := strings.Trim(m.Text, " \t\r\n")

		var response command.Notification

		cmd, args := util.StringPop(text)
		cmd = strings.ToLower(cmd)

		jot.Printf("message dispatch: message %s with cmd %s and args %v", m.Text, cmd, args)
		fn, ok := commands[cmd]
		if !ok {
			q = oq
			return
		}

		q, response = fn(oq, m.Channel, m.User, args)

		err = notify(response)
		if err != nil {
			return
		}

		err = persist(q)
		return
	}
}
