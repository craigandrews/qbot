package command

import (
	"github.com/doozr/qbot/queue"
)

// Boot kicks someone from the waiting list
func (c QueueCommands) Boot(q queue.Queue, ch, booter, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, ""}
	}

	position, name, _ := c.parsePosition(args)

	id := c.getIDFromName(name)
	if id == "" {
		return q, Notification{ch, c.response.BootNoEntry(booter, name)}
	}

	i, ok := c.findByPosition(q, position)
	if !ok {
		i, ok = c.findItemReverse(q, id)
		if !ok {
			return q, Notification{ch, c.response.BootNoEntry(booter, name)}
		}
	}

	if i.ID != id {
		return q, Notification{ch, c.response.NotOwned(booter, position, i.ID)}
	}

	if q.Active() == i {
		return q, Notification{ch, c.response.OustNotBoot(booter)}
	}

	q = q.Remove(i)
	c.logActivity(id, i.Reason, "booted by "+c.getNameIDPair(booter))
	return q, Notification{ch, c.response.Boot(booter, i)}
}
