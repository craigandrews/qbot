package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/doozr/guac"
)

func TestReceiveRunsReceiverInBackground(t *testing.T) {
	receiver := func(events guac.EventChan, done DoneChan) error {
		events <- "test"
		return nil
	}
	done := make(DoneChan)
	waitGroup := sync.WaitGroup{}

	events := receive(receiver, done, &waitGroup)
	select {
	case e := <-events:
		if e.(string) != "test" {
			t.Fatal("Expected test event")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Expected event within 2 seconds")
	}

	waitGroup.Wait()
}

func TestReceiveShutDownCleanlyWithErrors(t *testing.T) {
	receiver := func(events guac.EventChan, done DoneChan) error {
		events <- "test"
		return fmt.Errorf("Error!")
	}
	done := make(DoneChan)
	waitGroup := sync.WaitGroup{}

	events := receive(receiver, done, &waitGroup)
	select {
	case e := <-events:
		if e.(string) != "test" {
			t.Fatal("Expected test event")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Expected event within 2 seconds")
	}

	waitGroup.Wait()
}

type TestRealTimeReceiver struct {
	receive func() (interface{}, error)
}

func (r TestRealTimeReceiver) Receive() (interface{}, error) {
	return r.receive()
}

func TestReceiverPushesEventsToChannel(t *testing.T) {
	client := TestRealTimeReceiver{
		receive: func() (interface{}, error) {
			return "test event", nil
		},
	}
	receiver := createEventReceiver(client)

	events := make(guac.EventChan)
	done := make(DoneChan)
	go receiver(events, done)

	select {
	case e := <-events:
		if e.(string) != "test event" {
			t.Fatal("Expected test event ", e)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Expected event within 2 seconds")
	}

	close(done)
}

func TestReceiverQuitsOnError(t *testing.T) {
	client := TestRealTimeReceiver{
		receive: func() (interface{}, error) {
			return nil, fmt.Errorf("Error!")
		},
	}
	receiver := createEventReceiver(client)

	events := make(guac.EventChan)
	done := make(DoneChan)
	err := receiver(events, done)

	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestReceiverReturnsErrorOnNilEvent(t *testing.T) {
	client := TestRealTimeReceiver{
		receive: func() (interface{}, error) {
			return nil, nil
		},
	}
	receiver := createEventReceiver(client)

	events := make(guac.EventChan)
	done := make(DoneChan)
	err := receiver(events, done)

	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestReceiverShutsDownWhenDoneClosed(t *testing.T) {
	client := TestRealTimeReceiver{
		receive: func() (interface{}, error) {
			return "test event", nil
		},
	}
	receiver := createEventReceiver(client)

	events := make(guac.EventChan)
	done := make(DoneChan)
	close(done)
	err := receiver(events, done)

	if err != nil {
		t.Fatal("Expected no error ", err)
	}

	select {
	case e := <-events:
		t.Fatal("Expected nothing on queue ", e)
	default:
	}
}

func TestReceiverReturnsNoErrorWhenDoneClosed(t *testing.T) {
	client := TestRealTimeReceiver{
		receive: func() (interface{}, error) {
			return nil, fmt.Errorf("Error!")
		},
	}
	receiver := createEventReceiver(client)

	events := make(guac.EventChan)
	done := make(DoneChan)
	close(done)
	err := receiver(events, done)

	if err != nil {
		t.Fatal("Expected no error ", err)
	}

	select {
	case e := <-events:
		t.Fatal("Expected nothing on queue ", e)
	default:
	}
}
