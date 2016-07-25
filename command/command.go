package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/doozr/qbot/notification"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
)

// PendingOust contains an oust request that must be fulfilled
type PendingOust struct {
	Item      queue.Item
	Timestamp time.Time
}

// Command provides the API to the various commands supported by the bot
type Command struct {
	notification notification.Notification
	userCache    *usercache.UserCache
	pendingOusts map[string]PendingOust
}

// New returns a new Command instance
func New(n notification.Notification, uc *usercache.UserCache) Command {
	c := Command{n, uc, make(map[string]PendingOust)}
	return c
}

func (c Command) findItem(q queue.Queue, id, reason string) (item queue.Item, ok bool) {
	for ix := len(q) - 1; ix >= 0; ix-- {
		if q[ix].ID == id && strings.HasPrefix(q[ix].Reason, reason) {
			ok = true
			item = q[ix]
			break
		}
	}
	return
}

func (c Command) getIDFromName(name string) (id string) {
	id = ""
	if strings.HasPrefix(name, "<@") {
		id = strings.Trim(name, "<@>")
	} else {
		id = c.userCache.GetUserID(name)
	}
	return
}

// Join adds an item to the queue
func (c Command) Join(q queue.Queue, id, reason string) (queue.Queue, string) {
	i := queue.Item{ID: id, Reason: reason}

	if i.Reason == "" {
		return q, c.notification.JoinNoReason(i)
	}

	if q.Contains(i) {
		return q, ""
	}

	q = q.Add(i)
	if q.Active() == i {
		return q, c.notification.JoinActive(i)
	}

	return q, c.notification.Join(i)
}

// Leave removes an item from the queue
func (c Command) Leave(q queue.Queue, id, reason string) (queue.Queue, string) {
	i, ok := c.findItem(q, id, reason)
	if !ok {
		return q, c.notification.LeaveNoEntry(id, reason)
	}

	if q.Active() == i {
		return q, c.notification.LeaveActive(i)
	}

	if q.Contains(i) {
		q = q.Remove(i)
		return q, c.notification.Leave(i)
	}

	return q, ""
}

// Done removes the active user from the queue
func (c Command) Done(q queue.Queue, id string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, ""
	}

	i := q.Active()

	if i.ID != id {
		return q, c.notification.DoneNotActive(i)
	}

	q = q.Remove(i)
	if len(q) > 0 {
		return q, c.notification.Done(i, q)
	}
	return q, c.notification.DoneNoOthers(i)
}

// Yield allows the second place ahead of the active user
func (c Command) Yield(q queue.Queue, id string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, c.notification.YieldNotActive(queue.Item{ID: id, Reason: ""})
	}
	i := q.Active()
	if i.ID != id {
		return q, c.notification.YieldNotActive(queue.Item{ID: id, Reason: ""})
	}
	if len(q) < 2 {
		return q, c.notification.YieldNoOthers(i)
	}
	q = q.Yield()
	return q, c.notification.Yield(i, q)
}

// Barge adds a user to the front of the queue
func (c Command) Barge(q queue.Queue, id, reason string) (queue.Queue, string) {
	i := queue.Item{ID: id, Reason: reason}
	q = q.Barge(i)
	if q.Active() == i {
		return q, c.notification.JoinActive(i)
	}
	return q, c.notification.Barge(i)
}

// Boot kicks someone from the waiting list
func (c Command) Boot(q queue.Queue, booter, name, reason string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, ""
	}

	id := c.getIDFromName(name)
	i, ok := c.findItem(q, id, reason)
	if !ok {
		return q, c.notification.BootNoEntry(booter, name, reason)
	}

	if q.Active() == i {
		return q, c.notification.OustNotBoot(booter)
	}

	if q.Contains(i) {
		q = q.Remove(i)
		return q, c.notification.Boot(booter, i)
	}

	return q, ""
}

