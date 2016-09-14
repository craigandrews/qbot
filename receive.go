package main

import (
	"log"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
)

func receive(r guac.RealTimeClient, done chan struct{}) (events chan interface{}) {
	events = make(chan interface{})
	go func() {
		for {
			select {
			case <-done:
				log.Println("Terminating listener")
				close(events)
				return
			default:
				event, err := r.Receive()
				if err != nil {
					log.Print("Connection terminated")
					jot.Print("Error while receiving events: ", err)
					close(events)
					return
				}
				events <- event
			}
		}
	}()
	return
}
