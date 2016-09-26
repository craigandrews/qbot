package main

import (
	"strings"

	"github.com/doozr/guac"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/util"
)

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
			_, m.Text = util.StringPop(m.Text)
			return publicHandler(m)
		}
		return
	}
}
