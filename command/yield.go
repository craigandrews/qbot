package command

import "github.com/doozr/qbot/queue"

// Yield allows the second place ahead of the active user
func (c Command) Yield(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, ""}
	}
	i := q.Active()
	if i.ID != id {
		return q, Notification{ch, c.response.YieldNotActive(queue.Item{ID: id, Reason: ""})}
	}
	if len(q) < 2 {
		return q, Notification{ch, c.response.YieldNoOthers(i)}
	}
	q = q.Yield()
	n := q.Active()
	c.logActivity(id, i.Reason, "yielded")
	c.logActivity(n.ID, n.Reason, "is active")
	return q, Notification{ch, c.response.Yield(i, q)}
}
