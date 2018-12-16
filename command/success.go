package command

import "github.com/doozr/qbot/queue"

// Success removes the active user from the queue
func (c QueueCommands) Success(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	c.logActivity(id, "notification", "success")

	if len(q) == 0 {
		return q, Notification{ch, c.response.SuccessNotification(id, "")}
	}

	i := q.Active()
	q = q.Remove(i)
	c.logActivity(id, i.Reason, "done")

	if len(q) > 0 {
		n := q.Active()
		c.logActivity(n.ID, n.Reason, "is active")
		return q, Notification{ch, c.response.SuccessNotification(id, c.response.Done(i, q))}
	}

	return q, Notification{ch, c.response.SuccessNotification(id, c.response.DoneNoOthers(i))}
}
