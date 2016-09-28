package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
)

// Receive runs a Receiver instance in a goroutine and handles synchronisation
func Receive(receiver EventReceiver, done DoneChan, waitGroup *sync.WaitGroup) (events guac.EventChan) {
	events = make(guac.EventChan)

	waitGroup.Add(1)
	jot.Print("receive starting up")
	go func() {
		err := receiver(events, done)
		if err != nil {
			log.Print("Error receiving events: ", err)
		}

		close(events)
		jot.Print("receive done")
		waitGroup.Done()
	}()
	return
}

// EventReceiver receives events from Slack and pushes them to a channel
type EventReceiver func(guac.EventChan, DoneChan) error

// Receiver is anything with a Receive method for interface{}
type Receiver interface {
	Receive() (interface{}, error)
}

// CreateEventReceiver creates a default EventReceiver.
func CreateEventReceiver(client Receiver) EventReceiver {
	isDone := func(done DoneChan) bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}

	return func(events guac.EventChan, done DoneChan) (err error) {
		var event interface{}
		for {
			event, err = client.Receive()
			if isDone(done) {
				return nil
			}
			if err != nil {
				return
			}

			if event == nil {
				err = fmt.Errorf("Invalid null event received")
				return
			}

			jot.Print("receiver.listen: received ", event)
			events <- event
		}
	}
}
