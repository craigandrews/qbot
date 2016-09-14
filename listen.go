package main

import (
	"log"
	"strings"
	"time"

	"github.com/doozr/guac"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/dispatch"
	"github.com/doozr/qbot/notification"
	"github.com/doozr/qbot/util"
)

// Returns true if text begins with our name (case insensitive) or a link to our ID
//
// This means, if our name is "blork", we could be referenced in all these ways:
//
// * blork list
// * blork: list
// * Blork: list
// * @blork list
func messageIsForUs(r guac.RealTimeClient, message guac.MessageEvent) bool {
	return strings.HasPrefix(strings.ToLower(message.Text), r.Name()) ||
		strings.HasPrefix(message.Text, "<@"+r.ID()+">")
}

func listen(r guac.RealTimeClient, world World) {
	notifications := notification.New(world.UserCache)
	commands := command.New(notifications, world.UserCache)

	messageChan := make(dispatch.MessageChan, 100)
	defer close(messageChan)

	notifyChan := make(dispatch.NotifyChan, 5)
	defer close(notifyChan)

	world.WaitGroup.Add(2)
	go dispatch.Message(r.Name(), world.Q, commands, messageChan, world.SaveChan, notifyChan, world.WaitGroup)
	go dispatch.Notify(r, notifyChan, world.WaitGroup)

	// Get incoming stream of events
	//
	// When this channel closes, the connection is dead so we should stop
	events := receive(r, world.Done)

	log.Println("Ready to receive events")
	for {
		select {
		case <-world.Done:
			log.Println("Closing connection")
			r.Close()
			return

		case e, ok := <-events:
			if !ok {
				return
			}

			switch event := e.(type) {
			case guac.MessageEvent:
				isForUs := messageIsForUs(r, event)
				isDM := util.IsPrivateChannel(event.Channel)

				if isForUs {
					_, event.Text = util.StringPop(event.Text)
				}

				if isForUs || isDM {
					messageChan <- event
				}

			case guac.UserChangeEvent:
				world.UserChan <- event
			}

		case <-time.After(30 * time.Second):
			r.Ping()
		}
	}
}
