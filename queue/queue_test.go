package queue

import (
	"reflect"
	"testing"
)

var John = Item{"john", "done some coding"}
var Jimmy = Item{"jimmy", "fix some bugs"}
var Mick = Item{"mick", "refactoring"}
var Colin = Item{"colin", "adding bugs"}

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
	q := Queue{Mick, John}
	if len(q) != 2 {
		t.Errorf("Expected 2 entries in queue, but found %d", len(q))
	}
}

func TestAddImmutable(t *testing.T) {
	q := Queue{}
	q.Add(Mick)
	if len(q) != 0 {
		t.Errorf("Expected empty queue, but has length %d", len(q))
	}
}

func TestAdd(t *testing.T) {
	q := Queue{}
	q = q.Add(Mick)
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}

	active := q.Active()
	if active != Mick {
		t.Errorf("Expected active user to be 'mick' but was '%s'", active)
	}
}

func TestAddDuplicate(t *testing.T) {
	q := Queue{Mick}
	q = q.Add(Mick)
	if len(q) != 1 {
		t.Errorf("Expecte 1 entry in queue, but found %d", len(q))
	}
}

func TestAddWithDifferentReason(t *testing.T) {
	q := Queue{Mick}
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
	q := Queue{Mick}
	w := q.Waiting()
	if len(w) > 0 {
		t.Errorf("Expected 0 waiting, but found %d", len(w))
	}
}

func TestWaitingWhenMoreThanOne(t *testing.T) {
	q := Queue{Mick, John, Jimmy}
	expected := []Item{John, Jimmy}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected %v waiting, but found %v", expected, q.Waiting())
	}
}

func TestRemoveImmutable(t *testing.T) {
	q := Queue{Mick}
	q.Remove(Mick)
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}
}

func TestRemoveActive(t *testing.T) {
	q := Queue{Mick, John}
	q = q.Remove(Mick)
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}

	if q.Active() != John {
		t.Errorf("Expected 'john' to be active but was '%s'", q.Active())
	}
}

func TestRemoveMiddle(t *testing.T) {
	q := Queue{Mick, Jimmy, John}
	q = q.Remove(Jimmy)
	if len(q) != 2 {
		t.Errorf("Expected 2 entry in queue, but found %d", len(q))
	}
	if q.Contains(Jimmy) {
		t.Errorf("Expected 'jimmy' to be removed but was found")
	}
}

func TestRemoveLast(t *testing.T) {
	q := Queue{Mick, John}
	q = q.Remove(John)
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}
	if q.Contains(John) {
		t.Errorf("Expected 'john' to be removed but was found")
	}
}

func TestRemoveNotPresent(t *testing.T) {
	q := Queue{Mick, John}
	q = q.Remove(Jimmy)
	if len(q) != 2 {
		t.Errorf("Expected 2 entry in queue, but found %d", len(q))
	}
}

func TestYield(t *testing.T) {
	q := Queue{Mick, John}
	q = q.Yield()
	if q.Active() != John {
		t.Errorf("Expected 'john' to be active but was '%s'", q.Active())
	}
}

func TestYieldAlone(t *testing.T) {
	q := Queue{Mick}
	q = q.Yield()
	if q.Active() != Mick {
		t.Errorf("Expected 'mick' to be active but was '%s'", q.Active())
	}
}

func TestBargeWhenEmpty(t *testing.T) {
	q := Queue{}
	q = q.Barge(Mick)
	if q.Active() != Mick {
		t.Errorf("Expected 'mick' to be active but was '%s'", q.Active())
	}
}

func TestBargeWhenActive(t *testing.T) {
	q := Queue{Mick, John}
	q = q.Barge(Mick)
	if q.Active() != Mick {
		t.Errorf("Expected 'mick' to be active but was '%s'", q.Active())
	}
}

func TestBargeWhenOnlyOne(t *testing.T) {
	q := Queue{John}
	q = q.Barge(Mick)
	expected := []Item{Mick}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}

func TestBargeWhenOnlyTwoAndAlreadySecond(t *testing.T) {
	q := Queue{John, Mick}
	q = q.Barge(Mick)
	expected := []Item{Mick}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}

func TestBargeWhenOnlyTwo(t *testing.T) {
	q := Queue{John, Jimmy}
	q = q.Barge(Mick)
	expected := []Item{Mick, Jimmy}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}

func TestBargeWhenAlreadySecond(t *testing.T) {
	q := Queue{John, Mick}
	q = q.Barge(Mick)
	expected := []Item{Mick}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}

func TestBargeWhenNotPresent(t *testing.T) {
	q := Queue{John, Colin, Jimmy}
	q = q.Barge(Mick)
	expected := []Item{Mick, Colin, Jimmy}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}

func TestBargeWhenPresent(t *testing.T) {
	q := Queue{John, Colin, Mick, Jimmy}
	q = q.Barge(Mick)
	expected := []Item{Mick, Colin, Jimmy}
	if !reflect.DeepEqual(q.Waiting(), expected) {
		t.Errorf("Expected waiting to be %v but was %v", expected, q.Waiting())
	}
}
