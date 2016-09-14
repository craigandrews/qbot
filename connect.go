package main

import (
	"log"
	"time"

	"github.com/doozr/guac"
)

func mustConnect(w guac.WebClient, done chan struct{}) (r guac.RealTimeClient, ok bool) {
	backoffTimes := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		5 * time.Second,
		10 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}
	var backoff time.Duration

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
			// If we can increase backoff time, do so, otherwise stick with
			// whatever value we have (last value of backoffTimes)
			if len(backoffTimes) > 0 {
				backoff, backoffTimes = backoffTimes[0], backoffTimes[1:]
			}

			log.Printf("Error while connecting; retrying in %s: %v", backoff, err)
			time.Sleep(backoff)
		}
	}
}
