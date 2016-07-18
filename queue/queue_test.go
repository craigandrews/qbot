package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var John = Item{"john", "done some coding"}
var Jimmy = Item{"jimmy", "fix some bugs"}
var Mick = Item{"mick", "refactoring"}
var Colin = Item{"colin", "adding bugs"}

func TestCreateQueue(t *testing.T) {
	q := Queue{}
	assert.Equal(t, 0, len(q))

	defer func() {
		assert.NotNil(t, recover(), "Expected no active user but found one")
	}()
	q.Active()
}

func TestCreateQueueWithEntries(t *testing.T) {
	q := Queue{Mick, John}
	assert.Equal(t, 2, len(q))
}

func TestAddImmutable(t *testing.T) {
	q := Queue{}
	q.Add(Mick)
	assert.Equal(t, 0, len(q))
}

func TestAdd(t *testing.T) {
	q := Queue{}
	q = q.Add(Mick)
	assert.Equal(t, 1, len(q))

	active := q.Active()
	assert.Equal(t, Mick, active)
}

func TestAddDuplicate(t *testing.T) {
	q := Queue{Mick}
	q = q.Add(Mick)
	assert.Equal(t, 1, len(q))
}

func TestAddWithDifferentReason(t *testing.T) {
	q := Queue{Mick}
	q = q.Add(Item{"mick", "wrote tests"})
	assert.Equal(t, 2, len(q))
}

func TestWaitingWhenEmpty(t *testing.T) {
	q := Queue{}
	w := q.Waiting()
	assert.Equal(t, 0, len(w))
}

func TestWaitingWhenOnlyOne(t *testing.T) {
	q := Queue{Mick}
	w := q.Waiting()
	assert.Equal(t, 0, len(w))
}

func TestWaitingWhenMoreThanOne(t *testing.T) {
	q := Queue{Mick, John, Jimmy}
	expected := []Item{John, Jimmy}
	assert.Equal(t, expected, q.Waiting())
}

func TestRemoveImmutable(t *testing.T) {
	q := Queue{Mick}
	q.Remove(Mick)
	assert.Equal(t, 1, len(q))
}

func TestRemoveActive(t *testing.T) {
	q := Queue{Mick, John}
	q = q.Remove(Mick)
	assert.Equal(t, 1, len(q))
	assert.Equal(t, John, q.Active())
}

func TestRemoveMiddle(t *testing.T) {
	q := Queue{Mick, Jimmy, John}
	q = q.Remove(Jimmy)
	assert.Equal(t, 2, len(q))
	assert.NotContains(t, q, Jimmy)
}

func TestRemoveLast(t *testing.T) {
	q := Queue{Mick, John}
	q = q.Remove(John)
	assert.Equal(t, 1, len(q))
	assert.NotContains(t, q, John)
}

func TestRemoveNotPresent(t *testing.T) {
	q := Queue{Mick, John}
	q = q.Remove(Jimmy)
	assert.Equal(t, 2, len(q))
}

func TestYield(t *testing.T) {
	q := Queue{Mick, John}
	q = q.Yield()
	assert.Equal(t, John, q.Active())
}

func TestYieldAlone(t *testing.T) {
	q := Queue{Mick}
	q = q.Yield()
	assert.Equal(t, Mick, q.Active())
}

func TestBargeWhenEmpty(t *testing.T) {
	q := Queue{}
	q = q.Barge(Mick)
	assert.Equal(t, Mick, q.Active())
}

func TestBargeWhenActive(t *testing.T) {
	q := Queue{Mick, John}
	q = q.Barge(Mick)
	assert.Equal(t, Mick, q.Active())
}

func TestBargeWhenOnlyOne(t *testing.T) {
	q := Queue{John}
	q = q.Barge(Mick)
	expected := []Item{Mick}
	assert.Equal(t, expected, q.Waiting())
}

func TestBargeWhenOnlyTwoAndAlreadySecond(t *testing.T) {
	q := Queue{John, Mick}
	q = q.Barge(Mick)
	expected := []Item{Mick}
	assert.Equal(t, expected, q.Waiting())
}

func TestBargeWhenOnlyTwo(t *testing.T) {
	q := Queue{John, Jimmy}
	q = q.Barge(Mick)
	expected := []Item{Mick, Jimmy}
	assert.Equal(t, expected, q.Waiting())
}

func TestBargeWhenAlreadySecond(t *testing.T) {
	q := Queue{John, Mick}
	q = q.Barge(Mick)
	expected := []Item{Mick}
	assert.Equal(t, expected, q.Waiting())
}

func TestBargeWhenNotPresent(t *testing.T) {
	q := Queue{John, Colin, Jimmy}
	q = q.Barge(Mick)
	expected := []Item{Mick, Colin, Jimmy}
	assert.Equal(t, expected, q.Waiting())
}

func TestBargeWhenPresent(t *testing.T) {
	q := Queue{John, Colin, Mick, Jimmy}
	q = q.Barge(Mick)
	expected := []Item{Mick, Colin, Jimmy}
	assert.Equal(t, expected, q.Waiting())
}
