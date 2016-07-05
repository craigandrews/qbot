package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/doozr/qbot/notification"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
)

type PendingOust struct {
	Item      queue.Item
	Timestamp time.Time
}

type Command struct {
	Notification notification.Notification
	UserCache    *usercache.UserCache
	PendingOusts map[string]PendingOust
}

func New(n notification.Notification, uc *usercache.UserCache) Command {
	c := Command{n, uc, make(map[string]PendingOust)}
	return c
}

func (c Command) findItem(q queue.Queue, id, reason string) (item queue.Item) {
	for ix := len(q) - 1; ix >= 0; ix-- {
		if q[ix].Id == id && strings.HasPrefix(q[ix].Reason, reason) {
			item = q[ix]
			break
		}
	}
	return
}

// Join adds an item to the queue
func (c Command) Join(q queue.Queue, id, reason string) (queue.Queue, string) {
	i := queue.Item{id, reason}

	if i.Reason == "" {
		return q, c.Notification.JoinNoReason(i)
	}

	if q.Contains(i) {
		return q, ""
	}

	q = q.Add(i)
	if q.Active() == i {
		return q, c.Notification.JoinActive(i)
	}

	return q, c.Notification.Join(i)
}

// Leave removes an item from the queue
func (c Command) Leave(q queue.Queue, id, reason string) (queue.Queue, string) {
	i := c.findItem(q, id, reason)
	if i.Id == "" {
		return q, ""
	}

	active := q.Active() == i
	if q.Contains(i) {
		q = q.Remove(i)
		if active {
			if len(q) == 0 {
				return q, c.Notification.LeaveNoActive(i)
			}
			return q, c.Notification.LeaveActive(i, q)
		}
		return q, c.Notification.Leave(i)
	}
	return q, ""
}

// Done removes the active user from the queue
func (c Command) Done(q queue.Queue, id string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, ""
	}

	i := q.Active()

	if i.Id != id {
		return q, c.Notification.DoneNotActive(i)
	}

	q = q.Remove(i)
	if len(q) > 0 {
		return q, c.Notification.Done(i, q)
	}
	return q, c.Notification.DoneNoOthers(i)
}

// Yield allows the second place ahead of the active user
func (c Command) Yield(q queue.Queue, id string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, c.Notification.YieldNotActive(queue.Item{id, ""})
	}
	i := q.Active()
	if i.Id != id {
		return q, c.Notification.YieldNotActive(queue.Item{id, ""})
	}
	if len(q) < 2 {
		return q, c.Notification.YieldNoOthers(i)
	}
	q = q.Yield()
	return q, c.Notification.Yield(i, q)
}

// Barge adds a user to the front of the queue
func (c Command) Barge(q queue.Queue, id, reason string) (queue.Queue, string) {
	i := queue.Item{id, reason}
	q = q.Barge(i)
	if q.Active() == i {
		return q, c.Notification.JoinActive(i)
	}
	return q, c.Notification.Barge(i)
}

// Boot kicks someone from the waiting list
func (c Command) Boot(q queue.Queue, booter, name, reason string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, ""
	}

	id := c.UserCache.GetUserId(name)
	i := c.findItem(q, id, reason)
	if i.Id == "" {
		return q, ""
	}

	if q.Active() == i {
		return q, c.Notification.OustNotBoot(booter)
	}

	if q.Contains(i) {
		q = q.Remove(i)
		return q, c.Notification.Boot(booter, i)
	}
	return q, ""
}

// Oust boots the current token holder and gives it to the next person
func (c Command) Oust(q queue.Queue, ouster, name string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, ""
	}

	id := c.UserCache.GetUserId(name)
	if id == "" {
		return q, c.Notification.OustNotActive(ouster)
	}

	i := q.Active()
	if i.Id != id {
		return q, c.Notification.OustNotActive(ouster)
	}

	// If a previous request has been lodged in the last 30 seconds
	// and all is well then oust the active user
	pendingOust, ok := c.PendingOusts[ouster]
	if ok && pendingOust.Item == i {
		if time.Since(pendingOust.Timestamp).Seconds() < 30 && q.Active() == i {
			q = q.Remove(i)

			if len(q) == 0 {
				return q, c.Notification.OustNoOthers(ouster, i)
			}
			return q, c.Notification.Oust(ouster, i, q)
		}
	}

	c.PendingOusts[ouster] = PendingOust{i, time.Now()}

	return q, c.Notification.OustConfirm(ouster, i)
}

// List shows who has the token and who is waiting
func (c Command) List(q queue.Queue) string {
	if len(q) == 0 {
		return "Nobody has the token, and nobody is waiting"
	}

	a := q.Active()
	s := fmt.Sprintf("*%d: %s (%s) has the token*", 1, c.UserCache.GetUserName(a.Id), a.Reason)
	for ix, i := range q.Waiting() {
		s += fmt.Sprintf("\n%d: %s (%s)", ix+2, c.UserCache.GetUserName(i.Id), i.Reason)
	}
	return s
}

// Help provides much needed assistance
func (c Command) Help(name string) string {
	cmds := [][]string{
		[]string{"join <reason>", "Join the queue and give a reason why"},
		[]string{"leave", "Leave the queue (your most recent entry is removed)"},
		[]string{"leave <reason>", "Leave the queue (your most recent entry starting with <reason> is removed)"},
		[]string{"done", "Release the token once you are done with it"},
		[]string{"yield", "Release the token and swap places with next in line"},
		[]string{"barge <reason>", "Barge to the front of the queue so you get the token next (only with good reason!)"},
		[]string{"boot <name>", "Kick somebody out of the waiting list (their most recent entry is removed)"},
		[]string{"boot <name> <reason>", "Kick somebody out of the waiting list (their most recent entry starting with <reason> is removed"},
		[]string{"oust <name>", "Forcibly take the token from the token holder and kick them out of the queue (only with VERY good reason!)"},
		[]string{"list", "Show who has the token and who is waiting"},
		[]string{"help", "Show this text"},
	}
	s := fmt.Sprintf("Address each command to the bot (`%s: <command>`)\n\n", name)
	for _, vs := range cmds {
		s += fmt.Sprintf("`%s` - %s\n", vs[0], vs[1])
	}
	return s
}
