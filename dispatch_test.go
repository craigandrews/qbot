package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/doozr/guac"
)

func TestDispatchRunsReceiverInBackground(t *testing.T) {
	dispatcher := func(events guac.EventChan, done DoneChan) error {
		return nil
	}
	done := make(DoneChan)
	events := make(guac.EventChan)
	waitGroup := sync.WaitGroup{}

	abort := dispatch(dispatcher, events, done, &waitGroup)
	select {
	case <-abort:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected event within 2 seconds")
	}

	waitGroup.Wait()
}

func TestDispatchShutDownCleanlyWithErrors(t *testing.T) {
	dispatcher := func(events guac.EventChan, done DoneChan) error {
		return fmt.Errorf("Error!")
	}
	done := make(DoneChan)
	events := make(guac.EventChan)
	waitGroup := sync.WaitGroup{}

	abort := dispatch(dispatcher, events, done, &waitGroup)
	select {
	case <-abort:
	case <-time.After(2 * time.Second):
		t.Fatal("Expected event within 2 seconds")
	}

	waitGroup.Wait()
}

func TestDispatcherSendsMessagesToMessageHandler(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)
	var received *guac.MessageEvent

	handleMessage := func(event guac.MessageEvent) error {
		if received != nil {
			t.Fatal("Already received a message ", event)
		}
		received = &event
		close(done)
		return nil
	}
	handleUserChange := func(event guac.UserInfo) error {
		t.Fatal("Unexpected call to UserChangeHandler")
		return nil
	}
	dispatcher := createDispatcher(1*time.Millisecond, handleMessage, handleUserChange)

	go func() {
		events <- guac.MessageEvent{
			Text: "test event",
		}
	}()

	err := dispatcher(events, done)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

	if received == nil || received.Text != "test event" {
		t.Fatal("Did not receive expected message")
	}
}

func TestDispatcherReturnsErrorIfMessageFails(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)

	handleMessage := func(event guac.MessageEvent) error {
		return fmt.Errorf("Error!")
	}
	handleUserChange := func(event guac.UserInfo) error {
		t.Fatal("Unexpected call to UserChangeHandler")
		return nil
	}
	dispatcher := createDispatcher(1*time.Millisecond, handleMessage, handleUserChange)

	go func() {
		events <- guac.MessageEvent{
			Text: "test event",
		}
	}()

	err := dispatcher(events, done)
	if err == nil {
		t.Fatal("Expected error ", err)
	}
}

func TestDispatcherSendsUserChangesToUserChangeHandler(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)
	var received *guac.UserInfo

	handleMessage := func(event guac.MessageEvent) error {
		t.Fatal("Unexpected call to MessageHandler")
		return nil
	}
	handleUserChange := func(event guac.UserInfo) error {
		if received != nil {
			t.Fatal("Already received a user change ", event)
		}
		received = &event
		return nil
	}
	dispatcher := createDispatcher(1*time.Millisecond, handleMessage, handleUserChange)

	go func() {
		events <- guac.UserChangeEvent{
			UserInfo: guac.UserInfo{Name: "test event"},
		}
		close(done)
	}()

	err := dispatcher(events, done)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

	if received == nil || received.Name != "test event" {
		t.Fatal("Did not receive expected message")
	}
}

func TestDispatcherReturnsErrorIfUserChangeFails(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)

	handleMessage := func(event guac.MessageEvent) error {
		t.Fatal("Unexpected call to MessageHandler")
		return nil
	}
	handleUserChange := func(event guac.UserInfo) error {
		return fmt.Errorf("Error!")
	}
	dispatcher := createDispatcher(1*time.Millisecond, handleMessage, handleUserChange)

	go func() {
		events <- guac.UserChangeEvent{
			UserInfo: guac.UserInfo{Name: "test event"},
		}
	}()

	err := dispatcher(events, done)
	if err == nil {
		t.Fatal("Expected error ", err)
	}
}

func TestDispatcherReturnsErrorOnTimeout(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)

	handleMessage := func(event guac.MessageEvent) error {
		t.Fatal("Unexpected call to MessageHandler")
		return nil
	}
	handleUserChange := func(event guac.UserInfo) error {
		t.Fatal("Unexpected call to MessageHandler")
		return nil
	}
	dispatcher := createDispatcher(1*time.Millisecond, handleMessage, handleUserChange)

	err := dispatcher(events, done)
	if err == nil {
		t.Fatal("Expected error ", err)
	}
}

func TestDispatcherReturnsNoErrorIfDone(t *testing.T) {
	done := make(DoneChan)
	events := make(guac.EventChan)

	handleMessage := func(event guac.MessageEvent) error {
		t.Fatal("Unexpected call to MessageHandler")
		return nil
	}
	handleUserChange := func(event guac.UserInfo) error {
		t.Fatal("Unexpected call to MessageHandler")
		return nil
	}
	dispatcher := createDispatcher(1*time.Millisecond, handleMessage, handleUserChange)

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
	handleUserChange := func(event guac.UserInfo) error {
		t.Fatal("Unexpected call to MessageHandler")
		return nil
	}
	dispatcher := createDispatcher(1*time.Second, handleMessage, handleUserChange)

	// events is blocking so these things must be read in sequence
	go func() {
		events <- "unknown event"
		close(done)
	}()
	err := dispatcher(events, done)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

	select {
	case <-events:
		t.Fatal("Expected event to have been swallowed")
	default:
	}
}
