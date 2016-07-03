package notification

import (
	"github.com/doozr/qbot/queue"
	"fmt"
)

func finishedWithToken(i queue.Item) string {
	return fmt.Sprintf("@%s (%s) has finished with the token", i.Name, i.Reason)
}

func nowHasToken(i queue.Item) string {
	return fmt.Sprintf("@%s (%s) now has the token", i.Name, i.Reason)
}

func upForGrabs() string {
	return fmt.Sprint("The token is up for grabs")
}

func Join(i queue.Item) string {
	return fmt.Sprintf("@%s (%s) has joined the queue", i.Name, i.Reason)
}

func JoinActive(i queue.Item) string {
	return fmt.Sprintf("%s", nowHasToken(i))
}

func Leave(i queue.Item) string {
	return fmt.Sprintf("@%s (%s) has left the queue", i.Name, i.Reason)
}

func LeaveActive(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n\n%s", Leave(i), nowHasToken(a))
}

func LeaveNoActive(i queue.Item) string {
	return fmt.Sprintf("%s\n\n%s", Leave(i), upForGrabs())
}

func Done(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n\n%s",
		finishedWithToken(i), nowHasToken(a))
}

func DoneNoOthers(i queue.Item) string {
	return fmt.Sprintf("%s\n\nThe token is up for grabs%s",
		finishedWithToken(i), upForGrabs())
}

func DoneNotActive(i queue.Item) string {
	return fmt.Sprintf("@%s You cannot be done if you don't have the token")
}