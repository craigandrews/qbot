package queue

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
)

// Error is an error thrown when manipulating the queue
type Error struct {
	msg string
}

// Error returns the error message
func (e Error) Error() string {
	return e.msg
}

// Item represents a person with a job in the queue
type Item struct {
	ID     string
	Reason string
}

// Queue represents a list of waiting items
type Queue []Item

// Load populates a new queue from a JSON file
func Load(filename string) (Queue, error) {
	q := Queue{}
	if _, err := os.Stat(filename); err == nil {
		dat, err := ioutil.ReadFile(filename)
		if err != nil {
			return q, err
		}
		json.Unmarshal(dat, &q)
	}
	return q, nil
}

// Equal checks if the queue is the same as another queue
func (q Queue) Equal(other Queue) bool {
	return reflect.DeepEqual(q, other)
}

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
	panic(Error{"Queue is empty"})
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
			} else if ix == len(q)-1 {
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
		oq := q
		q = Queue{q[1], q[0]}
		if len(oq) > 2 {
			q = append(q, oq[2:]...)
		}
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

// Save serialises the queue to disk
func (q Queue) Save(filename string) (err error) {
	j, err := json.Marshal(q)
	err = ioutil.WriteFile(filename, j, 0644)
	return
}
