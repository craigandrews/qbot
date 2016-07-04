package notification

import (
	"github.com/doozr/qbot/queue"
	"fmt"
)

func item(i queue.Item) string {
	if i.Reason == "" {
		return fmt.Sprintf("<@%s>", i.Name)
	}
	return fmt.Sprintf("<@%s> (%s)", i.Name, i.Reason)
}

func finishedWithToken(i queue.Item) string {
	return fmt.Sprintf("%s has finished with the token", item(i))
}

func nowHasToken(i queue.Item) string {
	return fmt.Sprintf("%s now has the token", item(i))
}

func upForGrabs() string {
	return fmt.Sprint("The token is up for grabs")
}

func yielded(i queue.Item) string {
	return fmt.Sprintf("%s has yielded the token", item(i))
}

func ousted(ouster string, i queue.Item) string {
	return fmt.Sprintf("@%s ousted %s", ouster, item(i))
}

func Join(i queue.Item) string {
	return fmt.Sprintf("%s has joined the queue", item(i))
}

func JoinNoReason(i queue.Item) string {
	return fmt.Sprintf("@%s You must provide a reason for joining", i.Name)
}

func JoinActive(i queue.Item) string {
	return fmt.Sprintf("%s", nowHasToken(i))
}

func Leave(i queue.Item) string {
	return fmt.Sprintf("%s has left the queue", item(i))
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
	return fmt.Sprintf("%s\n\n%s",
		finishedWithToken(i), upForGrabs())
}

func DoneNotActive(i queue.Item) string {
	return fmt.Sprintf("@%s You cannot be done if you don't have the token", i.Name)
}

func Yield(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n\n%s", yielded(i), nowHasToken(a))
}

func YieldNoOthers(i queue.Item) string {
	return fmt.Sprintf("@%s You cannot yield if there is nobody waiting", i.Name)
}

func YieldNotActive(i queue.Item) string {
	return fmt.Sprintf("@%s You cannot yield if you do not have the token", i.Name)
}

func Barge(i queue.Item) string {
	return fmt.Sprintf("%s barged to the front", item(i))
}

func Boot(booter string, i queue.Item) string {
	return fmt.Sprintf("@%s booted %s from the list", booter, item(i))
}

func OustNotBoot(booter string) string {
	return fmt.Sprintf("@%s You must oust the token holder", booter)
}

func Oust(ouster string, i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n\n%s", ousted(ouster, i), nowHasToken(a))
}

func OustNotActive(ouster string) string {
	return fmt.Sprintf("@%s You can only oust the token holder", ouster)
}

func OustNoOthers(ouster string, i queue.Item) string {
	return ousted(ouster, i)
}
