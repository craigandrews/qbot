package qbot

import (
	"fmt"
	"sync"
	"time"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
	"github.com/doozr/qbot/queue"
)

// Dispatch runs a Dispatcher in a goroutine and handles synchronisation
func Dispatch(dispatcher Dispatcher, events guac.EventChan, done DoneChan, waitGroup *sync.WaitGroup) (abort chan error) {
	abort = make(chan error)

	waitGroup.Add(1)
	jot.Print("dispatch starting up")
	go func() {
		err := dispatcher(events, done)
		if err != nil {
			abort <- err
		}

		close(abort)
		jot.Print("dispatch done")
		waitGroup.Done()
	}()
	return
}

// Dispatcher sends incoming messages to the correct recipient.
type Dispatcher func(guac.EventChan, DoneChan) error

// CreateDispatcher creates a new Dispatcher instance.
func CreateDispatcher(q queue.Queue, timeout time.Duration, handleMessage MessageHandler, handleUserChange UserChangeHandler) Dispatcher {

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
					q, err = handleMessage(q, m)

				case guac.UserChangeEvent:
					handleUserChange(m.UserInfo)

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
