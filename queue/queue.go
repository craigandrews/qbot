package queue

// Item represents a person with a job in the queue
type Item struct {
	ID     string
	Reason string
}

// Queue represents a list of waiting items
type Queue []Item

// Equal checks if the queue is the same as another queue
func (q Queue) Equal(other Queue) bool {
	if len(q) != len(other) {
		return false
	}

	for ix := range q {
		if q[ix] != other[ix] {
			return false
		}
	}

	return true
}

// Add appends an item to the queue unless it already exists
func (q Queue) Add(i Item) Queue {
	if q.Contains(i) {
		return q
	}
	return append(q, i)
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

// Active returns the first item in the queue or empty item if the queue is empty
func (q Queue) Active() Item {
	if len(q) > 0 {
		return q[0]
	}
	return Item{}
}

// Waiting returns all items in the queue in order except the Active item
func (q Queue) Waiting() []Item {
	if len(q) > 1 {
		return q[1:]
	}
	return []Item{}
}

// Remove removes an item from the queue
func (q Queue) Remove(i Item) Queue {
	for ix := range q {
		if q[ix] == i {
			if ix == 0 {
				return q[1:]
			} else if ix == len(q)-1 {
				return q[:ix]
			}
			return append(q[:ix], q[ix+1:]...)
		}
	}
	return q
}

// Yield swaps the Active item with the first Waiting item
func (q Queue) Yield() Queue {
	if len(q) < 2 {
		return q
	}
	q[0], q[1] = q[1], q[0]
	return q
}

// Barge adds a new item to the second place in the queue, or moves an existing item to second place
func (q Queue) Barge(i Item) Queue {
	if q.Active() == i {
		return q
	}

	if len(q) < 2 {
		return q.Add(i)
	}

	w := q.Remove(i).Waiting()
	return append(Queue{q.Active(), i}, w...)
}

// Delegate swaps an item for another in the same position
func (q Queue) Delegate(i Item, n Item) Queue {
	for ix := range q {
		if q[ix] == i {
			q[ix] = n
		}
	}
	return q
}
