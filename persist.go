package main

import (
	"log"

	"github.com/doozr/jot"
	"github.com/doozr/qbot/queue"
)

// Persister handles exporting the queue to persistent media
type Persister func(queue.Queue) error

// NewPersister creates a new Persister
func createPersister(filename string) Persister {
	var oldQ queue.Queue
	return func(q queue.Queue) (err error) {
		jot.Print("persist: queue to save ", q)
		if !oldQ.Equal(q) {
			err = q.Save(filename)
			if err != nil {
				log.Printf("Error saving file to %s: %s", filename, err)
				return
			}
		}
		oldQ = q
		jot.Print("persist: saved to ", filename)
		return
	}
}
