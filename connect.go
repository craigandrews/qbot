package main

import (
	"log"
	"time"

	"github.com/doozr/guac"
)

func mustConnect(w guac.WebClient, done chan struct{}) (r guac.RealTimeClient, ok bool) {
	backoff := 1 * time.Second
	maxBackoff := 64 * time.Second

	log.Print("Connecting to Slack")
	var err error
	for {
		r, err = w.RealTime()
		if err == nil {
			log.Print("Connected")
			ok = true
			return
		}

		select {
		case <-done:
			log.Println("Cancelling reconnection")
			return
		default:
			log.Printf("Error while connecting; retrying in %s: %v", backoff, err)
			time.Sleep(backoff)
			if backoff <= maxBackoff {
				backoff *= 2
			}
		}
	}
}
