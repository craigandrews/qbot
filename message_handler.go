package main

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

// CommandMap is a dictionary of command strings to functions
type CommandMap map[string]command.CmdFn

// createMessageDirector creates a message handler that forwards messages to a public or private handler
func createMessageDirector(id string, name string, publicHandler MessageHandler, privateHandler MessageHandler) MessageHandler {
	isDirectedAtUs := func(text string) bool {
		return strings.HasPrefix(text, name) || strings.HasPrefix(text, "<@"+id+">")
	}

	isPrivateChannel := func(channel string) bool {
		return strings.HasPrefix(channel, "D")
	}

	return func(m guac.MessageEvent) (err error) {
		if isPrivateChannel(m.Channel) {
			return privateHandler(m)
		} else if isDirectedAtUs(m.Text) {
			return publicHandler(m)
		}
		return
	}
}

// NewMessageHandler creates a message handler that calls a command function
func createMessageHandler(q queue.Queue, commands CommandMap,
	notify Notifier, persist Persister) MessageHandler {
	return func(m guac.MessageEvent) (err error) {
		text := strings.Trim(m.Text, " \t\r\n")

		var response command.Notification

		cmd, args := util.StringPop(text)
		cmd = strings.ToLower(cmd)

		jot.Printf("message dispatch: message %s with cmd %s and args %v", m.Text, cmd, args)
		if fn, ok := commands[cmd]; ok {
			q, response = fn(q, m.Channel, m.User, args)
		}

		err = notify(response)
		if err == nil {
			err = persist(q)
		}
		return
	}
}
