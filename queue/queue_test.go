package queue

import (
	"testing"
	"reflect"
)

var john = Item{"john", "done some coding"}
var jimmy = Item{"jimmy", "fix some bugs"}
var mick = Item{"mick", "refactoring"}
var colin = Item{"colin", "adding bugs"}

func TestCreateQueue(t *testing.T) {
	q := Queue{}
	if len(q) != 0 {
		t.Errorf("Expected empty queue, but has length %d", len(q))
	}

	defer func() {
		if recover() == nil {
			t.Errorf("Expected no active user but one was found")
		}
	}()
	q.Active()
}

func TestCreateQueueWithEntries(t *testing.T) {
	q := Queue{mick, john}
	if len(q) != 2 {
		t.Errorf("Expected 2 entries in queue, but found %d", len(q))
	}
}

func TestAddImmutable(t *testing.T) {
	q := Queue{}
	q.Add(mick)
	if len(q) != 0 {
		t.Errorf("Expected empty queue, but has length %d", len(q))
	}
}

func TestAdd(t *testing.T) {
	q := Queue{}
	q = q.Add(mick)
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}

	active := q.Active()
	if active != mick {
		t.Errorf("Expected active user to be 'mick' but was '%s'", active)
	}
}

func TestAddDuplicate(t *testing.T) {
	q := Queue{mick}
	q = q.Add(mick)
	if len(q) != 1 {
		t.Errorf("Expecte 1 entry in queue, but found %d", len(q))
	}
}

func TestAddWithDifferentReason(t *testing.T) {
	q := Queue{mick}
	q = q.Add(Item{"mick", "wrote tests"})
	if len(q) != 2 {
		t.Errorf("Expecte 2 entry in queue, but found %d", len(q))
	}
}

func TestWaitingWhenEmpty(t *testing.T) {
	q := Queue{}
	w := q.Waiting()
	if len(w) > 0 {
		t.Errorf("Expected 0 waiting, but found %d", len(w))
	}
}

func TestWaitingWhenOnlyOne(t *testing.T) {
	q := Queue{mick}
	w := q.Waiting()
	if len(w) > 0 {
		t.Errorf("Expected 0 waiting, but found %d", len(w))
	}
}

func TestWaitingWhenMoreThanOne(t *testing.T) {
	q := Queue{mick, john, jimmy}
	expected := []Item{john, jimmy}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected %v waiting, but found %v", expected, q.Waiting())
	}
}

func TestRemoveImmutable(t *testing.T) {
	q := Queue{mick}
	q.Remove(mick)
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}
}

func TestRemoveActive(t *testing.T) {
	q := Queue{mick, john}
	q = q.Remove(mick)
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}

	if q.Active() != john {
		t.Errorf("Expected 'john' to be active but was '%s'", q.Active())
	}
}

func TestRemoveMiddle(t *testing.T) {
	q := Queue{mick, jimmy, john}
	q = q.Remove(jimmy)
	if len(q) != 2 {
		t.Errorf("Expected 2 entry in queue, but found %d", len(q))
	}
	if q.Contains(jimmy) {
		t.Errorf("Expected 'jimmy' to be removed but was found")
	}
}

func TestRemoveLast(t *testing.T) {
	q := Queue{mick, john}
	q = q.Remove(john)
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}
	if q.Contains(john) {
		t.Errorf("Expected 'john' to be removed but was found")
	}
}

func TestRemoveNotPresent(t *testing.T) {
	q := Queue{mick, john}
	q = q.Remove(jimmy)
	if len(q) != 2 {
		t.Errorf("Expected 2 entry in queue, but found %d", len(q))
	}
}

func TestYield(t *testing.T) {
	q := Queue{mick, john}
	q = q.Yield()
	if q.Active() != john {
		t.Errorf("Expected 'john' to be active but was '%s'", q.Active())
	}
}

func TestYieldAlone(t *testing.T) {
	q := Queue{mick}
	q = q.Yield()
	if q.Active() != mick {
		t.Errorf("Expected 'mick' to be active but was '%s'", q.Active())
	}
}

func TestBargeWhenEmpty(t *testing.T) {
	q := Queue{}
	q = q.Barge(mick)
	if q.Active() != mick {
		t.Errorf("Expected 'mick' to be active but was '%s'", q.Active())
	}
}

func TestBargeWhenActive(t *testing.T) {
	q := Queue{mick, john}
	q = q.Barge(mick)
	if q.Active() != mick {
		t.Errorf("Expected 'mick' to be active but was '%s'", q.Active())
	}
}

func TestBargeWhenOnlyOne(t *testing.T) {
	q := Queue{john}
	q = q.Barge(mick)
	expected := []Item{mick}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}

func TestBargeWhenOnlyTwoAndAlreadySecond(t *testing.T) {
	q := Queue{john, mick}
	q = q.Barge(mick)
	expected := []Item{mick}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}

func TestBargeWhenOnlyTwo(t *testing.T) {
	q := Queue{john, jimmy}
	q = q.Barge(mick)
	expected := []Item{mick, jimmy}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}

func TestBargeWhenAlreadySecond(t *testing.T) {
	q := Queue{john, mick}
	q = q.Barge(mick)
	expected := []Item{mick}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}

func TestBargeWhenNotPresent(t *testing.T) {
	q := Queue{john, colin, jimmy}
	q = q.Barge(mick)
	expected := []Item{mick, colin, jimmy}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}

func TestBargeWhenPresent(t *testing.T) {
	q := Queue{john, colin, mick, jimmy}
	q = q.Barge(mick)
	expected := []Item{mick, colin, jimmy}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}