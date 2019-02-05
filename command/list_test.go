package command_test

import (
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestList(t *testing.T) {
	cmd := command.New(id, name, userCache)
	testCommand(t, cmd.List, []CommandTest{
		{
			test:             "list all users who are waiting",
			startQueue:       queue.Queue([]queue.Item{{ID: "U123", Reason: "Active"}, {ID: "U456", Reason: "First"}, {ID: "U789", Reason: "Last"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			expectedQueue:    queue.Queue([]queue.Item{{ID: "U123", Reason: "Active"}, {ID: "U456", Reason: "First"}, {ID: "U789", Reason: "Last"}}),
			expectedResponse: "*1: craig (Active) has the token*\n2: edward (First)\n3: andrew (Last)",
		},
		{
			test:             "list active users if nobody waiting",
			startQueue:       queue.Queue([]queue.Item{{ID: "U123", Reason: "Active"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			expectedQueue:    queue.Queue([]queue.Item{{ID: "U123", Reason: "Active"}}),
			expectedResponse: "*1: craig (Active) has the token*",
		},
		{
			test:             "gives friendly message if queue empty",
			startQueue:       queue.Queue{},
			channel:          "C1A2B3C",
			user:             "U789",
			expectedQueue:    queue.Queue{},
			expectedResponse: "Nobody has the token, and nobody is waiting",
		},
	})
}
