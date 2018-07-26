package command

import "github.com/doozr/qbot/queue"

// Replace swaps the entry at a given position for another one
func (c QueueCommands) Replace(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	position, reason, ok := c.parsePosition(args)
	if !ok {
		return q, Notification{ch, c.response.BadIndex(id)}
	}

	i := queue.Item{ID: id, Reason: reason}

	o, ok := c.findByPosition(q, position)
	if !ok {
		return q, Notification{ch, c.response.BadIndex(id)}
	}

	if i.ID != o.ID {
		return q, Notification{ch, c.response.NotOwned(i.ID, position, o.ID)}
	}

	if reason == "" {
		return q, Notification{ch, c.response.ReplaceNoReason(i)}
	}

	q = q.Delegate(o, i)
	c.logActivity(id, reason, "replaced")
	if q.Active() == i {
		return q, Notification{ch, c.response.JoinActive(i)}
	}
	return q, Notification{ch, c.response.Join(i, position)}
}
