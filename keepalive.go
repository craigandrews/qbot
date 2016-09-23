package main

import (
	"sync"
	"time"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
)

// keepalive sends a ping request every 30 seconds
func keepalive(client guac.RealTimeClient, done DoneChan, waitGroup *sync.WaitGroup) {
	jot.Print("qbot.keepalive starting up")
	waitGroup.Add(1)
	go func() {
		for {
			select {
			case <-done:
				jot.Print("qbot.keepalive done")
				waitGroup.Done()
				return
			case <-time.After(30 * time.Second):
				jot.Print("keepalive: ping")
				client.Ping()
			}
		}
	}()
}
