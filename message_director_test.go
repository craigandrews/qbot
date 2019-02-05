package qbot_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/doozr/guac"
	. "github.com/doozr/qbot"
	"github.com/doozr/qbot/queue"
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
	privateHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		received = m
		return q, nil
	}
	publicHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to public handler")
		return q, nil
	}

	event := getTestMessageEvent("U4321", "D1A2B3C", "This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	director(queue.Queue{}, event)

	if !reflect.DeepEqual(event, received) {
		t.Fatal("Event does not match ", event, received)
	}
}

func TestPrivateMessageDoesNotGetNewQueue(t *testing.T) {
	var expected = queue.Queue{}
	privateHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		return queue.Queue([]queue.Item{{ID: "U123", Reason: "Tomato"}}), nil
	}
	publicHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to public handler")
		return q, nil
	}

	event := getTestMessageEvent("U4321", "D1A2B3C", "This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	received, _ := director(expected, event)

	if !received.Equal(expected) {
		t.Fatal("Queue should match original ", expected, received)
	}
}

func TestErrorReturnedWhenPrivateMessageFails(t *testing.T) {
	privateHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		return q, fmt.Errorf("Error!")
	}
	publicHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to public handler")
		return q, nil
	}

	event := getTestMessageEvent("U4321", "D1A2B3C", "This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	_, err := director(queue.Queue{}, event)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestPublicMessageWithNameIsRouted(t *testing.T) {
	var received guac.MessageEvent
	privateHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to private handler")
		return q, nil
	}
	publicHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		received = m
		return q, nil
	}

	event := getTestMessageEvent("U4321", "C1A2B3C", "myname: This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	director(queue.Queue{}, event)

	expected := getTestMessageEvent("U4321", "C1A2B3C", "This is a message")
	if !reflect.DeepEqual(expected, received) {
		t.Fatal("Event does not match ", event, received)
	}
}

func TestPublicMessageWithIDIsRouted(t *testing.T) {
	var received guac.MessageEvent
	privateHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to private handler")
		return q, nil
	}
	publicHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		received = m
		return q, nil
	}

	event := getTestMessageEvent("U4321", "C1A2B3C", "<@U123> This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	director(queue.Queue{}, event)

	expected := getTestMessageEvent("U4321", "C1A2B3C", "This is a message")
	if !reflect.DeepEqual(expected, received) {
		t.Fatal("Event does not match ", event, received)
	}
}

func TestPublicMessageGetsNewQueue(t *testing.T) {
	expected := queue.Queue([]queue.Item{{ID: "U123", Reason: "Tomato"}})
	privateHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to private handler")
		return q, nil
	}
	publicHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		return expected, nil
	}

	event := getTestMessageEvent("U4321", "C1A2B3C", "<@U123> This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	received, _ := director(queue.Queue{}, event)

	if !reflect.DeepEqual(expected, received) {
		t.Fatal("Queue does not match ", event, received)
	}
}

func TestErrorReturnedIfPublicMessageFailed(t *testing.T) {
	privateHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to private handler")
		return q, nil
	}
	publicHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		return q, fmt.Errorf("Error!")
	}

	event := getTestMessageEvent("U4321", "C1A2B3C", "<@U123> This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	_, err := director(queue.Queue{}, event)
	if err == nil {
		t.Fatal("Expected error ", err)
	}
}

func TestPublicMessageWithoutNameOrIDIsNotRouted(t *testing.T) {
	privateHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to private handler")
		return q, nil
	}
	publicHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to public handler")
		return q, nil
	}

	event := getTestMessageEvent("U4321", "C1A2B3C", "This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	director(queue.Queue{}, event)
}

func TestPublicMessageWithoutNameOrIDReturnsSameQueue(t *testing.T) {
	q := queue.Queue([]queue.Item{{ID: "U123", Reason: "Tomato"}})
	privateHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to private handler")
		return q, nil
	}
	publicHandler := func(q queue.Queue, m guac.MessageEvent) (queue.Queue, error) {
		t.Fatal("Unexpected call to public handler")
		return q, nil
	}

	event := getTestMessageEvent("U4321", "C1A2B3C", "This is a message")
	director := CreateMessageDirector("U123", "myname", publicHandler, privateHandler)
	received, _ := director(q, event)

	if !q.Equal(received) {
		t.Fatal("Unexpected queue", q, received)
	}
}
