package main

import (
	"sync"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
)

func receive(r guac.RealTimeClient, done chan struct{}, waitGroup *sync.WaitGroup) (events chan interface{}) {
	jot.Print("qbot.receive started")
	events = make(chan interface{})
	waitGroup.Add(1)

	go func() {
		defer func() {
			close(events)
			waitGroup.Done()
			jot.Print("qbot.receive done")
		}()

		for {
			select {
			case <-done:
				jot.Println("qbot.receive: terminating listener")
				return
			default:
				event, err := r.Receive()
				if err != nil {
					jot.Print("qbot.receive: error while receiving events: ", err)
					return
				}
				jot.Print("qbot.receive: event: ", event)
				events <- event
			}
		}
	}()
	return
}
