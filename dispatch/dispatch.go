package dispatch

import (
	"github.com/doozr/guac"
	"github.com/doozr/qbot/queue"
)

// Notification represents a message to a channel
type Notification struct {
	Channel string
	Message string
}

// MessageChan is a stream of Slack real-time messages
type MessageChan chan guac.MessageEvent

// SaveChan is a stream of queue instances to persist
type SaveChan chan queue.Queue

// NotifyChan is a stream of notifications
type NotifyChan chan Notification

// UserChan is a stream of user info updates
type UserChan chan guac.UserInfo
