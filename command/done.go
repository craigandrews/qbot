package command

import "github.com/doozr/qbot/queue"

// Done removes the active user from the queue
func (c Command) Done(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, ""}
	}

	i := q.Active()

	if i.ID != id {
		return q, Notification{ch, c.response.DoneNotActive(id)}
	}

	q = q.Remove(i)
	c.logActivity(id, i.Reason, "done")
	if len(q) > 0 {
		n := q.Active()
		c.logActivity(n.ID, n.Reason, "is active")
		return q, Notification{ch, c.response.Done(i, q)}
	}
	return q, Notification{ch, c.response.DoneNoOthers(i)}
}
