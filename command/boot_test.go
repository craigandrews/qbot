package command_test

import (
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestBoot(t *testing.T) {
	cmd := command.New(id, name, userCache)
	testCommand(t, cmd.Boot, []CommandTest{
		{
			test:             "remove last entry if no position provided",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}, {"U456", "Last"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "edward",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			expectedResponse: "<@U789|andrew> booted <@U456|edward> (Last) from the list",
		},
		{
			test:             "remove entry at position if provided",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}, {"U456", "Fitbit"}, {"U456", "Last"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "2 edward",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}, {"U456", "Last"}}),
			expectedResponse: "<@U789|andrew> booted <@U456|edward> (Fitbit) from the list",
		},
		{
			test:             "advise to use oust if target has the token",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "craig",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			expectedResponse: "<@U789|andrew> You must oust the token holder",
		},
		{
			test:             "do not boot if target does not own entry",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			channel:          "C1A2B3C",
			user:             "U456",
			args:             "2 andrew",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			expectedResponse: "<@U456|edward> No entry with for andrew with a reason that starts with 'something' was found",
		},
		{
			test:             "do not boot if target does not have an entry and no position specified",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			channel:          "C1A2B3C",
			user:             "U456",
			args:             "andrew",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			expectedResponse: "<@U456|edward> No entry with for andrew with a reason that starts with 'something' was found",
		},
		{
			test:             "do nothing if the queue is empty",
			startQueue:       queue.Queue{},
			channel:          "C1A2B3C",
			user:             "U456",
			args:             "andrew something",
			expectedQueue:    queue.Queue{},
			expectedResponse: "",
		},
	})
}
