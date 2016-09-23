package main

import (
	"fmt"
	"time"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
	"github.com/doozr/qbot/queue"
)

// Notification represents a message to a channel
type Notification struct {
	Channel string
	Message string
}

// MessageHandler handles an incoming message event
type MessageHandler func(guac.MessageEvent) error

// Notifier sends notifications to channels or users
type Notifier func(Notification) error

// Persister handles exporting the queue to persistent media
type Persister func(queue.Queue) error

// UserChangeHandler handles incoming user change events
type UserChangeHandler func(guac.UserChangeEvent) error

// Dispatcher sends incoming messages to the correct recipient
type Dispatcher func(guac.EventChan, DoneChan) error

// createDispatcher creates a new Dispatcher instance
func createDispatcher(client guac.RealTimeClient, timeout time.Duration,
	handleMessage MessageHandler, handleUserChange UserChangeHandler) Dispatcher {

	return func(events guac.EventChan, done DoneChan) (err error) {
		for {
			jot.Print("dispatcher awaiting event")
			select {
			case <-done:
				jot.Print("dispatcher shutting down")
				return

			case event, ok := <-events:
				if !ok {
					jot.Print("dispatcher: closing abort channel")
					return
				}

				switch m := event.(type) {
				case guac.MessageEvent:
					jot.Print("dispatcher received message: ", m)
					err = handleMessage(m)

				case guac.UserChangeEvent:
					err = handleUserChange(m)

				case guac.PingPongEvent:
					jot.Print("dispatcher: pong")
				}

			case <-time.After(timeout):
				err = fmt.Errorf("No activity for %s - shutting down", timeout)
			}

			if err != nil {
				return
			}
		}
	}
}
