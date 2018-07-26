package command

import (
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/util"
)

// Delegate hands over a place in the queue to someone else
func (c QueueCommands) Delegate(q queue.Queue, ch, owner, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, c.response.DelegateNoEntry(owner)}
	}

	name, reason := util.StringPop(args)
	id := c.getIDFromName(name)
	if id == "" {
		return q, Notification{ch, c.response.DelegateNoSuchUser(owner, name)}
	}

	i, ok := c.findItemReverse(q, owner, reason)
	if !ok {
		return q, Notification{ch, c.response.DelegateNoEntry(owner)}
	}

	isActive := q.Active() == i
	n := queue.Item{ID: id, Reason: i.Reason}

	if id == c.id {
		if isActive {
			return q, Notification{ch, c.response.RefuseTokenActive(i, n)}
		}
		return q, Notification{ch, c.response.RefuseToken()}
	}

	q = q.Delegate(i, n)
	c.logActivity(owner, i.Reason, "delegated to "+c.getNameIDPair(id))

	if isActive {
		c.logActivity(n.ID, n.Reason, "is active")
		return q, Notification{ch, c.response.DelegateActive(i, n)}
	}
	return q, Notification{ch, c.response.Delegate(i, id)}
}
