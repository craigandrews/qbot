package command_test

import (
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestSuccess(t *testing.T) {
	cmd := command.New(id, name, userCache)
	testCommand(t, cmd.Success, []CommandTest{
		{
			test:             "drop token",
			startQueue:       queue.Queue([]queue.Item{{ID: "U123", Reason: "Banana"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			expectedQueue:    queue.Queue{},
			expectedResponse: "Received a success notification from <@U789|andrew>\n<@U123|craig> (Banana) has finished with the token\nThe token is up for grabs",
		},
		{
			test:             "drop token and give it to the next in line",
			startQueue:       queue.Queue([]queue.Item{{ID: "U123", Reason: "Banana"}, {ID: "U456", Reason: "Next up"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			expectedQueue:    queue.Queue([]queue.Item{{ID: "U456", Reason: "Next up"}}),
			expectedResponse: "Received a success notification from <@U789|andrew>\n<@U123|craig> (Banana) has finished with the token\n*<@U456|edward> (Next up) now has the token*",
		},
		{
			test:             "does nothing if the queue is empty",
			startQueue:       queue.Queue{},
			channel:          "C1A2B3C",
			user:             "U789",
			expectedQueue:    queue.Queue{},
			expectedResponse: "Received a success notification from <@U789|andrew>",
		},
	})
}
