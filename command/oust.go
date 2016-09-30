package command

import "github.com/doozr/qbot/queue"

// Oust boots the current token holder and gives it to the next person
func (c Command) Oust(q queue.Queue, ch, ouster, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, ""}
	}

	if args == "" {
		return q, Notification{ch, c.response.OustNoTarget(ouster)}
	}

	id := c.getIDFromName(args)
	if id == "" {
		return q, Notification{ch, c.response.OustNotActive(ouster)}
	}

	i := q.Active()
	if i.ID != id {
		return q, Notification{ch, c.response.OustNotActive(ouster)}
	}

	if len(q) == 1 {
		q = q.Remove(i)
		c.logActivity(i.ID, i.Reason, "ousted by "+c.getNameIDPair(ouster))
		return q, Notification{ch, c.response.OustNoOthers(ouster, i)}
	}

	q = q.Yield()
	n := q.Active()
	c.logActivity(n.ID, n.Reason, "is active")
	return q, Notification{ch, c.response.Oust(ouster, i, q)}
}
