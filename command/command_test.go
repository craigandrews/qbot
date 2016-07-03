package command

import (
	"testing"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/notification"
)

var John = queue.Item{"john", "done some coding"}
var Jimmy = queue.Item{"jimmy", "fix some bugs"}
var Mick = queue.Item{"mick", "refactoring"}

func TestJoinEmptyQueue(t *testing.T) {
	q := queue.Queue{}
	q, m := Join(q, Mick.Name, Mick.Reason)
	if len(q) < 1 || q.Active() != Mick {
		t.Errorf("Expected %v to be active but queue was empty", Mick)
	}
	if m != notification.Active(Mick) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.Active(Mick), m)
	}
}

func TestJoin(t *testing.T) {
	q := queue.Queue{John, Jimmy}
	q, m := Join(q, Mick.Name, Mick.Reason)
	if len(q) != 3 || !q.Contains(Mick) {
		t.Errorf("Expected %v to be in the queue but was not found", Mick)
	}
	if m != notification.Join(Mick) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.Join(Mick), m)
	}
}

func TestJoinAlreadyExists(t *testing.T) {
	q := queue.Queue{John, Jimmy, Mick}
	q, m := Join(q, Mick.Name, Mick.Reason)
	if len(q) != 3 || !q.Contains(Mick) {
		t.Errorf("Expected %v to be in the queue but was not found", Mick)
	}
	if m != "" {
		t.Errorf("Expected notification of '%s' but got '%s'", "", m)
	}
}

func TestLeaveWhenNotPresent(t *testing.T) {
	q := queue.Queue{John, Jimmy}
	q, m := Leave(q, Mick.Name, "")
	if len(q) != 2 {
		t.Errorf("Expected 2 items but got %d", len(q))
	}
	if m != "" {
		t.Errorf("Expected empty message but got '%s'", m)
	}
}

func TestLeave(t *testing.T) {
	q := queue.Queue{John, Mick, Jimmy}
	q, m := Leave(q, Mick.Name, "")
	if q.Contains(Mick) {
		t.Errorf("Expected %v to be missing but was present", Mick)
	}
	if m != notification.Leave(Mick) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.Leave(Mick), m)
	}
}

func TestLeaveWithMulti(t *testing.T) {
	i := queue.Item{"mick", "potato"}
	q := queue.Queue{John, Mick, Jimmy, i}
	q, m := Leave(q, Mick.Name, "")
	if q.Contains(i) {
		t.Errorf("Expected %v to be missing but was present", i)
	}
	if !q.Contains(Mick) {
		t.Errorf("Expected %v to be present but was missing", Mick)
	}
	if m != notification.Leave(i) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.Leave(i), m)
	}
}

func TestLeaveWithPrefix(t *testing.T) {
	i := queue.Item{"mick", "potato"}
	q := queue.Queue{John, Mick, Jimmy, i}
	q, m := Leave(q, Mick.Name, "refac")
	if q.Contains(Mick) {
		t.Errorf("Expected %v to be missing but was present", Mick)
	}
	if !q.Contains(i) {
		t.Errorf("Expected %v to be present but was missing", i)
	}
	if m != notification.Leave(Mick) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.Leave(i), m)
	}
}

func TestLeaveWhenActive(t *testing.T) {
	q := queue.Queue{Mick, John, Jimmy}
	q, m := Leave(q, Mick.Name, "")
	if len(q) != 2 {
		t.Errorf("Expected 2 items but got %d", len(q))
	}
	if q.Active() != John {
		t.Errorf("Expected %v to be active but was %v", John, q.Active())
	}
	if m != notification.LeaveActive(Mick, q) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.LeaveActive(Mick, q), m)
	}
}

func TestLeaveWhenActiveAndAlone(t *testing.T) {
	q := queue.Queue{Mick}
	q, m := Leave(q, Mick.Name, "")
	if len(q) != 0 {
		t.Errorf("Expected 0 items but got %d", len(q))
	}
	if m != notification.LeaveNoActive(Mick) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.LeaveNoActive(Mick), m)
	}
}

func TestDone(t *testing.T) {
	q := queue.Queue{Mick, John}
	q, m := Done(q, Mick.Name)
	if len(q) != 1 {
		t.Errorf("Expected 1 items but got %d", len(q))
	}
	if q.Active() != John {
		t.Errorf("Expected %v to be active but was %v", John, q.Active())
	}
	if m != notification.Done(Mick, q) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.Done(Mick, q), m)
	}
}

func TestDoneNoOthers(t *testing.T) {
	q := queue.Queue{Mick}
	q, m := Done(q, Mick.Name)
	if len(q) != 0 {
		t.Errorf("Expected 0 items but got %d", len(q))
	}
	if m != notification.DoneNoOthers(Mick) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.DoneNoOthers(Mick), m)
	}
}

func TestDoneNotActive(t *testing.T) {
	q := queue.Queue{John, Mick}
	q, m := Done(q, Mick.Name)
	if len(q) != 2 {
		t.Errorf("Expected 2 items but got %d", len(q))
	}
	if q.Active() != John {
		t.Errorf("Expected %v to be active but was %v", John, q.Active())
	}
	if m != notification.DoneNotActive(Mick) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.DoneNotActive(Mick), m)
	}
}
