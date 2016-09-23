package dispatch

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
)

// Dispatcher sends incoming messages to the correct recipient
type Dispatcher struct {
	handleMessage    MessageHandler
	handleUserChange UserChangeHandler
}

// New Listener instance
func New(messageHandler MessageHandler, userChangeHandler UserChangeHandler) Dispatcher {
	return Dispatcher{
		handleMessage:    messageHandler,
		handleUserChange: userChangeHandler,
	}
}

// Listen for incoming messages and dispatch them
func (d Dispatcher) Listen(client guac.RealTimeClient, timeout time.Duration, done chan struct{}, waitGroup *sync.WaitGroup) (abort chan error) {
	abort = make(chan error)
	events := d.receive(client, done, waitGroup)

	waitGroup.Add(1)
	jot.Print("dispatcher.listen starting up")
	go func() {
		defer func() {
			waitGroup.Done()
			close(abort)
			jot.Print("dispatcher.listen done")
		}()

		var err error
		for {
			jot.Print("dispatcher.listen awaiting event")
			select {
			case <-done:
				jot.Print("dispatcher.listen shutting down")
				return

			case event, ok := <-events:
				if !ok {
					jot.Print("dispatcher.listen: closing abort channel")
					return
				}

				switch m := event.(type) {
				case guac.MessageEvent:
					jot.Print("dispatcher.listen received message: ", m)
					err = d.handleMessage(m)

				case guac.UserChangeEvent:
					err = d.handleUserChange(m)

				case guac.PingPongEvent:
					jot.Print("dispatcher.listen: pong")
				}

			case <-time.After(timeout):
				err = fmt.Errorf("No activity for %s - shutting down", timeout)
				return
			}

			if err != nil {
				abort <- err
				return
			}
		}
	}()
	return abort
}

func (d Dispatcher) receive(client guac.RealTimeClient, done chan struct{}, waitGroup *sync.WaitGroup) (events chan interface{}) {
	events = make(chan interface{})
	jot.Print("dispatcher.receive started")

	waitGroup.Add(1)
	jot.Print("dispatcher.receive starting up")
	go func() {
		defer func() {
			waitGroup.Done()
			close(events)
			jot.Print("dispatcher.receive done")
		}()

		for {
			select {
			case <-done:
				jot.Print("dispatcher.receive shutting down")
				return
			default:
				event, err := client.Receive()
				if err != nil {
					log.Print("Error receiving event: ", err)
					return
				}

				if event == nil {
					log.Print("Nil event received")
					return
				}

				jot.Print("dispatcher.receive: received ", event)
				events <- event
			}
		}
	}()

	return
}
