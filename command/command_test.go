package command

import (
	"testing"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/notification"
	"fmt"
)

var John = queue.Item{"john", "done some coding"}
var Jimmy = queue.Item{"jimmy", "fix some bugs"}
var Mick = queue.Item{"mick", "refactoring"}

func assertQueueLength(q queue.Queue, n int) string {
	if len(q) != n {
		return fmt.Sprintf("Expected queue length to be %d but was %d", n, len(q))
	}
	return ""
}

func assertContains(q queue.Queue, i queue.Item) string {
	if !q.Contains(i) {
		return fmt.Sprintf("Expected %v to contain %v", q, i)
	}
	return ""
}

func assertNotContains(q queue.Queue, i queue.Item) string {
	if q.Contains(i) {
		return fmt.Sprintf("Expected %v not to contain %v", q, i)
	}
	return ""
}

func assertActive(q queue.Queue, i queue.Item) string {
	if q.Active() != i {
		return fmt.Sprintf("Expected %v to be active but was %v", i, q.Active())
	}
	return ""
}

func assertNotification(expected, actual string) string {
	if expected != actual {
		fmt.Sprintf("Expected notification of '%s' but got '%s'", expected, actual)
	}
	return ""
}

func TestJoinEmptyQueue(t *testing.T) {
	q := queue.Queue{}
	q, m := Join(q, Mick.Name, Mick.Reason)
	if e := assertQueueLength(q, 1); e != "" {
		t.Error(e)
	}
	if e := assertActive(q, Mick); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.JoinActive(Mick), m); e != "" {
		t.Error(e)
	}
}

func TestJoin(t *testing.T) {
	q := queue.Queue{John, Jimmy}
	q, m := Join(q, Mick.Name, Mick.Reason)
	if e := assertQueueLength(q, 3); e != "" {
		t.Error(e)
	}
	if e := assertContains(q, Mick); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.Join(Mick), m); e != "" {
		t.Error(e)
	}
}

func TestJoinAlreadyExists(t *testing.T) {
	q := queue.Queue{John, Jimmy, Mick}
	q, m := Join(q, Mick.Name, Mick.Reason)
	if e := assertQueueLength(q, 3); e != "" {
		t.Error(e)
	}
	if e := assertContains(q, Mick); e != "" {
		t.Error(e)
	}
	if e := assertNotification("", m); e != "" {
		t.Error(e)
	}
}

func TestLeaveWhenNotPresent(t *testing.T) {
	q := queue.Queue{John, Jimmy}
	q, m := Leave(q, Mick.Name, "")
	if e := assertQueueLength(q, 2); e != "" {
		t.Error(e)
	}
	if e := assertNotContains(q, Mick); e != "" {
		t.Error(e)
	}
	if e := assertNotification("", m); e != "" {
		t.Error(e)
	}
}

func TestLeave(t *testing.T) {
	q := queue.Queue{John, Mick, Jimmy}
	q, m := Leave(q, Mick.Name, "")
	if e := assertNotContains(q, Mick); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.Leave(Mick), m); e != "" {
		t.Error(e)
	}
}

func TestLeaveWithMulti(t *testing.T) {
	i := queue.Item{"mick", "potato"}
	q := queue.Queue{John, Mick, Jimmy, i}
	q, m := Leave(q, Mick.Name, "")
	if e := assertNotContains(q, i); e != "" {
		t.Error(e)
	}
	if e := assertContains(q, Mick); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.Leave(i), m); e != "" {
		t.Error(e)
	}
}

func TestLeaveWithPrefix(t *testing.T) {
	i := queue.Item{"mick", "potato"}
	q := queue.Queue{John, Mick, Jimmy, i}
	q, m := Leave(q, Mick.Name, "refac")
	if e := assertNotContains(q, Mick); e != "" {
		t.Error(e)
	}
	if e := assertContains(q, i); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.Leave(Mick), m); e != "" {
		t.Error(e)
	}
}

func TestLeaveWhenActive(t *testing.T) {
	q := queue.Queue{Mick, John, Jimmy}
	q, m := Leave(q, Mick.Name, "")
	if e := assertQueueLength( q, 2); e != "" {
		t.Error(e)
	}
	if e := assertActive(q, John); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.LeaveActive(Mick, q), m); e != "" {
		t.Error(e)
	}
}

func TestLeaveWhenActiveAndAlone(t *testing.T) {
	q := queue.Queue{Mick}
	q, m := Leave(q, Mick.Name, "")
	if e := assertQueueLength(q, 0); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.LeaveNoActive(Mick), m); e != "" {
		t.Error(e)
	}
}

func TestDone(t *testing.T) {
	q := queue.Queue{Mick, John}
	q, m := Done(q, Mick.Name)
	if e := assertQueueLength(q, 1); e != "" {
		t.Error(e)
	}
	if e := assertActive(q, John); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.Done(Mick, q), m); e != "" {
		t.Error(e)
	}
}

func TestDoneNoOthers(t *testing.T) {
	q := queue.Queue{Mick}
	q, m := Done(q, Mick.Name)
	if e := assertQueueLength(q, 0); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.DoneNoOthers(Mick), m); e != "" {
		t.Error(e)
	}
}

func TestDoneNotActive(t *testing.T) {
	q := queue.Queue{John, Mick}
	q, m := Done(q, Mick.Name)
	if e := assertQueueLength(q, 2); e != "" {
		t.Error(e)
	}
	if e := assertActive(q, John); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.DoneNotActive(Mick), m); e != "" {
		t.Error(e)
	}
}

func TestYield(t *testing.T) {
	q := queue.Queue{Mick, John}
	q, m := Yield(q, Mick.Name)
	if e := assertQueueLength(q, 2); e != "" {
		t.Error(e)
	}
	if e := assertActive(q, John); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.Yield(Mick, q), m); e != "" {
		t.Error(e)
	}
}

func TestYieldNoOthers(t *testing.T) {
	q := queue.Queue{Mick}
	q, m := Yield(q, Mick.Name)
	if e := assertQueueLength(q, 1); e != "" {
		t.Error(e)
	}
	if e := assertActive(q, Mick); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.YieldNoOthers(Mick), m); e != "" {
		t.Error(e)
	}
}

func TestYieldNotActive(t *testing.T) {
	q := queue.Queue{John, Mick}
	q, m := Yield(q, Mick.Name)
	if e := assertQueueLength(q, 2); e != "" {
		t.Error(e)
	}
	if e := assertActive(q, John); e != "" {
		t.Error(e)
	}
	if e := assertNotification(notification.YieldNotActive(Mick), m); e != "" {
		t.Error(e)
	}
}

func TestYieldEmpty(t *testing.T) {
	q := queue.Queue{}
	q, m := Yield(q, Mick.Name)
	if e := assertNotification(notification.YieldNotActive(Mick), m); e != "" {
		t.Error(e)
	}
}