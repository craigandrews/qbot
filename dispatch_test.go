package qbot_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/doozr/guac"
	. "github.com/doozr/qbot"
	"github.com/doozr/qbot/queue"
)

func testDispatchCleanShutDown(t *testing.T, dispatcher Dispatcher) {
	done := make(DoneChan)
	events := make(guac.EventChan)
	waitGroup := sync.WaitGroup{}

	abort := Dispatch(dispatcher, events, done, &waitGroup)
	select {
	case <-abort:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected abort within 2 seconds")
	}

	waitGroup.Wait()
}

func TestDispatchRunsDispatcherInBackground(t *testing.T) {
	called := false
	dispatcher := func(events guac.EventChan, done DoneChan) error {
		called = true
		return nil
	}

	testDispatchCleanShutDown(t, dispatcher)
	if !called {
		t.Fatal("Dispatcher not called")
	}
}

func TestDispatchShutDownCleanlyWhenDispatcherReturnsError(t *testing.T) {
	dispatcher := func(events guac.EventChan, done DoneChan) error {
		return fmt.Errorf("Error!")
	}
	testDispatchCleanShutDown(t, dispatcher)
}

func testMessageDispatch(handleMessage MessageHandler, handleUserChange UserChangeHandler) error {
	done := make(DoneChan)
	events := make(guac.EventChan)

	go func() {
		events <- guac.MessageEvent{
			Text: "test event",
		}
		close(done)
	}()

	dispatcher := CreateDispatcher(queue.Queue{}, 1*time.Millisecond, handleMessage, handleUserChange)
	return dispatcher(events, done)
}

func TestDispatcherSendsMessagesToMessageHandler(t *testing.T) {
	var received *guac.MessageEvent
	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		received = &event
		return q, nil
	}
	handleUserChange := func(event guac.UserInfo) {
	}

	testMessageDispatch(handleMessage, handleUserChange)
	if received == nil || received.Text != "test event" {
		t.Fatal("Did not receive expected message")
	}
}

func TestDispatcherSendsMessagesOnlyOnce(t *testing.T) {
	calls := 0
	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		calls++
		return q, nil
	}
	handleUserChange := func(event guac.UserInfo) {
	}

	testMessageDispatch(handleMessage, handleUserChange)
	if calls != 1 {
		t.Fatalf("Expected handler to be called exactly once, was called %d times", calls)
	}
}

func TestDispatcherDoesNotSendMessageToUserChangeHandler(t *testing.T) {
	calls := 0
	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		return q, nil
	}
	handleUserChange := func(event guac.UserInfo) {
		calls++
	}

	testMessageDispatch(handleMessage, handleUserChange)
	if calls != 0 {
		t.Fatalf("Expected handler not to be called, was called %d times", calls)
	}
}

func TestDispatcherReturnsErrorIfMessageFails(t *testing.T) {
	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		return q, fmt.Errorf("Error!")
	}
	handleUserChange := func(event guac.UserInfo) {
	}

	err := testMessageDispatch(handleMessage, handleUserChange)
	if err == nil {
		t.Fatal("Expected error ", err)
	}
}

func testUserChangeDispatch(handleMessage MessageHandler, handleUserChange UserChangeHandler) error {
	done := make(DoneChan)
	events := make(guac.EventChan)

	go func() {
		events <- guac.UserChangeEvent{
			UserInfo: guac.UserInfo{Name: "test event"},
		}
		close(done)
	}()

	dispatcher := CreateDispatcher(queue.Queue{}, 1*time.Millisecond, handleMessage, handleUserChange)
	return dispatcher(events, done)
}

func TestDispatcherSendsUserChangesToUserChangeHandler(t *testing.T) {
	var received *guac.UserInfo
	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		return q, nil
	}
	handleUserChange := func(event guac.UserInfo) {
		received = &event
	}

	testUserChangeDispatch(handleMessage, handleUserChange)
	if received == nil || received.Name != "test event" {
		t.Fatal("Did not receive expected message")
	}
}

func TestDispatcherDoesNotSendUserChangeToMessageHandler(t *testing.T) {
	calls := 0
	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		calls++
		return q, nil
	}
	handleUserChange := func(event guac.UserInfo) {
	}

	testUserChangeDispatch(handleMessage, handleUserChange)
	if calls != 0 {
		t.Fatalf("Expected handler not to be called, was called %d times", calls)
	}
}

func TestDispatcherSendsUserChangeOnlyOnce(t *testing.T) {
	calls := 0
	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		return q, nil
	}
	handleUserChange := func(event guac.UserInfo) {
		calls++
	}

	testUserChangeDispatch(handleMessage, handleUserChange)
	if calls != 1 {
		t.Fatalf("Expected handler to be called exactly once, was called %d times", calls)
	}
}

func TestDispatcherReturnsErrorOnTimeout(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)

	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		return q, nil
	}
	handleUserChange := func(event guac.UserInfo) {
	}
	dispatcher := CreateDispatcher(queue.Queue{}, 1*time.Millisecond, handleMessage, handleUserChange)

	err := dispatcher(events, done)
	if err == nil {
		t.Fatal("Expected error ", err)
	}
}

func TestDispatcherReturnsNoErrorIfDone(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)

	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		return q, nil
	}
	handleUserChange := func(event guac.UserInfo) {
	}
	dispatcher := CreateDispatcher(queue.Queue{}, 1*time.Millisecond, handleMessage, handleUserChange)

	close(done)
	err := dispatcher(events, done)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}
}

func TestDispatcherSwallowsUnknownEvents(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)

	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to MessageHandler")
		return q, nil
	}
	handleUserChange := func(event guac.UserInfo) {
		t.Fatal("Unexpected call to MessageHandler")
	}
	dispatcher := CreateDispatcher(queue.Queue{}, 1*time.Second, handleMessage, handleUserChange)

	// events is blocking so these things must be read in sequence
	go func() {
		events <- "unknown event"
		close(done)
	}()
	dispatcher(events, done)

	select {
	case <-events:
		t.Fatal("Expected event to have been swallowed")
	default:
	}
}

func TestDispatcherPassesUpdatedQueueToMessageHandler(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)
	expectedQueue := queue.Queue([]queue.Item{queue.Item{ID: "U123", Reason: "Tomato"}})
	var receivedQueue queue.Queue

	called := false
	handleMessage := func(q queue.Queue, event guac.MessageEvent) (queue.Queue, error) {
		if !called {
			called = true
			return expectedQueue, nil
		}
		receivedQueue = q
		return q, nil
	}
	handleUserChange := func(event guac.UserInfo) {
		t.Fatal("Unexpected call to MessageHandler")
	}
	dispatcher := CreateDispatcher(queue.Queue{}, 1*time.Second, handleMessage, handleUserChange)

	// events is blocking so these things must be read in sequence
	go func() {
		events <- guac.MessageEvent{
			Text: "test event",
		}
		events <- guac.MessageEvent{
			Text: "a second event",
		}
		close(done)
	}()
	dispatcher(events, done)

	if !receivedQueue.Equal(expectedQueue) {
		t.Fatal("Unexpected queue received on second call", expectedQueue, receivedQueue)
	}
}
