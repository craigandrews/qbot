package command

import (
	"fmt"

	"github.com/doozr/qbot/queue"
)

// List shows who has the token and who is waiting
func (c Command) List(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, "Nobody has the token, and nobody is waiting"}
	}

	a := q.Active()
	s := fmt.Sprintf("*%d: %s (%s) has the token*", 1, c.userCache.GetUserName(a.ID), a.Reason)
	for ix, i := range q.Waiting() {
		s += fmt.Sprintf("\n%d: %s (%s)", ix+2, c.userCache.GetUserName(i.ID), i.Reason)
	}
	return q, Notification{ch, s}
}
