package command_test

import (
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestBarge(t *testing.T) {
	cmd := command.New(name, userCache)
	testCommand(t, cmd.Barge, []CommandTest{
		{
			test:             "make active if queue empty",
			startQueue:       queue.Queue{},
			channel:          "C1A2B3C",
			user:             "U123",
			reason:           "Banana",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}}),
			expectedResponse: "*<@U123|craig> (Banana) now has the token*",
		},
		{
			test:             "add in 2nd place if queue not empty",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			reason:           "Barging",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U789", "Barging"}, {"U456", "Next up"}}),
			expectedResponse: "<@U789|andrew> (Barging) barged to the front\n<@U123|craig> (Banana) still has the token",
		},
		{
			test:             "move to 2nd place if already lower in queue",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}, {"U789", "Needs barge"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			reason:           "Needs barge",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U789", "Needs barge"}, {"U456", "Next up"}}),
			expectedResponse: "<@U789|andrew> (Needs barge) barged to the front\n<@U123|craig> (Banana) still has the token",
		},
	})
}
