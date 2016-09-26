package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/doozr/guac"
)

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

	event := guac.MessageEvent{
		Type:    "message",
		ID:      1234,
		User:    "U4321",
		Channel: "D1A2B3C",
		Text:    "This is a message",
	}
	director := createMessageDirector("U123", "myname", publicHandler, privateHandler)
	err := director(event)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

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

	event := guac.MessageEvent{
		Type:    "message",
		ID:      1234,
		User:    "U4321",
		Channel: "D1A2B3C",
		Text:    "This is a message",
	}
	director := createMessageDirector("U123", "myname", publicHandler, privateHandler)
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

	event := guac.MessageEvent{
		Type:    "message",
		ID:      1234,
		User:    "U4321",
		Channel: "C1A2B3C",
		Text:    "myname: This is a message",
	}
	director := createMessageDirector("U123", "myname", publicHandler, privateHandler)
	err := director(event)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

	expected := event
	expected.Text = "This is a message"
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

	event := guac.MessageEvent{
		Type:    "message",
		ID:      1234,
		User:    "U4321",
		Channel: "C1A2B3C",
		Text:    "<@U123> This is a message",
	}
	director := createMessageDirector("U123", "myname", publicHandler, privateHandler)
	err := director(event)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

	expected := event
	expected.Text = "This is a message"
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

	event := guac.MessageEvent{
		Type:    "message",
		ID:      1234,
		User:    "U4321",
		Channel: "C1A2B3C",
		Text:    "<@U123> This is a message",
	}
	director := createMessageDirector("U123", "myname", publicHandler, privateHandler)
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

	event := guac.MessageEvent{
		Type:    "message",
		ID:      1234,
		User:    "U4321",
		Channel: "C1A2B3C",
		Text:    "This is a message",
	}
	director := createMessageDirector("U123", "myname", publicHandler, privateHandler)
	err := director(event)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}
}
