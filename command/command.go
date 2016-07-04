package command

import (
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/notification"
	"strings"
	"fmt"
)

func findItem(q queue.Queue, name, reason string) queue.Item {
	var i queue.Item
	for ix := len(q) - 1; ix >= 0; ix-- {
		if q[ix].Name == name && strings.HasPrefix(q[ix].Reason, reason) {
			return q[ix]
		}
	}
	return i
}

// Join adds an item to the queue
func Join(q queue.Queue, name, reason string) (queue.Queue, string) {
	i := queue.Item{name, reason}

	if i.Reason == "" {
		return q, notification.JoinNoReason(i)
	}

	if q.Contains(i) {
		return q, ""
	}

	q = q.Add(i)
	if q.Active() == i {
		return q, notification.JoinActive(i)
	}

	return q, notification.Join(i)
}

// Leave removes an item from the queue
func Leave(q queue.Queue, name, reason string) (queue.Queue, string) {
	i := findItem(q, name, reason)
	if i.Name == "" {
		return q, ""
	}

	active := q.Active() == i
	if q.Contains(i) {
		q = q.Remove(i)
		if active {
			if len(q) == 0 {
				return q, notification.LeaveNoActive(i)
			}
			return q, notification.LeaveActive(i, q)
		}
		return q, notification.Leave(i)
	}
	return q, ""
}

// Done removes the active user from the queue
func Done(q queue.Queue, name string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, ""
	}

	i := q.Active()

	if i.Name != name {
		return q, notification.DoneNotActive(i)
	}

	q = q.Remove(i)
	if len(q) > 0 {
		return q, notification.Done(i, q)
	}
	return q, notification.DoneNoOthers(i)
}

// Yield allows the second place ahead of the active user
func Yield(q queue.Queue, name string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, notification.YieldNotActive(queue.Item{name, ""})
	}
	i := q.Active()
	if i.Name != name {
		return q, notification.YieldNotActive(queue.Item{name, ""})
	}
	if len(q) < 2 {
		return q, notification.YieldNoOthers(i)
	}
	q = q.Yield()
	return q, notification.Yield(i, q)
}

// Barge adds a user to the front of the queue
func Barge(q queue.Queue, name, reason string) (queue.Queue, string) {
	i := queue.Item{name, reason}
	q = q.Barge(i)
	if q.Active() == i {
		return q, notification.JoinActive(i)
	}
	return q, notification.Barge(i)
}

// Boot kicks someone from the waiting list
func Boot(q queue.Queue, booter, name, reason string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, ""
	}

	i := findItem(q, name, reason)
	if i.Name == "" {
		return q, ""
	}

	if q.Active() == i {
		return q, notification.OustNotBoot(booter)
	}

	if q.Contains(i) {
		q = q.Remove(i)
		return q, notification.Boot(booter, i)
	}
	return q, ""
}

// Oust boots the current token holder and gives it to the next person
func Oust(q queue.Queue, ouster, name, reason string) (queue.Queue, string) {
	if len(q) == 0 {
		return q, ""
	}

	i := findItem(q, name, reason)
	if i.Name == "" {
		return q, ""
	}

	if q.Active() != i {
		return q, notification.OustNotActive(ouster)
	}

	q = q.Remove(i)

	if len(q) == 0 {
		return q, notification.OustNoOthers(ouster, i)
	}
	return q, notification.Oust(ouster, i, q)
}

// List shows who has the token and who is waiting
func List(q queue.Queue) string {
	if len(q) == 0 {
		return "Nobody has the token, and nobody is waiting"
	}

	a := q.Active()
	s := fmt.Sprintf("%s (%s) has the token", a.Name, a.Reason)
	if len(q) == 1 {
		return fmt.Sprintf("%s, and nobody is waiting", s)
	}

	s += ", and waiting their turn are:"
	for ix, i := range q.Waiting() {
		s += fmt.Sprintf("\n%d: %s (%s)", ix + 1, i.Name, i.Reason)
	}
	return s
}

// Help provides much needed assistance
func Help(name string) string {
	cmds := [][]string{
		[]string{"join <reason>", "Join the queue and give a reason why"},
		[]string{"leave", "Leave the queue (your most recent entry is removed)"},
		[]string{"leave <reason>", "Leave the queue (your most recent entry starting with <reason> is removed)"},
		[]string{"done", "Release the token once you are done with it"},
		[]string{"yield", "Release the token and swap places with next in line"},
		[]string{"barge <reason>", "Barge to the front of the queue so you get the token next (only with good reason!)"},
		[]string{"boot <name>", "Kick somebody out of the waiting list (their most recent entry is removed)"},
		[]string{"boot <name> <reason>", "Kick somebody out of the waiting list (their most recent entry starting with <reason> is removed"},
		[]string{"oust", "Forcibly take the token from the token holder and kick them out of the queue (only with VERY good reason!)"},
		[]string{"list", "Show who has the token and who is waiting"},
		[]string{"help", "Show this text"},
	}
	s := ""
	for _, vs := range cmds {
		s += fmt.Sprintf("`%s: %s` - %s\n", name, vs[0], vs[1])
	}
	return s
}
