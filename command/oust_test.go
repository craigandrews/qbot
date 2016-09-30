package command_test

import (
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestOust(t *testing.T) {
	cmd := command.New(name, userCache)
	testCommand(t, cmd.Oust, []CommandTest{
		{
			test:             "swap active with next in line",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "craig",
			expectedQueue:    queue.Queue([]queue.Item{{"U456", "First"}, {"U123", "Active"}}),
			expectedResponse: "<@U789|andrew> ousted <@U123|craig> (Active)\n*<@U456|edward> (First) now has the token*",
		},
		{
			test:             "remove if nobody waiting",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Active"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "craig",
			expectedQueue:    queue.Queue{},
			expectedResponse: "<@U789|andrew> ousted <@U123|craig> (Active)",
		},
		{
			test:             "warns if target not active",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "edward",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			expectedResponse: "<@U789|andrew> You can only oust the token holder",
		},
		{
			test:             "warns if target not valid",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "banana",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			expectedResponse: "<@U789|andrew> You can only oust the token holder",
		},
		{
			test:             "warns if target missing",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Active"}, {"U456", "First"}}),
			expectedResponse: "<@U789|andrew> You must specify who you want to oust",
		},
		{
			test:             "does nothing if queue empty",
			startQueue:       queue.Queue{},
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "craig",
			expectedQueue:    queue.Queue{},
			expectedResponse: "",
		},
	})
}
