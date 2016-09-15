package main

import (
	"sync"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
)

func receive(r guac.RealTimeClient, done chan struct{}, waitGroup *sync.WaitGroup) (events chan interface{}) {
	jot.Print("receive started")
	events = make(chan interface{})
	waitGroup.Add(1)

	go func() {
		defer func() {
			close(events)
			waitGroup.Done()
			jot.Print("receive done")
		}()

		for {
			select {
			case <-done:
				jot.Println("Terminating listener")
				return
			default:
				event, err := r.Receive()
				if err != nil {
					jot.Print("Error while receiving events: ", err)
					return
				}
				jot.Print("Received event: ", event)
				events <- event
			}
		}
	}()
	return
}
