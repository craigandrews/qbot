package main

import (
	"sync"
	"time"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
)

func keepalive(client guac.RealTimeClient, done DoneChan, waitGroup *sync.WaitGroup) {
	jot.Print("qbot.keepalive starting up")
	defer func() {
		waitGroup.Done()
		jot.Print("qbot.keepalive done")
	}()

	for {
		select {
		case <-done:
			jot.Print("qbot.keepalive shutting down")
			return
		case <-time.After(30 * time.Second):
			jot.Print("keepalive: ping")
			client.Ping()
		}
	}
}
