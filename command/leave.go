package command

import "github.com/doozr/qbot/queue"

// Leave removes an item from the queue
func (c QueueCommands) Leave(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	i, ok := c.findItem(q, id, args)
	if !ok {
		return q, Notification{ch, c.response.LeaveNoEntry(id, args)}
	}

	if q.Active() == i {
		return q, Notification{ch, c.response.LeaveActive(i)}
	}

	if q.Contains(i) {
		q = q.Remove(i)
		c.logActivity(i.ID, i.Reason, "left the queue")
		return q, Notification{ch, c.response.Leave(i)}
	}

	return q, Notification{ch, ""}
}
