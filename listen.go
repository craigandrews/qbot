package main

import (
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
	done chan struct{}, waitGroup *sync.WaitGroup) {

	jot.Print("qbot.listen started")
	defer func() {
		jot.Print("qbot.listen done")
		waitGroup.Done()
	}()

	for {
		// read each incoming message
		select {
		case <-done:
			return

		case event := <-client.Receive():
			switch m := event.(type) {
			case guac.MessageEvent:
				directedAtUs := strings.HasPrefix(m.Text, name) || strings.HasPrefix(m.Text, "<@"+client.ID()+">")
				if directedAtUs {
					_, m.Text = util.StringPop(m.Text)
					messageChan <- m
				} else if util.IsPrivateChannel(m.Channel) {
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
}
