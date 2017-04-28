package command

import (
	"fmt"
	"log"
	"strings"

	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
)

// Command is a function for a command
type Command func(q queue.Queue, channel string, user string, args string) (queue.Queue, Notification)

// Notification represents a message to a channel
type Notification struct {
	Channel string
	Message string
}

// QueueCommands provides the API to the various commands supported by the bot
type QueueCommands struct {
	id        string
	name      string
	response  responses
	userCache usercache.UserCache
}

// New returns a new Command instance
func New(id string, name string, uc usercache.UserCache) QueueCommands {
	r := responses{uc}
	c := QueueCommands{id, name, r, uc}
	return c
}

func (c QueueCommands) isFuzzyMatch(i queue.Item, id, reason string) (match bool) {
	match = i.ID == id && strings.HasPrefix(i.Reason, reason)
	return
}

func (c QueueCommands) findItemReverse(q queue.Queue, id, reason string) (item queue.Item, ok bool) {
	for ix := len(q) - 1; ix >= 0; ix-- {
		if c.isFuzzyMatch(q[ix], id, reason) {
			ok = true
			item = q[ix]
			break
		}
	}
	return
}

func (c QueueCommands) findItem(q queue.Queue, id, reason string) (item queue.Item, ok bool) {
	for _, i := range q {
		if c.isFuzzyMatch(i, id, reason) {
			ok = true
			item = i
			break
		}
	}
	return
}

func (c QueueCommands) getIDFromName(name string) (id string) {
	id = ""
	if strings.HasPrefix(name, "<@") {
		id = strings.Trim(name, "<@>")
	} else {
		id = c.userCache.GetUserID(name)
	}
	return
}

func (c QueueCommands) getNameIDPair(id string) (pair string) {
	name := c.userCache.GetUserName(id)
	return fmt.Sprintf("<%s|%s>", id, name)
}

func (c QueueCommands) logActivity(id, reason, text string) {
	log.Printf("%s (%s) %s", c.getNameIDPair(id), reason, text)
}

func cmdList(cmds [][]string) string {
	c := ""
	for _, vs := range cmds {
		c += fmt.Sprintf("`%s` - %s\n", vs[0], vs[1])
	}
	return c
}

// Help provides much needed assistance
func (c QueueCommands) Help(q queue.Queue, ch, id, args string) (queue.Queue, Notification) {
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
		[]string{"delegate <user>", "Delegate your place to someone else (your most recent entry is delegated)"},
		[]string{"delegate <user> <reason prefix>", "Delegate your place to someone else (match the entry with reason that starts with <reason prefix>)"},
	})

	s += "\n*If you need to get rid of somebody who is in the way:*\n"
	s += cmdList([][]string{
		[]string{"oust <name>", "Force the token holder to yield to the next in line"},
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
