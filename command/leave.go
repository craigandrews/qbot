package command

import "github.com/doozr/qbot/queue"

// Leave removes an item from the queue
func (c QueueCommands) Leave(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, ""}
	}

	position, _, ok := c.parsePosition(args)

	var i queue.Item
	if ok {
		i, ok = c.findByPosition(q, position)
		if !ok {
			return q, Notification{ch, c.response.BadIndex(id)}
		}
	} else {
		i, ok = c.findItemReverse(q, id)
		if !ok {
			return q, Notification{ch, c.response.LeaveNoEntry(id)}
		}
	}

	if i.ID != id {
		return q, Notification{ch, c.response.NotOwned(id, position)}
	}

	if q.Active() == i {
		return q, Notification{ch, c.response.LeaveActive(i)}
	}

	q = q.Remove(i)
	c.logActivity(i.ID, i.Reason, "left the queue")
	return q, Notification{ch, c.response.Leave(i)}
}
