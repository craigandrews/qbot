package notification

import (
	"fmt"

	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
)

type Notification struct {
	UserCache *usercache.UserCache
}

func New(userCache *usercache.UserCache) Notification {
	n := Notification{ userCache }
	return n
}

func (n Notification) GetUserName(id string) (username string) {
	username = n.UserCache.GetUserName(id)
	return
}

func (n Notification) GetUserId(name string) (id string) {
	id = n.UserCache.GetUserId(name)
	return
}

func (n Notification) link(i string) string {
	return fmt.Sprintf("<@%s|%s>", i, n.GetUserName(i))
}

func (n Notification) item(i queue.Item) string {
	if i.Reason == "" {
		return n.link(i.Id)
	}
	return fmt.Sprintf("%s (%s)", n.link(i.Id), i.Reason)
}

func (n Notification) finishedWithToken(i queue.Item) string {
	return fmt.Sprintf("%s has finished with the token", n.item(i))
}

func (n Notification) nowHasToken(i queue.Item) string {
	return fmt.Sprintf("*%s now has the token*", n.item(i))
}

func (n Notification) upForGrabs() string {
	return fmt.Sprint("The token is up for grabs")
}

func (n Notification) yielded(i queue.Item) string {
	return fmt.Sprintf("%s has yielded the token", n.item(i))
}

func (n Notification) ousted(ouster string, i queue.Item) string {
	return fmt.Sprintf("%s ousted %s", n.link(ouster), n.item(i))
}

func (n Notification) Join(i queue.Item) string {
	return fmt.Sprintf("%s has joined the queue", n.item(i))
}

func (n Notification) JoinNoReason(i queue.Item) string {
	return fmt.Sprintf("%s You must provide a reason for joining", n.link(i.Id))
}

func (n Notification) JoinActive(i queue.Item) string {
	return fmt.Sprintf("%s", n.nowHasToken(i))
}

func (n Notification) Leave(i queue.Item) string {
	return fmt.Sprintf("%s has left the queue", n.item(i))
}

func (n Notification) LeaveActive(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n%s", n.Leave(i), n.nowHasToken(a))
}

func (n Notification) LeaveNoActive(i queue.Item) string {
	return fmt.Sprintf("%s\n%s", n.Leave(i), n.upForGrabs())
}

func (n Notification) Done(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n%s",
		n.finishedWithToken(i), n.nowHasToken(a))
}

func (n Notification) DoneNoOthers(i queue.Item) string {
	return fmt.Sprintf("%s\n%s",
		n.finishedWithToken(i), n.upForGrabs())
}

func (n Notification) DoneNotActive(i queue.Item) string {
	return fmt.Sprintf("%s You cannot be done if you don't have the token", n.link(i.Id))
}

func (n Notification) Yield(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n%s", n.yielded(i), n.nowHasToken(a))
}

func (n Notification) YieldNoOthers(i queue.Item) string {
	return fmt.Sprintf("%s You cannot yield if there is nobody waiting", n.link(i.Id))
}

func (n Notification) YieldNotActive(i queue.Item) string {
	return fmt.Sprintf("%s You cannot yield if you do not have the token", n.link(i.Id))
}

func (n Notification) Barge(i queue.Item) string {
	return fmt.Sprintf("%s barged to the front", n.item(i))
}

func (n Notification) Boot(booter string, i queue.Item) string {
	return fmt.Sprintf("%s booted %s from the list", n.link(booter), n.item(i))
}

func (n Notification) OustNotBoot(booter string) string {
	return fmt.Sprintf("%s You must oust the token holder", n.link(booter))
}

func (n Notification) Oust(ouster string, i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n%s", n.ousted(ouster, i), n.nowHasToken(a))
}

func (n Notification) OustNotActive(ouster string) string {
	return fmt.Sprintf("%s You can only oust the token holder", n.link(ouster))
}

func (n Notification) OustNoOthers(ouster string, i queue.Item) string {
	return n.ousted(ouster, i)
}
