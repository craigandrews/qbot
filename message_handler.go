package qbot

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

// CreateMessageHandler creates a message handler that calls a command function
func CreateMessageHandler(q queue.Queue, commands CommandMap,
	notify Notifier, persist Persister) MessageHandler {
	return func(m guac.MessageEvent) (err error) {
		text := strings.Trim(m.Text, " \t\r\n")

		var response command.Notification

		cmd, args := util.StringPop(text)
		cmd = strings.ToLower(cmd)

		jot.Printf("message dispatch: message %s with cmd %s and args %v", m.Text, cmd, args)
		if fn, ok := commands[cmd]; ok {
			q, response = fn(q, m.Channel, m.User, args)

			if err = notify(response); err == nil {
				err = persist(q)
			}
		}
		return
	}
}
