package queue

type Queue []string
type QueueError struct {
	msg string
}

func (e QueueError) Error() string {
	return e.msg
}

func (q Queue) Add(u string) Queue {
	if q.Contains(u) {
		return q
	}
	q = append(q, u)
	return q
}

func (q Queue) Contains(u string) bool {
	for _, n := range q {
		if n == u {
			return true
		}
	}
	return false
}

func (q Queue) Active() string {
	if len(q) > 0 {
		return q[0]
	}
	return ""
}

func (q Queue) Remove(u string) Queue {
	for ix := range q {
		if q[ix] == u {
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

func (q Queue) Yield() Queue {
	if len(q) > 1 {
		q[0], q[1] = q[1], q[0]
	}
	return q
}