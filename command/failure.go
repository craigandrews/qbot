package command

import "github.com/doozr/qbot/queue"

// Failure notifies token holder and next in line of a problem
func (c QueueCommands) Failure(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	c.logActivity(id, "notification", "failure")

	if len(q) == 0 {
		return q, Notification{ch, c.response.FailureNotificationEmptyQueue(id, args)}
	}

	var ids []string
	if len(q) > 1 && q[0].ID != q[1].ID {
		ids = []string{q[0].ID, q[1].ID}
	} else {
		ids = []string{q[0].ID}
	}

	return q, Notification{ch, c.response.FailureNotification(id, ids, args)}
}
