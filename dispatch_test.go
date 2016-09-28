package qbot_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/doozr/guac"
	. "github.com/doozr/qbot"
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

	dispatcher := CreateDispatcher(1*time.Millisecond, handleMessage, handleUserChange)
	return dispatcher(events, done)
}

func TestDispatcherSendsMessagesToMessageHandler(t *testing.T) {
	var received *guac.MessageEvent
	handleMessage := func(event guac.MessageEvent) error {
		received = &event
		return nil
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
	handleMessage := func(event guac.MessageEvent) error {
		calls++
		return nil
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
	handleMessage := func(event guac.MessageEvent) error {
		return nil
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
	handleMessage := func(event guac.MessageEvent) error {
		return fmt.Errorf("Error!")
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

	dispatcher := CreateDispatcher(1*time.Millisecond, handleMessage, handleUserChange)
	return dispatcher(events, done)
}

func TestDispatcherSendsUserChangesToUserChangeHandler(t *testing.T) {
	var received *guac.UserInfo
	handleMessage := func(event guac.MessageEvent) error {
		return nil
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
	handleMessage := func(event guac.MessageEvent) error {
		calls++
		return nil
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
	handleMessage := func(event guac.MessageEvent) error {
		return nil
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

	handleMessage := func(event guac.MessageEvent) error {
		return nil
	}
	handleUserChange := func(event guac.UserInfo) {
	}
	dispatcher := CreateDispatcher(1*time.Millisecond, handleMessage, handleUserChange)

	err := dispatcher(events, done)
	if err == nil {
		t.Fatal("Expected error ", err)
	}
}

func TestDispatcherReturnsNoErrorIfDone(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)

	handleMessage := func(event guac.MessageEvent) error {
		return nil
	}
	handleUserChange := func(event guac.UserInfo) {
	}
	dispatcher := CreateDispatcher(1*time.Millisecond, handleMessage, handleUserChange)

	close(done)
	err := dispatcher(events, done)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}
}

func TestDispatcherSwallowsUnknownEvents(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)

	handleMessage := func(event guac.MessageEvent) error {
		t.Fatal("Unexpected call to MessageHandler")
		return nil
	}
	handleUserChange := func(event guac.UserInfo) {
		t.Fatal("Unexpected call to MessageHandler")
	}
	dispatcher := CreateDispatcher(1*time.Second, handleMessage, handleUserChange)

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
