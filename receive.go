package main

import (
	"log"

	"github.com/doozr/guac"
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
					log.Print("Error while receiving events: ", err)
					close(events)
					return
				}
				events <- event
			}
		}
	}()
	return
}