// Oust boots the current token holder and gives it to the next person
func (c Command) Oust(q queue.Queue, ouster, name string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, ""
	}

	id := c.getIDFromName(name)
	if id == "" {
		return q, c.notification.OustNotActive(ouster)
	}

	i := q.Active()
	if i.ID != id {
		return q, c.notification.OustNotActive(ouster)
	}

	// If a previous request has been lodged in the last 30 seconds
	// and all is well then oust the active user
	pendingOust, ok := c.pendingOusts[ouster]
	if ok && pendingOust.Item == i {
		if time.Since(pendingOust.Timestamp).Seconds() < 30 && q.Active() == i {
			q = q.Remove(i)

			if len(q) == 0 {
				return q, c.notification.OustNoOthers(ouster, i)
			}
			return q, c.notification.Oust(ouster, i, q)
		}
	}

	c.pendingOusts[ouster] = PendingOust{i, time.Now()}

	return q, c.notification.OustConfirm(ouster, i)
}

// List shows who has the token and who is waiting
func (c Command) List(q queue.Queue) string {
	if len(q) == 0 {
		return "Nobody has the token, and nobody is waiting"
	}

	a := q.Active()
	s := fmt.Sprintf("*%d: %s (%s) has the token*", 1, c.userCache.GetUserName(a.ID), a.Reason)
	for ix, i := range q.Waiting() {
		s += fmt.Sprintf("\n%d: %s (%s)", ix+2, c.userCache.GetUserName(i.ID), i.Reason)
	}
	return s
}

func cmdList(cmds [][]string) string {
	c := ""
	for _, vs := range cmds {
		c += fmt.Sprintf("`%s` - %s\n", vs[0], vs[1])
	}
	return c
}

// Help provides brief assistance
func (c Command) Help(name string) string {
	s := fmt.Sprintf("Address each command to the bot (`%s: <command>`)\n\n", name)

	s += cmdList([][]string{
		[]string{"list", "Show who has the token and who is waiting"},
		[]string{"join <reason>", "Join the queue and give a reason why"},
		[]string{"done", "Release the token once you are done with it"},
		[]string{"yield", "Relinquish the token and swap places with the next in line"},
		[]string{"leave <reason>", "Leave the queue (your most recent entry starting with <reason> is removed)"},
		[]string{"help", "Show this text"},
		[]string{"morehelp", "Show more detailed help and extra actions"},
	})
	return s
}

// MoreHelp provides much needed assistance
func (c Command) MoreHelp(name string) string {
	s := fmt.Sprintf("Address each command to the bot (`%s: <command>`)\n\n", name)

	s += "*If you don't have the token and need it:*\n"
	s += cmdList([][]string{
		[]string{"join <reason>", "Join the queue and give a reason why"},
		[]string{"barge <reason>", "Barge to the front of the queue so you get the token next (only with good reason!)"},
	})

	s += "\n*If you have the token and have done with it:*\n"
	s += cmdList([][]string{
		[]string{"done", "Release the token once you are done with it"},
		[]string{"yield", "Release the token and swap places with next in line"},
	})

	s += "\n*If you are in the queue and need to leave:*\n"
	s += cmdList([][]string{
		[]string{"leave", "Leave the queue (your most recent entry is removed)"},
		[]string{"leave <reason prefix>", "Leave the queue (match the entry with reason that starts with <reason prefix>)"},
	})

	s += "\n*If you need to get rid of somebody who is in the way:*\n"
	s += cmdList([][]string{
		[]string{"oust <name>", "Forcibly take the token from the token holder and kick them out of the queue (only with VERY good reason!)"},
		[]string{"boot <name>", "Kick somebody out of the waiting list (their most recent entry is removed)"},
		[]string{"boot <name> <reason prefix>", "Kick somebody out of the waiting list (match the entry with reason that starts with <reason prefix>)"},
	})

	s += "\n*Other useful things to know:*\n"
	s += cmdList([][]string{
		[]string{"list", "Show who has the token and who is waiting"},
		[]string{"help", "Show this text"},
	})
	return s
}
