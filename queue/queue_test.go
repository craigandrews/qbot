package queue

import "testing"

func TestCreateQueue(t *testing.T) {
	q := Queue{}
	if len(q) != 0 {
		t.Errorf("Expected empty queue, but has length %d", len(q))
	}

	active := q.Active()
	if active != "" {
		t.Errorf("Expected no active user but was '%s'", active)
	}
}

func TestCreateQueueWithEntries(t *testing.T) {
	q := Queue{"mick", "john"}
	if len(q) != 2 {
		t.Errorf("Expected 2 entries in queue, but found %d", len(q))
	}
}

func TestAddImmutable(t *testing.T) {
	q := Queue{}
	q.Add("mick")
	if len(q) != 0 {
		t.Errorf("Expected empty queue, but has length %d", len(q))
	}
}

func TestAdd(t *testing.T) {
	q := Queue{}
	q = q.Add("mick")
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}

	active := q.Active()
	if active != "mick" {
		t.Errorf("Expected active user to be 'mick' but was '%s'", active)
	}
}

func TestAddDuplicate(t *testing.T) {
	q := Queue{"mick"}
	q = q.Add("mick")
	if len(q) != 1 {
		t.Errorf("Expecte 1 entry in queue, but found %d", len(q))
	}
}

func TestRemoveImmutable(t *testing.T) {
	q := Queue{"mick"}
	q.Remove("mick")
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}
}

func TestRemoveActive(t *testing.T) {
	q := Queue{"mick", "john"}
	q = q.Remove("mick")
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}
	if q.Active() != "john" {
		t.Errorf("Expected 'john' to be active but was '%s'", q.Active())
	}
}

func TestRemoveMiddle(t *testing.T) {
	q := Queue{"mick", "jimmy", "john"}
	q = q.Remove("jimmy")
	if len(q) != 2 {
		t.Errorf("Expected 2 entry in queue, but found %d", len(q))
	}
	if q.Contains("jimmy") {
		t.Errorf("Expected 'jimmy' to be removed but was found")
	}
}

func TestRemoveLast(t *testing.T) {
	q := Queue{"mick", "john"}
	q = q.Remove("john")
	if len(q) != 1 {
		t.Errorf("Expected 1 entry in queue, but found %d", len(q))
	}
	if q.Contains("john") {
		t.Errorf("Expected 'john' to be removed but was found")
	}
}

func TestRemoveNotPresent(t *testing.T) {
	q := Queue{"mick", "john"}
	q = q.Remove("jimmy")
	if len(q) != 2 {
		t.Errorf("Expected 2 entry in queue, but found %d", len(q))
	}
}

func TestYield(t *testing.T) {
	q := Queue{"mick", "john"}
	q = q.Yield()
	if q.Active() != "john" {
		t.Errorf("Expected 'john' to be active but was '%s'", q.Active())
	}
}

func TestYieldAlone(t *testing.T) {
	q := Queue{"mick"}
	q = q.Yield()
	if q.Active() != "mick" {
		t.Errorf("Expected 'john' to be active but was '%s'", q.Active())
	}
}