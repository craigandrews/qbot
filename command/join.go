package command

import "github.com/doozr/qbot/queue"

// Join adds an item to the queue
func (c Command) Join(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	i := queue.Item{ID: id, Reason: args}

	if i.Reason == "" {
		return q, Notification{ch, c.response.JoinNoReason(i)}
	}

	if q.Contains(i) {
		return q, Notification{ch, ""}
	}

	q = q.Add(i)
	c.logActivity(id, args, "joined")
	if q.Active() == i {
		c.logActivity(id, args, "is active")
		return q, Notification{ch, c.response.JoinActive(i)}
	}

	return q, Notification{ch, c.response.Join(i)}
}
