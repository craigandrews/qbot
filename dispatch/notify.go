package dispatch

import (
	"log"
	"sync"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
	"github.com/doozr/qbot/util"
)

// Notify handles sending messages to the Slack channel after a command runs
func Notify(client guac.RealTimeClient, notifyChan NotifyChan, waitGroup *sync.WaitGroup) {

	jot.Print("notify dispatch started")
	defer func() {
		waitGroup.Done()
		jot.Print("notify dispatch done")
	}()

	for n := range notifyChan {
		if util.IsUser(n.Channel) {
			channel, err := client.IMOpen(n.Channel)
			if err != nil {
				log.Printf("Could not get IM channel for user %s: %s", n.Channel, err)
			} else {
				n.Channel = channel
			}
		}

		err := client.PostMessage(n.Channel, n.Message)
		if err != nil {
			log.Printf("Error when sending: %s", err)
		}
	}
}
