package command

import "github.com/doozr/qbot/queue"

// Barge adds a user to the front of the queue
func (c QueueCommands) Barge(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	var i queue.Item
	i = queue.Item{ID: id, Reason: args}

	position, _, ok := c.parsePosition(args)
	if ok {
		i, ok = c.findByPosition(q, position)
		if !ok {
			return q, Notification{ch, c.response.BadIndex(id)}
		}
	} else if i.Reason == "" {
		n, ok := c.findItem(q, id)
		if !ok {
			return q, Notification{ch, c.response.JoinNoReason(i)}
		}
		i = n
	}

	if i.ID != id {
		return q, Notification{ch, c.response.NotOwned(id, position)}
	}

	q = q.Barge(i)
	c.logActivity(id, args, "barged")
	if q.Active() == i {
		return q, Notification{ch, c.response.JoinActive(i)}
	}
	return q, Notification{ch, c.response.Barge(i, q.Active())}
}
