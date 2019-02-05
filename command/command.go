package command

import (
	"fmt"
	"log"
	"strconv"
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

func (c QueueCommands) findItem(q queue.Queue, id string) (item queue.Item, ok bool) {
	for _, i := range q {
		if i.ID == id {
			ok = true
			item = i
			break
		}
	}
	return
}

func (c QueueCommands) findItemReverse(q queue.Queue, id string) (item queue.Item, ok bool) {
	for ix := len(q) - 1; ix >= 0; ix-- {
		if q[ix].ID == id {
			ok = true
			item = q[ix]
			return
		}
	}
	return
}

func (c QueueCommands) findByPosition(q queue.Queue, position int) (item queue.Item, ok bool) {
	if position < 1 || position > len(q) {
		return
	}

	return q[position-1], true
}

func (c QueueCommands) parsePosition(args string) (position int, remainder string, ok bool) {
	remainder = args
	fields := strings.Fields(args)

	// No fields? Quit now
	if len(fields) == 0 {
		return
	}

	// First field is not an integer? Again, stop
	position, err := strconv.Atoi(fields[0])
	if err != nil {
		return
	}

	// Anything after this point is OK
	ok = true

	// If the integer is the only field, return an empty string
	if len(fields) == 1 {
		remainder = ""
		return
	}

	// String the position part from the start of the string
	positionLength := len(fields[0])
	remainder = strings.Trim(args[positionLength:], " ")
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
		{"join <reason>", "Join the queue and give a reason why"},
		{"barge <reason>", "Barge to the front of the queue so you get the token next (only with good reason!)"},
		{"barge <position>", "Barge the entry at the given position to the front of the queue"},
	})

	s += "\n*If you have the token and have done with it:*\n"
	s += cmdList([][]string{
		{"done", "Release the token once you are done with it"},
		{"drop", "Drop the token and leave the queue (note: actually just an alias of `done`)"},
		{"yield", "Release the token and swap places with next in line"},
	})

	s += "\n*If you are in the queue and need to change something:*\n"
	s += cmdList([][]string{
		{"delegate <user>", "Delegate your place to someone else (your most recent entry is delegated)"},
		{"delegate <user> <reason prefix>", "Delegate your place to someone else (match the entry with reason that starts with <reason prefix>)"},
		{"replace <position> <reason>", "Replace the reason of a queue entry you own"},
	})

	s += "\n*If you are in the queue and need to leave:*\n"
	s += cmdList([][]string{
		{"leave", "Leave the queue (your most recent entry is removed)"},
		{"leave <position>", "Leave the queue (match the entry at the given position)"},
	})

	s += "\n*If you need to get rid of somebody who is in the way:*\n"
	s += cmdList([][]string{
		{"oust <name>", "Force the token holder to yield to the next in line"},
		{"boot <name>", "Kick somebody out of the waiting list (their most recent entry is removed)"},
		{"boot <position> <name>", "Kick somebody out of the waiting list (match the entry at the given position)"},
	})

	s += "\n*Notifications (from automated systems):*\n"
	s += cmdList([][]string{
		{"success", "Notify the token holder and next in line of success and remove the token holder"},
		{"failure <message>", "Notify the token holder and next in line of failure with a custom error message"},
	})

	s += "\n*Other useful things to know:*\n"
	s += cmdList([][]string{
		{"list", "Show who has the token and who is waiting"},
		{"help", "Show this text"},
	})
	return q, Notification{id, s}
}
