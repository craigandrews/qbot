package command

import (
	"fmt"

	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
	"github.com/doozr/qbot/util"
)

// responses builds mad-libbed reply strings
type responses struct {
	UserCache usercache.UserCache
}

func (n responses) getUserName(id string) (username string) {
	username = n.UserCache.GetUserName(id)
	return
}

func (n responses) link(i string) string {
	return fmt.Sprintf("<@%s|%s>", i, n.getUserName(i))
}

func (n responses) item(i queue.Item) string {
	return fmt.Sprintf("%s (%s)", n.link(i.ID), i.Reason)
}

func (n responses) finishedWithToken(i queue.Item) string {
	return fmt.Sprintf("%s has finished with the token", n.item(i))
}

func (n responses) nowHasToken(i queue.Item) string {
	return fmt.Sprintf("*%s now has the token*", n.item(i))
}

func (n responses) upForGrabs() string {
	return fmt.Sprint("The token is up for grabs")
}

func (n responses) yielded(i queue.Item) string {
	return fmt.Sprintf("%s has yielded the token", n.item(i))
}

func (n responses) ousted(ouster string, i queue.Item) string {
	return fmt.Sprintf("%s ousted %s", n.link(ouster), n.item(i))
}

// BadIndex is an attempt to manipulate the queue with an out of range index
func (n responses) BadIndex(i queue.Item, position int) string {
	return fmt.Sprintf("%s That's not a valid position in the queue", n.link(i.ID))
}

// NotOwned is an attempt to change a queue entry that the user does not own
func (n responses) NotOwned(i queue.Item, position int, o queue.Item) string {
	suffix := util.Suffix(position)
	ordinal := fmt.Sprintf("%d%s", position, suffix)
	return fmt.Sprintf("%s Not replacing because %s is %s in line", n.link(i.ID), n.link(o.ID), ordinal)
}

// Join is a successful join to the queue
func (n responses) Join(i queue.Item, position int) string {
	ordinal := "next"
	if position > 2 {
		suffix := util.Suffix(position)
		ordinal = fmt.Sprintf("%d%s", position, suffix)
	}
	return fmt.Sprintf("%s is now %s in line", n.item(i), ordinal)
}

// ReplaceNoReason tells the user that a reason is required when replacing
func (n responses) ReplaceNoReason(i queue.Item) string {
	return fmt.Sprintf("%s You must provide a new reason", n.link(i.ID))
}

// JoinNoReason tells the user that a reason is required on join
func (n responses) JoinNoReason(i queue.Item) string {
	return fmt.Sprintf("%s You must provide a reason for joining", n.link(i.ID))
}

// JoinActive tells the user that they have immediately taken the token on join
func (n responses) JoinActive(i queue.Item) string {
	return fmt.Sprintf("%s", n.nowHasToken(i))
}

// Leave is a successful leave from the queue
func (n responses) Leave(i queue.Item) string {
	return fmt.Sprintf("%s has left the queue", n.item(i))
}

// LeaveActive tells the user they should use done rather than leave if they have the token
func (n responses) LeaveActive(i queue.Item) string {
	return fmt.Sprintf("%s You have the token, did you mean `done` or `drop`?", n.link(i.ID))
}

// LeaveNoEntry tells the user that an entry with the requested reason does not exist
func (n responses) LeaveNoEntry(id, reason string) string {
	return fmt.Sprintf("%s No entry with a reason that starts with '%s' was found", n.link(id), reason)
}

// Done is a successful drop of the token
func (n responses) Done(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n%s",
		n.finishedWithToken(i), n.nowHasToken(a))
}

// DoneNoOthers is a successful drop of the token when nobody can pick it up
func (n responses) DoneNoOthers(i queue.Item) string {
	return fmt.Sprintf("%s\n%s",
		n.finishedWithToken(i), n.upForGrabs())
}

// DoneNotActive tells the user that they must have the token to use done
func (n responses) DoneNotActive(user string) string {
	return fmt.Sprintf("%s You cannot be done if you don't have the token", n.link(user))
}

// Yield is a successful passing of the token to next in line
func (n responses) Yield(i queue.Item, q queue.Queue) string {
	a := q.Active()
	return fmt.Sprintf("%s\n%s", n.yielded(i), n.nowHasToken(a))
}

// YieldNoOthers tells the user that they cannot yield if nobody is waiting
func (n responses) YieldNoOthers(i queue.Item) string {
	return fmt.Sprintf("%s You cannot yield if there is nobody waiting", n.link(i.ID))
}

// YieldNotActive tells the user they that must have the token to yield
func (n responses) YieldNotActive(i queue.Item) string {
	return fmt.Sprintf("%s You cannot yield if you do not have the token", n.link(i.ID))
}

// Barge is a successful barge to the front of the queue
func (n responses) Barge(i queue.Item, a queue.Item) string {
	return fmt.Sprintf("%s barged to the front\n%s still has the token", n.item(i), n.item(a))
}

// Boot is a successful force remove from the queue
func (n responses) Boot(booter string, i queue.Item) string {
	return fmt.Sprintf("%s booted %s from the list", n.link(booter), n.item(i))
}

// BootNoEntry tells the user that an entry with the requested user and reason does not exist
func (n responses) BootNoEntry(id, name, reason string) string {
	return fmt.Sprintf("%s No entry with for %s with a reason that starts with '%s' was found", n.link(id), name, reason)
}

// OustNotBoot tells the user that they can't boot the token holder
func (n responses) OustNotBoot(booter string) string {
	return fmt.Sprintf("%s You must oust the token holder", n.link(booter))
}

// Oust is a successful oust
func (n responses) Oust(ouster string, i queue.Item, a queue.Item) string {
	return fmt.Sprintf("%s\n%s", n.ousted(ouster, i), n.nowHasToken(a))
}

// OustNotActive tells the user they can only oust the token holder
func (n responses) OustNotActive(ouster string) string {
	return fmt.Sprintf("%s You can only oust the token holder", n.link(ouster))
}

// OustNoTarget tells the user they must specify a target to oust
func (n responses) OustNoTarget(ouster string) string {
	return fmt.Sprintf("%s You must specify who you want to oust", n.link(ouster))
}

// OustNoOthers is a successful oust when nobody can pick up the token
func (n responses) OustNoOthers(ouster string, i queue.Item) string {
	return fmt.Sprintf("%s\n%s", n.ousted(ouster, i), n.upForGrabs())
}

func (n responses) Delegate(i queue.Item, target string) string {
	return fmt.Sprintf("%s has delegated to %s", n.item(i), n.link(target))
}

func (n responses) DelegateActive(i queue.Item, ni queue.Item) string {
	return fmt.Sprintf("%s\n%s", n.Delegate(i, ni.ID), n.nowHasToken(ni))
}

func (n responses) DelegateNoSuchUser(delegator, target string) string {
	return fmt.Sprintf("%s You cannot delegate to %s because they don't exist", n.link(delegator), target)
}

func (n responses) DelegateNoEntry(delegator string) string {
	return fmt.Sprintf("%s You cannot delegate if you are not in the queue", n.link(delegator))
}

func (n responses) RefuseToken() string {
	return "What am I going to do with the token?"
}

func (n responses) RefuseTokenActive(i queue.Item, ni queue.Item) string {
	return fmt.Sprintf("%s\n:zap: :zap: AT LAST! ULTIMATE POWER! :zap: :zap:\n\nJust kidding ... I don't need the token, you can have it back\n%s", n.DelegateActive(i, ni), n.DelegateActive(ni, i))
}
