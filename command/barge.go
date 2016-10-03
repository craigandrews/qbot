package command

import "github.com/doozr/qbot/queue"

// Barge adds a user to the front of the queue
func (c QueueCommands) Barge(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	i := queue.Item{ID: id, Reason: args}
	q = q.Barge(i)
	if q.Active() == i {
		return q, Notification{ch, c.response.JoinActive(i)}
	}
	c.logActivity(id, args, "barged")
	return q, Notification{ch, c.response.Barge(i, q.Active())}
}
