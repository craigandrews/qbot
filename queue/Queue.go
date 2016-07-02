package queue

type QueueError struct {
	msg string
}

func (e QueueError) Error() string {
	return e.msg
}

// Item represents a person with a job in the queue
type Item struct {
	Name   string
	Reason string
}

// Queue represents a list of waiting items
type Queue []Item

// Add appends an item to the queue unless it already exists
func (q Queue) Add(i Item) Queue {
	if q.Contains(i) {
		return q
	}
	q = append(q, i)
	return q
}

// Contains returns true if the item exists in the queue
func (q Queue) Contains(i Item) bool {
	for _, n := range q {
		if n == i {
			return true
		}
	}
	return false
}

// Active returns the first item in the queue or panics if the queue is empty
func (q Queue) Active() Item {
	if len(q) > 0 {
		return q[0]
	}
	panic(QueueError{"Queue is empty"})
}

// Waiting returns all items in the queue in order except the Active item
func (q Queue) Waiting() []Item {
	if len(q) > 0 {
		return q[1:]
	}
	return []Item(q)
}

// Remove removes an item from the queue
func (q Queue) Remove(i Item) Queue {
	for ix := range q {
		if q[ix] == i {
			if ix == 0 {
				return Queue(q[1:])
			} else if ix == len(q) - 1 {
				return Queue(q[:ix])
			}
			return Queue(append(q[:ix], q[ix+1:]...))
		}
	}
	return q
}

// Yield swaps the Active item with the first Waiting item
func (q Queue) Yield() Queue {
	if len(q) > 1 {
		q[0], q[1] = q[1], q[0]
	}
	return q
}

// Barge adds a new item to the second place in the queue, or moves an existing item to second place
func (q Queue) Barge(i Item) Queue {
	if len(q) > 1 && q.Active() != i {
		w := q.Remove(i).Waiting()
		q := Queue{q.Active(), i}
		return Queue(append(q, w...))
	}
	return q.Add(i)
}