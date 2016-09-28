package qbot

import (
	"encoding/json"
	"log"
	"os"

	"github.com/doozr/jot"
	"github.com/doozr/qbot/queue"
)

// Persister handles exporting the queue to persistent media.
type Persister func(queue.Queue) error

// WriteFile is a type to allow replacement of the WriteFile function for reasons.
type WriteFile func(string, []byte, os.FileMode) error

// CreatePersister creates a new Persister.
func CreatePersister(writeFile WriteFile, filename string, oldQ queue.Queue) Persister {
	return func(q queue.Queue) (err error) {
		jot.Print("persist: queue to save ", q)
		if oldQ.Equal(q) {
			return
		}

		j, err := json.Marshal(q)
		if err != nil {
			log.Print("Error serialising qeuue: ", err)
			return
		}

		jot.Printf("queue: writing queue JSON %v to %s", string(j), filename)
		err = writeFile(filename, j, 0644)
		if err != nil {
			log.Printf("Error saving file to %s: %s", filename, err)
			return
		}

		oldQ = q
		jot.Print("persist: saved to ", filename)
		return
	}
}
