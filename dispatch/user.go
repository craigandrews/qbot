package dispatch

import (
	"log"
	"sync"

	"github.com/doozr/jot"
)

// Save handles serialising the queue to disk
func Save(filename string, saveChan SaveChan, waitGroup *sync.WaitGroup) {

	jot.Print("save dispatch started")
	defer func() {
		waitGroup.Done()
		jot.Print("save dispatch done")
	}()

	for q := range saveChan {
		jot.Print("save dispath: queue to save ", q)
		err := q.Save(filename)
		if err != nil {
			log.Printf("Error saving file to %s: %s", filename, err)
		} else {
			jot.Print("save dispatch: saved to ", filename)
		}
	}
}
