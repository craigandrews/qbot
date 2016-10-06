package qbot

import (
	"github.com/doozr/guac"
	"github.com/doozr/qbot/queue"
)

// CreateMessagePersister creates a message handler that call another and persists the result
func CreateMessagePersister(persist Persister, fn MessageHandler) MessageHandler {
	return func(oq queue.Queue, m guac.MessageEvent) (q queue.Queue, err error) {
		q, err = fn(oq, m)
		if err != nil {
			return
		}

		err = persist(q)
		return
	}
}
