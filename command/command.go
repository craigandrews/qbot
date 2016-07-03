package command

import (
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/notification"
	"strings"
)

// Join adds an item to the queue
func Join(q queue.Queue, name, reason string) (queue.Queue, string) {
	i := queue.Item{name, reason}
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
	var i queue.Item
	for ix := len(q) - 1; ix >= 0; ix-- {
		if q[ix].Name == name && strings.HasPrefix(q[ix].Reason, reason) {
			i = q[ix]
			break
		}
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

// Done
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

// Yield
// Barge
// Boot
// Oust
// List
// Help
