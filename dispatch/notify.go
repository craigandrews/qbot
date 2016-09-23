package dispatch

import (
	"fmt"

	"github.com/doozr/guac"
	"github.com/doozr/qbot/util"
)

// Notifier sends notifications to channels or users
type Notifier func(Notification) error

// NewNotifier creates a new Notifier
func NewNotifier(client guac.RealTimeClient) Notifier {
	return func(notification Notification) error {
		if util.IsUser(notification.Channel) {
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
