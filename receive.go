package main

import (
	"fmt"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
)

// Receiver of events from Slack
type Receiver func(guac.EventChan, DoneChan) error

// New receiver instance.
func createReceiver(client guac.RealTimeClient) Receiver {
	return func(events guac.EventChan, done DoneChan) (err error) {
		var event interface{}
		for {
			event, err = client.Receive()
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
