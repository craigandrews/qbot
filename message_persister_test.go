package qbot_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/doozr/guac"
	. "github.com/doozr/qbot"
	"github.com/doozr/qbot/queue"
)

func TestPassesOnParameters(t *testing.T) {
	var receivedQueue queue.Queue
	var receivedEvent guac.MessageEvent
	fn := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		receivedQueue = q
		receivedEvent = m
		return q, nil
	}

	persist := func(q queue.Queue) error {
		return nil
	}

	event := makeTestEvent("text")
	expectedQueue := queue.Queue([]queue.Item{{"U123", "Tomato"}})

	handler := CreateMessagePersister(persist, fn)
	handler(expectedQueue, event)

	if !expectedQueue.Equal(receivedQueue) {
		t.Fatal("Unexpected qeueu", expectedQueue, receivedQueue)
	}

	if !reflect.DeepEqual(event, receivedEvent) {
		t.Fatal("Unexpected event", event, receivedEvent)
	}
}

func TestPersistsReturnedQueue(t *testing.T) {
	expectedQueue := queue.Queue([]queue.Item{{"U123", "Tomato"}})

	fn := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		return expectedQueue, nil
	}

	var receivedQueue queue.Queue
	persist := func(q queue.Queue) error {
		receivedQueue = q
		return nil
	}

	event := makeTestEvent("text")

	handler := CreateMessagePersister(persist, fn)
	handler(queue.Queue{}, event)

	if !expectedQueue.Equal(receivedQueue) {
		t.Fatal("Unexpected qeueu", expectedQueue, receivedQueue)
	}
}

func TestReturnsReturnedQueue(t *testing.T) {
	expectedQueue := queue.Queue([]queue.Item{{"U123", "Tomato"}})

	fn := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		return expectedQueue, nil
	}

	persist := func(q queue.Queue) error {
		return nil
	}

	event := makeTestEvent("text")

	handler := CreateMessagePersister(persist, fn)
	receivedQueue, _ := handler(queue.Queue{}, event)

	if !expectedQueue.Equal(receivedQueue) {
		t.Fatal("Unexpected qeueu", expectedQueue, receivedQueue)
	}
}

func TestDoesNotPersistOnError(t *testing.T) {
	fn := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		return nil, fmt.Errorf("Error!")
	}

	calls := 0
	persist := func(q queue.Queue) error {
		calls++
		return nil
	}

	event := makeTestEvent("text")

	handler := CreateMessagePersister(persist, fn)
	_, err := handler(queue.Queue{}, event)

	if calls != 0 {
		t.Fatal("Expected 0 calls, recieved ", calls)
	}

	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestReturnsPersistError(t *testing.T) {
	fn := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		return q, nil
	}

	persist := func(q queue.Queue) error {
		return fmt.Errorf("Error!")
	}

	event := makeTestEvent("text")

	handler := CreateMessagePersister(persist, fn)
	_, err := handler(queue.Queue{}, event)

	if err == nil {
		t.Fatal("Expected error")
	}
}
