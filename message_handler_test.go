package qbot_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/doozr/guac"
	. "github.com/doozr/qbot"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func makeTestEvent(text string) guac.MessageEvent {
	return guac.MessageEvent{
		Type:    "message",
		ID:      123,
		Channel: "C1234",
		User:    "U1234",
		Text:    text,
	}
}

func TestDispatchesMessage(t *testing.T) {
	initialQueue := queue.Queue{}
	event := makeTestEvent("test the args")

	commands := map[string]command.Command{
		"test": func(q queue.Queue, channel string, user string, args string) (queue.Queue, command.Notification) {
			q = q.Add(queue.Item{ID: user, Reason: args})
			n := command.Notification{Channel: channel, Message: "This is a message"}
			return q, n
		},
	}

	var receivedNotification command.Notification
	notify := func(n command.Notification) error {
		receivedNotification = n
		return nil
	}

	var receivedQueue queue.Queue
	persist := func(q queue.Queue) error {
		receivedQueue = q
		return nil
	}

	handler := CreateMessageHandler(commands, notify, persist)
	_, err := handler(initialQueue, event)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

	expectedNotification := command.Notification{
		Channel: "C1234",
		Message: "This is a message",
	}
	if !reflect.DeepEqual(expectedNotification, receivedNotification) {
		t.Fatal("Received unexpected notification ", expectedNotification, receivedNotification)
	}

	expectedQueue := queue.Queue([]queue.Item{
		queue.Item{ID: "U1234", Reason: "the args"},
	})
	if !receivedQueue.Equal(expectedQueue) {
		t.Fatal("Received unexpected qeueue", expectedQueue, receivedQueue)
	}
}

func TestDispatchCaseInsensitive(t *testing.T) {
	initialQueue := queue.Queue{}
	event := makeTestEvent("TEST UPPER CASE")

	calls := 0
	commands := map[string]command.Command{
		"test": func(q queue.Queue, channel string, user string, args string) (queue.Queue, command.Notification) {
			calls++
			return q, command.Notification{Channel: channel, Message: "response"}
		},
	}

	notify := func(n command.Notification) error {
		return nil
	}

	persist := func(q queue.Queue) error {
		return nil
	}

	handler := CreateMessageHandler(commands, notify, persist)
	handler(initialQueue, event)

	if calls != 1 {
		t.Fatalf("Expected command to be called exactly once, was called %d times", calls)
	}
}

func TestDoesNothingIfNoMatchingCommand(t *testing.T) {
	initialQueue := queue.Queue{}
	event := makeTestEvent("NOT FOUND")

	commands := map[string]command.Command{
		"test": func(q queue.Queue, channel string, user string, args string) (queue.Queue, command.Notification) {
			t.Fatal("Unexpected call to command")
			return q, command.Notification{}
		},
	}

	notify := func(n command.Notification) error {
		t.Fatal("Unexpected call to notify")
		return nil
	}

	persist := func(q queue.Queue) error {
		t.Fatal("Unexpected called to persist")
		return nil
	}

	handler := CreateMessageHandler(commands, notify, persist)
	handler(initialQueue, event)
}

func TestDoesNotPersistIfNotifyFails(t *testing.T) {
	initialQueue := queue.Queue{}
	event := makeTestEvent("test with errors")

	commands := map[string]command.Command{
		"test": func(q queue.Queue, channel string, user string, args string) (queue.Queue, command.Notification) {
			return q, command.Notification{Channel: channel, Message: "response"}
		},
	}

	notify := func(n command.Notification) error {
		return fmt.Errorf("Error!")
	}

	persist := func(q queue.Queue) error {
		t.Fatal("Unexpected called to persist")
		return nil
	}

	handler := CreateMessageHandler(commands, notify, persist)
	_, err := handler(initialQueue, event)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestReturnsErrorIfPersistFails(t *testing.T) {
	initialQueue := queue.Queue{}
	event := makeTestEvent("test with errors")

	commands := map[string]command.Command{
		"test": func(q queue.Queue, channel string, user string, args string) (queue.Queue, command.Notification) {
			return q, command.Notification{Channel: channel, Message: "response"}
		},
	}

	notify := func(n command.Notification) error {
		return nil
	}

	persist := func(q queue.Queue) error {
		return fmt.Errorf("Error!")
	}

	handler := CreateMessageHandler(commands, notify, persist)
	_, err := handler(initialQueue, event)
	if err == nil {
		t.Fatal("Expected error")
	}
}

// returns new queue from command
// returns same queue if invalid command
