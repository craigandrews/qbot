package qbot

import (
	"strings"

	"github.com/doozr/guac"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/util"
)

// CreateMessageDirector creates a message handler that forwards messages to a public or private handler.
func CreateMessageDirector(id string, name string, publicHandler MessageHandler, privateHandler MessageHandler) MessageHandler {
	isDirectedAtUs := func(text string) bool {
		return strings.HasPrefix(text, name) || strings.HasPrefix(text, "<@"+id+">")
	}

	isPrivateChannel := func(channel string) bool {
		return strings.HasPrefix(channel, "D")
	}

	return func(oq queue.Queue, m guac.MessageEvent) (q queue.Queue, err error) {
		q = oq
		if isPrivateChannel(m.Channel) {
			// Private channels should never cause state change
			_, err = privateHandler(q, m)
		} else if isDirectedAtUs(m.Text) {
			// Public channels can cause state change
			_, m.Text = util.StringPop(m.Text)
			q, err = publicHandler(q, m)
		}
		return
	}
}
