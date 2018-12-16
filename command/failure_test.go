package command_test

import (
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestFailure(t *testing.T) {
	cmd := command.New(id, name, userCache)
	testCommand(t, cmd.Failure, []CommandTest{
		{
			test:             "tell token holder if only one in queue",
			startQueue:       queue.Queue([]queue.Item{{ID: "U123", Reason: "Banana"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "it broke",
			expectedQueue:    queue.Queue([]queue.Item{{ID: "U123", Reason: "Banana"}}),
			expectedResponse: "<@U123|craig> Received a failure notification from <@U789|andrew>: it broke",
		},
		{
			test:             "tell token holder and next in line",
			startQueue:       queue.Queue([]queue.Item{{ID: "U123", Reason: "Banana"}, {ID: "U456", Reason: "Next up"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "it broke",
			expectedQueue:    queue.Queue([]queue.Item{{ID: "U123", Reason: "Banana"}, {ID: "U456", Reason: "Next up"}}),
			expectedResponse: "<@U123|craig> <@U456|edward> Received a failure notification from <@U789|andrew>: it broke",
		},
		{
			test:             "does nothing if the queue is empty",
			startQueue:       queue.Queue{},
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "it broke",
			expectedQueue:    queue.Queue{},
			expectedResponse: "Received a failure notification from <@U789|andrew>: it broke",
		},
	})
}
