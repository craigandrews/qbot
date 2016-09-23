package main

import (
	"fmt"
	"strings"

	"github.com/doozr/guac"
)

func isUser(channel string) bool {
	return strings.HasPrefix(channel, "U")
}

// NewNotifier creates a new Notifier
func createNotifier(client guac.RealTimeClient) Notifier {
	return func(notification Notification) error {
		if isUser(notification.Channel) {
			channel, err := client.IMOpen(notification.Channel)
			if err != nil {
				return fmt.Errorf("Could not get IM channel for user %s: %s", notification.Channel, err)
			}
			notification.Channel = channel
		}

		err := client.PostMessage(notification.Channel, notification.Message)
		return err
	}
}
