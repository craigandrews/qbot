package main

import (
	"fmt"
	"strings"

	"github.com/doozr/qbot/command"
)

// Notifier sends notifications to channels or users
type Notifier func(command.Notification) error

// IMOpener is a function that opens an IM with a given user
type IMOpener func(string) (string, error)

// MessagePoster is a function that posts a message to a channel
type MessagePoster func(string, string) error

func isUser(channel string) bool {
	return strings.HasPrefix(channel, "U")
}

// CreateNotifier creates a new Notifier
func CreateNotifier(openIM IMOpener, postMessage MessagePoster) Notifier {
	openChannelIfUser := func(user string) (channel string, err error) {
		if !isUser(user) {
			channel = user
			return
		}
		channel, err = openIM(user)
		return
	}

	return func(notification command.Notification) (err error) {
		if notification.Message == "" {
			return
		}
		channel, err := openChannelIfUser(notification.Channel)
		if err != nil {
			return fmt.Errorf("Could not get open channel for %s: %s", notification.Channel, err)
		}

		err = postMessage(channel, notification.Message)
		return err
	}
}
