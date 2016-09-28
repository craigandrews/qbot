package main

import (
	"sync"
	"time"

	"github.com/doozr/jot"
)

// Pinger is a thing that pings
type Pinger func() error

// After is a thing that returns a channel that emits the time after a duration
type After func(time.Duration) <-chan time.Time

// StartKeepAlive sends a ping request every 30 seconds
func StartKeepAlive(ping Pinger, after After, done DoneChan, waitGroup *sync.WaitGroup) {
	jot.Print("qbot.keepalive starting up")
	waitGroup.Add(1)
	go func() {
		for {
			select {
			case <-done:
				jot.Print("qbot.keepalive done")
				waitGroup.Done()
				return
			case <-after(30 * time.Second):
				jot.Print("keepalive: ping")
				ping()
			}
		}
	}()
}
