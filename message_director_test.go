package main_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/doozr/guac"
	. "github.com/doozr/qbot"
)

func getTestMessageEvent(user, channel, text string) guac.MessageEvent {
	return guac.MessageEvent{
		Type:    "message",
		ID:      1234,
		User:    user,
		Channel: channel,
		Text:    text,
	}
}

func TestPrivateMessageIsRouted(t *testing.T) {
	var received guac.MessageEvent
	privateHandler := func(m guac.MessageEvent) error {
		received = m
		return nil
	}
	publicHandler := func(m guac.MessageEvent) error {
		t.Fatal("Unexpected call to public handler")
		return nil
	}

	event := getTestMessageEvent("U4321", "D1A2B3C", "This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	director(event)

	if !reflect.DeepEqual(event, received) {
		t.Fatal("Event does not match ", event, received)
	}
}

func TestErrorReturnedWhenPrivateMessageFails(t *testing.T) {
	privateHandler := func(m guac.MessageEvent) error {
		return fmt.Errorf("Error!")
	}
	publicHandler := func(m guac.MessageEvent) error {
		t.Fatal("Unexpected call to public handler")
		return nil
	}

	event := getTestMessageEvent("U4321", "D1A2B3C", "This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	err := director(event)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestPublicMessageWithNameIsRouted(t *testing.T) {
	var received guac.MessageEvent
	privateHandler := func(m guac.MessageEvent) error {
		t.Fatal("Unexpected call to private handler")
		return nil
	}
	publicHandler := func(m guac.MessageEvent) error {
		received = m
		return nil
	}

	event := getTestMessageEvent("U4321", "C1A2B3C", "myname: This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	director(event)

	expected := getTestMessageEvent("U4321", "C1A2B3C", "This is a message")
	if !reflect.DeepEqual(expected, received) {
		t.Fatal("Event does not match ", event, received)
	}
}

func TestPublicMessageWithIDIsRouted(t *testing.T) {
	var received guac.MessageEvent
	privateHandler := func(m guac.MessageEvent) error {
		t.Fatal("Unexpected call to private handler")
		return nil
	}
	publicHandler := func(m guac.MessageEvent) error {
		received = m
		return nil
	}

	event := getTestMessageEvent("U4321", "C1A2B3C", "<@U123> This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	director(event)

	expected := getTestMessageEvent("U4321", "C1A2B3C", "This is a message")
	if !reflect.DeepEqual(expected, received) {
		t.Fatal("Event does not match ", event, received)
	}
}

func TestErrorReturnedIfPublicMessageFailed(t *testing.T) {
	privateHandler := func(m guac.MessageEvent) error {
		t.Fatal("Unexpected call to private handler")
		return nil
	}
	publicHandler := func(m guac.MessageEvent) error {
		return fmt.Errorf("Error!")
	}

	event := getTestMessageEvent("U4321", "C1A2B3C", "<@U123> This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	err := director(event)
	if err == nil {
		t.Fatal("Expected error ", err)
	}
}

func TestPublicMessageWithoutNameOrIDIsNotRouted(t *testing.T) {
	privateHandler := func(m guac.MessageEvent) error {
		t.Fatal("Unexpected call to private handler")
		return nil
	}
	publicHandler := func(m guac.MessageEvent) error {
		t.Fatal("Unexpected call to public handler")
		return nil
	}

	event := getTestMessageEvent("U4321", "C1A2B3C", "This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	director(event)
}
