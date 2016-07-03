package notification

import (
	"github.com/doozr/qbot/queue"
	"fmt"
)

func Join(i queue.Item) string {
	return fmt.Sprintf("@%s (%s) has joined the queue", i.Name, i.Reason)
}

func Active(i queue.Item) string {
	return fmt.Sprintf("@%s (%s) now has the token", i.Name, i.Reason)
}

func Leave(i queue.Item) string {
	return fmt.Sprintf("@%s (%s) has left the queue", i.Name, i.Reason)
}

func LeaveActive(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n\n@%s (%s) now has the token", Leave(i), a.Name, a.Reason)
}

func LeaveNoActive(i queue.Item) string {
	return fmt.Sprintf("%s\n\nThe token is up for grabs", Leave(i))
}

func Done(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("@%s (%s) has finished with the token\n\n@%s (%s) now has the token",
		i.Name, i.Reason, a.Name, a.Reason)
}

func DoneNoOthers(i queue.Item) string {
	return fmt.Sprintf("@%s (%s) has finished with the token\n\nThe token is up for grabs",
		i.Name, i.Reason)
}

func DoneNotActive(i queue.Item) string {
	return fmt.Sprintf("@%s You cannot be done if you don't have the token")
}