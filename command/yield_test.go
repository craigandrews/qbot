package command_test

import (
	"testing"

	. "github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestYield(t *testing.T) {
	command := New(name, userCache)
	testCommand(t, command.Yield, []CommandTest{
		{
			test:             "do not yield if nobody can receive it",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}}),
			expectedResponse: "<@U123|craig> You cannot yield if there is nobody waiting",
		},
		{
			test:             "yield token and give it to the next in line",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "Next up"}, {"U123", "Banana"}}),
			expectedResponse: "<@U123|craig> (Banana) has yielded the token\n*<@U456|edward> (Next up) now has the token*",
		},
		{
			test:             "warns if user does not have the token",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			channel:          "C1A2B3C",
			user:             "U456",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			expectedResponse: "<@U456|edward> You cannot yield if you do not have the token",
		},
		{
			test:             "does nothing if the queue is empty",
			startQueue:       queue.Queue{},
			channel:          "C1A2B3C",
			user:             "U456",
			expectedQueue:    queue.Queue{},
			expectedResponse: "",
		},
	})
}
