package command

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
	"github.com/doozr/qbot/util"
)

// PendingOust contains an oust request that must be fulfilled
type PendingOust struct {
	Item      queue.Item
	Timestamp time.Time
}

// CmdFn is a function for a command
type CmdFn func(q queue.Queue, channel string, user string, args string) (queue.Queue, Notification)

// Notification represents a message to a channel
type Notification struct {
	Channel string
	Message string
}

// Command provides the API to the various commands supported by the bot
type Command struct {
	name         string
	response     responses
	userCache    usercache.UserCache
	pendingOusts map[string]PendingOust
}

// New returns a new Command instance
func New(name string, uc usercache.UserCache) Command {
	r := responses{uc}
	c := Command{name, r, uc, make(map[string]PendingOust)}
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

func (c Command) getNameIDPair(id string) (pair string) {
	name := c.userCache.GetUserName(id)
	return fmt.Sprintf("<%s|%s>", id, name)
}

func (c Command) logActivity(id, reason, text string) {
	log.Printf("%s (%s) %s", c.getNameIDPair(id), reason, text)
}

// Done removes the active user from the queue
func (c Command) Done(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, ""}
	}

	i := q.Active()

	if i.ID != id {
		return q, Notification{ch, c.response.DoneNotActive(id)}
	}

	q = q.Remove(i)
	c.logActivity(id, i.Reason, "done")
	if len(q) > 0 {
		n := q.Active()
		c.logActivity(n.ID, n.Reason, "is active")
		return q, Notification{ch, c.response.Done(i, q)}
	}
	return q, Notification{ch, c.response.DoneNoOthers(i)}
}

// Yield allows the second place ahead of the active user
func (c Command) Yield(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, c.response.YieldNotActive(queue.Item{ID: id, Reason: ""})}
	}
	i := q.Active()
	if i.ID != id {
		return q, Notification{ch, c.response.YieldNotActive(queue.Item{ID: id, Reason: ""})}
	}
	if len(q) < 2 {
		return q, Notification{ch, c.response.YieldNoOthers(i)}
	}
	q = q.Yield()
	n := q.Active()
	c.logActivity(id, i.Reason, "yielded")
	c.logActivity(n.ID, n.Reason, "is active")
	return q, Notification{ch, c.response.Yield(i, q)}
}

// Barge adds a user to the front of the queue
func (c Command) Barge(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	i := queue.Item{ID: id, Reason: args}
	q = q.Barge(i)
	if q.Active() == i {
		return q, Notification{ch, c.response.JoinActive(i)}
	}
	c.logActivity(id, args, "barged")
	return q, Notification{ch, c.response.Barge(i, q.Active())}
}

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

	if q.Contains(i) {
		q = q.Remove(i)
		c.logActivity(id, reason, "booted by "+c.getNameIDPair(booter))
		return q, Notification{ch, c.response.Boot(booter, i)}
	}

	return q, Notification{ch, ""}
}

// Oust boots the current token holder and gives it to the next person
func (c Command) Oust(q queue.Queue, ch, ouster, args string) (queue.Queue, Notification) {
	if len(q) == 0 {
		return q, Notification{ch, ""}
	}

	id := c.getIDFromName(args)
	if id == "" {
		return q, Notification{ch, c.response.OustNotActive(ouster)}
	}

	i := q.Active()
	if i.ID != id {
		return q, Notification{ch, c.response.OustNotActive(ouster)}
	}

	// If a previous request has been lodged in the last 30 seconds
	// and all is well then oust the active user
	pendingOust, ok := c.pendingOusts[ouster]
	if ok && pendingOust.Item == i {
		if time.Since(pendingOust.Timestamp).Seconds() < 30 && q.Active() == i {
			q = q.Remove(i)
			c.logActivity(i.ID, i.Reason, "ousted by "+c.getNameIDPair(ouster))
			if len(q) == 0 {
				return q, Notification{ch, c.response.OustNoOthers(ouster, i)}
			}
			n := q.Active()
			c.logActivity(n.ID, n.Reason, "is active")
			return q, Notification{ch, c.response.Oust(ouster, i, q)}
		}
	}

	c.pendingOusts[ouster] = PendingOust{i, time.Now()}

	return q, Notification{ch, c.response.OustConfirm(ouster, i)}
}

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

func cmdList(cmds [][]string) string {
	c := ""
	for _, vs := range cmds {
		c += fmt.Sprintf("`%s` - %s\n", vs[0], vs[1])
	}
	return c
}

// Help provides brief assistance
func (c Command) Help(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	s := fmt.Sprintf("Address each command to the bot (`%s: <command>`)\n\n", c.name)

	s += cmdList([][]string{
		[]string{"list", "Show who has the token and who is waiting"},
		[]string{"join <reason>", "Join the queue and give a reason why"},
		[]string{"done", "Release the token once you are done with it"},
		[]string{"drop", "Drop the token and leave the queue (note: actually just an alias of `done`)"},
		[]string{"yield", "Relinquish the token and swap places with the next in line"},
		[]string{"leave <reason>", "Leave the queue (your most recent entry starting with <reason> is removed)"},
		[]string{"help", "Show this text"},
		[]string{"morehelp", "Show more detailed help and extra actions"},
	})
	return q, Notification{id, s}
}

// MoreHelp provides much needed assistance
func (c Command) MoreHelp(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
	s := fmt.Sprintf("Address each command to the bot (`%s: <command>`)\n\n", c.name)

	s += "*If you don't have the token and need it:*\n"
	s += cmdList([][]string{
		[]string{"join <reason>", "Join the queue and give a reason why"},
		[]string{"barge <reason>", "Barge to the front of the queue so you get the token next (only with good reason!)"},
	})

	s += "\n*If you have the token and have done with it:*\n"
	s += cmdList([][]string{
		[]string{"done", "Release the token once you are done with it"},
		[]string{"drop", "Drop the token and leave the queue (note: actually just an alias of `done`)"},
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
	return q, Notification{id, s}
}
