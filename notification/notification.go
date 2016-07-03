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

func yielded(i queue.Item) string {
	return fmt.Sprintf("@%s (%s) has yielded the token", i.Name, i.Reason)
}

func ousted(ouster string, i queue.Item) string {
	return fmt.Sprintf("@%s ousted @%s (%s)", ouster, i.Name, i.Reason)
}

func Join(i queue.Item) string {
	return fmt.Sprintf("@here @%s (%s) has joined the queue", i.Name, i.Reason)
}

func JoinActive(i queue.Item) string {
	return fmt.Sprintf("@here %s", nowHasToken(i))
}

func Leave(i queue.Item) string {
	return fmt.Sprintf("@here @%s (%s) has left the queue", i.Name, i.Reason)
}

func LeaveActive(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("@here %s\n\n%s", Leave(i), nowHasToken(a))
}

func LeaveNoActive(i queue.Item) string {
	return fmt.Sprintf("@here %s\n\n%s", Leave(i), upForGrabs())
}

func Done(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("@here %s\n\n%s",
		finishedWithToken(i), nowHasToken(a))
}

func DoneNoOthers(i queue.Item) string {
	return fmt.Sprintf("@here %s\n\nThe token is up for grabs%s",
		finishedWithToken(i), upForGrabs())
}

func DoneNotActive(i queue.Item) string {
	return fmt.Sprintf("@%s You cannot be done if you don't have the token")
}

func Yield(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("@here %s\n\n%s", yielded(i), nowHasToken(a))
}

func YieldNoOthers(i queue.Item) string {
	return fmt.Sprintf("@%s You cannot yield if there is nobody waiting", i.Name)
}

func YieldNotActive(i queue.Item) string {
	return fmt.Sprintf("@%s You cannot yield if you do not have the token", i.Name)
}

func Barge(i queue.Item) string {
	return fmt.Sprintf("@here @%s (%s) barged to the front", i.Name, i.Reason)
}

func Boot(booter string, i queue.Item) string {
	return fmt.Sprintf("@here @%s booted %s (%s) from the list", booter, i.Name, i.Reason)
}

func OustNotBoot(booter string) string {
	return fmt.Sprintf("@%s You must oust the token holder", booter)
}

func Oust(ouster string, i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("@here %s\n\n%s", ousted(ouster, i), nowHasToken(a))
}

func OustNotActive(ouster string) string {
	return fmt.Sprintf("@%s You can only oust the token holder", ouster)
}

func OustNoOthers(ouster string, i queue.Item) string {
	return fmt.Sprintf("@here %s", ousted(ouster, i))
}