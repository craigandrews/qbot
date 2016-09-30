package command_test

import (
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestJoin(t *testing.T) {
	cmd := command.New(name, userCache)
	testCommand(t, cmd.Join, []CommandTest{
		{
			test:             "join as active when queue is empty",
			startQueue:       queue.Queue{},
			channel:          "C1A2B3C",
			user:             "U123",
			reason:           "Banana",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}}),
			expectedResponse: "*<@U123|craig> (Banana) now has the token*",
		},
		{
			test:             "join as inactive when queue is not empty",
			startQueue:       queue.Queue([]queue.Item{{"U456", "Already here"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			reason:           "Banana",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "Already here"}, {"U123", "Banana"}}),
			expectedResponse: "<@U123|craig> (Banana) has joined the queue",
		},
		{
			test:             "do nothing when entry already exists",
			startQueue:       queue.Queue([]queue.Item{{"U456", "Already here"}}),
			channel:          "C1A2B3C",
			user:             "U456",
			reason:           "Already here",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "Already here"}}),
			expectedResponse: "",
		},
		{
			test:             "join as inactive when same user exists with different reason",
			startQueue:       queue.Queue([]queue.Item{{"U456", "Already here"}}),
			channel:          "C1A2B3C",
			user:             "U456",
			reason:           "Banana",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "Already here"}, {"U456", "Banana"}}),
			expectedResponse: "<@U456|edward> (Banana) has joined the queue",
		},
	})
}
