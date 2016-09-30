package command

import (
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/util"
)

// Boot kicks someone from the waiting list
func (c Command) Boot(q queue.Queue, ch, booter, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, ""}
	}

	name, reason := util.StringPop(args)
	id := c.getIDFromName(name)
	i, ok := c.findItem(q, id, reason)
	if !ok {
		return q, Notification{ch, c.response.BootNoEntry(booter, name, reason)}
	}

	if q.Active() == i {
		return q, Notification{ch, c.response.OustNotBoot(booter)}
	}

	q = q.Remove(i)
	c.logActivity(id, reason, "booted by "+c.getNameIDPair(booter))
	return q, Notification{ch, c.response.Boot(booter, i)}
}
