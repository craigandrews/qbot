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
	q, m := Join(q, Mick)
	if len(q) < 1 || q.Active() != Mick {
		t.Errorf("Expected %v to be active but queue was empty", Mick)
	}
	if m != notification.Active(Mick) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.Active(Mick), m)
	}
}

func TestJoin(t *testing.T) {
	q := queue.Queue{John, Jimmy}
	q, m := Join(q, Mick)
	if len(q) != 3 || !q.Contains(Mick) {
		t.Errorf("Expected %v to be in the queue but was not found", Mick)
	}
	if m != notification.Join(Mick) {
		t.Errorf("Expected notification of '%s' but got '%s'", notification.Join(Mick), m)
	}
}

func TestJoinAlreadyExists(t *testing.T) {
	q := queue.Queue{John, Jimmy, Mick}
	q, m := Join(q, Mick)
	if len(q) != 3 || !q.Contains(Mick) {
		t.Errorf("Expected %v to be in the queue but was not found", Mick)
	}
	if m != "" {
		t.Errorf("Expected notification of '%s' but got '%s'", "", m)
	}
}