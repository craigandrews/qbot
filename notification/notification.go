package notification

import (
	"fmt"

	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
)

// Notification contains the notification API
type Notification struct {
	UserCache *usercache.UserCache
}

// New returns a new Notification instance
func New(userCache *usercache.UserCache) Notification {
	n := Notification{userCache}
	return n
}

func (n Notification) getUserName(id string) (username string) {
	username = n.UserCache.GetUserName(id)
	return
}

func (n Notification) getUserID(name string) (id string) {
	id = n.UserCache.GetUserID(name)
	return
}

func (n Notification) link(i string) string {
	return fmt.Sprintf("<@%s|%s>", i, n.getUserName(i))
}

func (n Notification) item(i queue.Item) string {
	if i.Reason == "" {
		return n.link(i.ID)
	}
	return fmt.Sprintf("%s (%s)", n.link(i.ID), i.Reason)
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

// Join is a successful join to the queue
func (n Notification) Join(i queue.Item) string {
	return fmt.Sprintf("%s has joined the queue", n.item(i))
}

// JoinNoReason tells the user that a reason is required on join
func (n Notification) JoinNoReason(i queue.Item) string {
	return fmt.Sprintf("%s You must provide a reason for joining", n.link(i.ID))
}

// JoinActive tells the user that they have immediately taken the token on join
func (n Notification) JoinActive(i queue.Item) string {
	return fmt.Sprintf("%s", n.nowHasToken(i))
}

// Leave is a successful leave from the queue
func (n Notification) Leave(i queue.Item) string {
	return fmt.Sprintf("%s has left the queue", n.item(i))
}

// LeaveActive tells the user they should use done rather than leave if they have the token
func (n Notification) LeaveActive(i queue.Item) string {
	return fmt.Sprintf("%s You have the token, did you mean 'done'?", n.link(i.ID))
}

// LeaveNoEntry tells the user that an entry with the requested reason does not exist
func (n Notification) LeaveNoEntry(id, reason string) string {
	return fmt.Sprintf("%s No entry with a reason that starts with '%s' was found", n.link(id), reason)
}

// Done is a successful drop of the token
func (n Notification) Done(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n%s",
		n.finishedWithToken(i), n.nowHasToken(a))
}

// DoneNoOthers is a successful drop of the token when nobody can pick it up
func (n Notification) DoneNoOthers(i queue.Item) string {
	return fmt.Sprintf("%s\n%s",
		n.finishedWithToken(i), n.upForGrabs())
}

// DoneNotActive tells the user that they must have the token to use done
func (n Notification) DoneNotActive(user string) string {
	return fmt.Sprintf("%s You cannot be done if you don't have the token", n.link(user))
}

// Yield is a successful passing of the token to next in line
func (n Notification) Yield(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n%s", n.yielded(i), n.nowHasToken(a))
}

// YieldNoOthers tells the user that they cannot yield if nobody is waiting
func (n Notification) YieldNoOthers(i queue.Item) string {
	return fmt.Sprintf("%s You cannot yield if there is nobody waiting", n.link(i.ID))
}

// YieldNotActive tells the user they that must have the token to yield
func (n Notification) YieldNotActive(i queue.Item) string {
	return fmt.Sprintf("%s You cannot yield if you do not have the token", n.link(i.ID))
}

// Barge is a successful barge to the front of the queue
func (n Notification) Barge(i queue.Item, a queue.Item) string {
	return fmt.Sprintf("%s barged to the front\n%s still has the token", n.item(i), n.item(a))
}

// Boot is a successful force remove from the queue
func (n Notification) Boot(booter string, i queue.Item) string {
	return fmt.Sprintf("%s booted %s from the list", n.link(booter), n.item(i))
}

// BootNoEntry tells the user that an entry with the requested user and reason does not exist
func (n Notification) BootNoEntry(id, name, reason string) string {
	return fmt.Sprintf("%s No entry with for %s with a reason that starts with '%s' was found", n.link(id), name, reason)
}

// OustNotBoot tells the user that they can't boot the token holder
func (n Notification) OustNotBoot(booter string) string {
	return fmt.Sprintf("%s You must oust the token holder", n.link(booter))
}

// Oust is a successful oust
func (n Notification) Oust(ouster string, i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n%s", n.ousted(ouster, i), n.nowHasToken(a))
}

// OustNotActive tells the user they can only oust the token holder
func (n Notification) OustNotActive(ouster string) string {
	return fmt.Sprintf("%s You can only oust the token holder", n.link(ouster))
}

// OustNoOthers is a successful oust when nobody can pick up the token
func (n Notification) OustNoOthers(ouster string, i queue.Item) string {
	return n.ousted(ouster, i)
}

// OustConfirm asks the ouster if they are sure
func (n Notification) OustConfirm(ouster string, i queue.Item) string {
	return fmt.Sprintf("%s Are you sure you want to oust %s?\n(Repeat this command within 30 seconds to confirm)",
		n.link(ouster), n.link(i.ID))
}
