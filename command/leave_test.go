package command_test

import (
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestLeave(t *testing.T) {
	cmd := command.New(id, name, userCache)
	testCommand(t, cmd.Leave, []CommandTest{
		{
			test:             "do nothing if not present",
			startQueue:       queue.Queue([]queue.Item{{"U456", "Already here"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			args:             "Banana",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "Already here"}}),
			expectedResponse: "<@U123|craig> No entry with a reason that starts with 'Banana' was found",
		},
		{
			test:             "remove last instance of user if no reason specified",
			startQueue:       queue.Queue([]queue.Item{{"U456", "Already here"}, {"U123", "First"}, {"U123", "Last"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			args:             "",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "Already here"}, {"U123", "First"}}),
			expectedResponse: "<@U123|craig> (Last) has left the queue",
		},
		{
			test:             "remove last instance of user matching prefix",
			startQueue:       queue.Queue([]queue.Item{{"U456", "Already here"}, {"U123", "Fitbit"}, {"U123", "First"}, {"U123", "Last"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			args:             "Fi",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "Already here"}, {"U123", "Fitbit"}, {"U123", "Last"}}),
			expectedResponse: "<@U123|craig> (First) has left the queue",
		},
		{
			test:             "do nothing if reason does not match",
			startQueue:       queue.Queue([]queue.Item{{"U456", "Already here"}, {"U123", "First"}, {"U123", "Last"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			args:             "No",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "Already here"}, {"U123", "First"}, {"U123", "Last"}}),
			expectedResponse: "<@U123|craig> No entry with a reason that starts with 'No' was found",
		},
		{
			test:             "do nothing if reason matches a different user",
			startQueue:       queue.Queue([]queue.Item{{"U456", "Already here"}, {"U789", "First"}, {"U123", "Last"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			args:             "Fi",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "Already here"}, {"U789", "First"}, {"U123", "Last"}}),
			expectedResponse: "<@U123|craig> No entry with a reason that starts with 'Fi' was found",
		},
		{
			test:             "warns if entry to leave is active",
			startQueue:       queue.Queue([]queue.Item{{"U456", "Already here"}}),
			channel:          "C1A2B3C",
			user:             "U456",
			args:             "",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "Already here"}}),
			expectedResponse: "<@U456|edward> You have the token, did you mean `done` or `drop`?",
		},
	})
}
