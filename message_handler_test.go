package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/doozr/guac"
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

	expectedQueue := queue.Queue{
		queue.Item{ID: "U1234", Reason: "Some reason"},
	}
	expectedNotification := command.Notification{Channel: "C1234", Message: "This is a message"}
	commands := map[string]command.CmdFn{
		"test": func(q queue.Queue, channel string, user string, args string) (queue.Queue, command.Notification) {
			if !q.Equal(initialQueue) {
				t.Fatal("Incorrect queue passed to command ", initialQueue, q)
			}
			if channel != event.Channel {
				t.Fatal("Incorrect channel passed to command ", event.Channel, channel)
			}
			if user != event.User {
				t.Fatal("Incorrect user passed to command ", event.User, user)
			}
			if args != "the args" {
				t.Fatal("Incorrect args passed to command ", "the args", args)
			}
			return expectedQueue, expectedNotification
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

	handler := createMessageHandler(initialQueue, commands, notify, persist)
	err := handler(event)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

	if !reflect.DeepEqual(expectedNotification, receivedNotification) {
		t.Fatal("Received unexpected notification ", expectedNotification, receivedNotification)
	}

	if !receivedQueue.Equal(expectedQueue) {
		t.Fatal("Received unexpected qeueue", expectedQueue, receivedQueue)
	}
}

func TestDispatchCaseInsensitive(t *testing.T) {
	initialQueue := queue.Queue{}
	event := makeTestEvent("TEST UPPER CASE")

	calls := 0
	commands := map[string]command.CmdFn{
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

	handler := createMessageHandler(initialQueue, commands, notify, persist)
	err := handler(event)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

	if calls != 1 {
		t.Fatalf("Expected command to be called exactly once, was called %d times", calls)
	}
}

func TestDoesNothingIfNoMatchingCommand(t *testing.T) {
	initialQueue := queue.Queue{}
	event := makeTestEvent("NOT FOUND")

	commands := map[string]command.CmdFn{
		"test": func(q queue.Queue, channel string, user string, args string) (queue.Queue, command.Notification) {
			t.Fatal("Unexpected call to command")
			return q, command.Notification{Channel: channel, Message: "response"}
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

	handler := createMessageHandler(initialQueue, commands, notify, persist)
	err := handler(event)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}
}

func TestDoesNotPersistIfNotifyFails(t *testing.T) {
	initialQueue := queue.Queue{}
	event := makeTestEvent("test with errors")

	commands := map[string]command.CmdFn{
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

	handler := createMessageHandler(initialQueue, commands, notify, persist)
	err := handler(event)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestReturnsErrorIfPersistFails(t *testing.T) {
	initialQueue := queue.Queue{}
	event := makeTestEvent("test with errors")

	commands := map[string]command.CmdFn{
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

	handler := createMessageHandler(initialQueue, commands, notify, persist)
	err := handler(event)
	if err == nil {
		t.Fatal("Expected error")
	}
}
