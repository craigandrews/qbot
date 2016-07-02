package command

import (
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/notification"
)

// Join adds an item to the queue
func Join(q queue.Queue, i queue.Item) (queue.Queue, string){
	if q.Contains(i) {
		return q, ""
	}

	q = q.Add(i)
	if q.Active() == i {
		return q, notification.Active(i)
	}

	return q, notification.Join(i)
}

// Leave
// Done
// Yield
// Barge
// Boot
// Oust
// List
// Help
