package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
	"github.com/doozr/qbot/dispatch"
	"github.com/doozr/qbot/util"
)

func listen(name string, client guac.RealTimeClient,
	messageChan dispatch.MessageChan, userChan dispatch.UserChan,
	done chan struct{}, waitGroup *sync.WaitGroup) (abort chan struct{}) {

	abort = make(chan struct{})

	jot.Print("qbot.listen started")
	waitGroup.Add(1)
	go func() {
		defer func() {
			jot.Print("qbot.listen done")
			waitGroup.Done()
		}()

		events := receive(client, done, waitGroup)
		for {
			jot.Print("qbot.listen awaiting event")
			select {
			case <-done:
				return

			case event, ok := <-events:
				if !ok {
					log.Print("Incoming event stream closed")
					jot.Print("qbot.listen: closing abort channel")
					close(abort)
					return
				}
				switch m := event.(type) {
				case guac.MessageEvent:
					jot.Print("qbot.listen received message: ", m)
					directedAtUs := strings.HasPrefix(m.Text, name) || strings.HasPrefix(m.Text, "<@"+client.ID()+">")
					jot.Print("qbot.listen message directed at us? ", name)
					if directedAtUs {
						jot.Printf("qbot.listen received public message from %s in channel %s: %v", m.User, m.Channel, m.Text)
						_, m.Text = util.StringPop(m.Text)
						messageChan <- m
					} else if util.IsPrivateChannel(m.Channel) {
						jot.Printf("qbot.listen received private message from %s in channel %s: %v", m.User, m.Channel, m.Text)
						messageChan <- m
					}

				case guac.UserChangeEvent:
					userChan <- m.UserInfo

				case guac.PingPongEvent:
					jot.Print("qbot.listen: pong")
				}

			case <-time.After(30 * time.Second):
				jot.Print("qbot.listen: ping")
				client.Ping()
			}
		}
	}()

	return
}

func receive(client guac.RealTimeClient, done chan struct{}, waitGroup *sync.WaitGroup) (events chan interface{}) {
	events = make(chan interface{})
	jot.Print("qbot.receive started")

	waitGroup.Add(1)
	go func() {
		defer func() {
			close(events)
			waitGroup.Done()
			jot.Print("qbot.receive done")
		}()

		for {
			select {
			case <-done:
				return

			default:
				event, err := client.Receive()
				if err != nil {
					log.Print("Error receiving event: ", err)
					return
				}

				if event == nil {
					log.Print("Nil event received")
					return
				}

				jot.Print("qbot.receive: received ", event)
				events <- event
			}
		}
	}()

	return
}
