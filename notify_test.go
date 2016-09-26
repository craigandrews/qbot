package main

import (
	"fmt"
	"testing"
)

func TestNotifySuccess(t *testing.T) {
	var channelTargeted string
	var messageSent string
	openIM := func(c string) (string, error) {
		t.Fatal("Unexpected call to openIM")
		return "", nil
	}
	postMessage := func(c string, m string) error {
		channelTargeted = c
		messageSent = m
		return nil
	}

	notify := createNotifier(openIM, postMessage)
	err := notify(Notification{
		Channel: "C123456",
		Message: "This is a message",
	})
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	if channelTargeted != "C123456" {
		t.Fatal("Unexpected channel: ", channelTargeted)
	}

	if messageSent != "This is a message" {
		t.Fatal("Unexpected message: ", messageSent)
	}
}

func TestNotifyUserSuccess(t *testing.T) {
	var channelTargeted string
	var messageSent string
	openIM := func(c string) (string, error) {
		if c != "U654321" {
			t.Fatal("Unexpected user: ", c)
		}
		return "C123456", nil
	}
	postMessage := func(c string, m string) error {
		channelTargeted = c
		messageSent = m
		return nil
	}

	notify := createNotifier(openIM, postMessage)
	err := notify(Notification{
		Channel: "U654321",
		Message: "This is a message",
	})
	if err != nil {
		t.Fatal("Unexpected error: ", err)
	}

	if channelTargeted != "C123456" {
		t.Fatal("Unexpected channel: ", channelTargeted)
	}

	if messageSent != "This is a message" {
		t.Fatal("Unexpected message: ", messageSent)
	}
}

func TestErrorOnChannelOpenFailure(t *testing.T) {
	openIM := func(c string) (string, error) {
		return "", fmt.Errorf("Error!")
	}
	postMessage := func(c string, m string) error {
		t.Fatal("Unexpected call to postMessage")
		return nil
	}

	notify := createNotifier(openIM, postMessage)
	err := notify(Notification{
		Channel: "U654321",
		Message: "This is a message",
	})
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestErrorOnPostFailure(t *testing.T) {
	openIM := func(c string) (string, error) {
		t.Fatal("Unexpected call to openIM")
		return "", nil
	}
	postMessage := func(c string, m string) error {
		return fmt.Errorf("Error!")
	}

	notify := createNotifier(openIM, postMessage)
	err := notify(Notification{
		Channel: "C123456",
		Message: "This is a message",
	})
	if err == nil {
		t.Fatal("Expected error: ", err)
	}
}
