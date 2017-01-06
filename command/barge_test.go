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
			args:             "Banana",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}}),
			expectedResponse: "*<@U123|craig> (Banana) now has the token*",
		},
		{
			test:             "add in 2nd place if queue not empty",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "Barging",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U789", "Barging"}, {"U456", "Next up"}}),
			expectedResponse: "<@U789|andrew> (Barging) barged to the front\n<@U123|craig> (Banana) still has the token",
		},
		{
			test:             "move to 2nd place if already lower in queue",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}, {"U789", "Needs barge"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "Needs barge",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U789", "Needs barge"}, {"U456", "Next up"}}),
			expectedResponse: "<@U789|andrew> (Needs barge) barged to the front\n<@U123|craig> (Banana) still has the token",
		},
		{
			test:             "error if no reason and not in queue",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			expectedResponse: "<@U789|andrew> You must provide a reason for joining",
		},
		{
			test:             "move to 2nd place if already lower in queue and no reason specified",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}, {"U789", "Needs barge"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U789", "Needs barge"}, {"U456", "Next up"}}),
			expectedResponse: "<@U789|andrew> (Needs barge) barged to the front\n<@U123|craig> (Banana) still has the token",
		},
		{
			test:             "move to highest to 2nd place if multiple lower in queue and no reason specified",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}, {"U789", "Needs barge"}, {"U789", "No barge"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U789", "Needs barge"}, {"U456", "Next up"}, {"U789", "No barge"}}),
			expectedResponse: "<@U789|andrew> (Needs barge) barged to the front\n<@U123|craig> (Banana) still has the token",
		},
		{
			test:             "leave active if already active",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			args:             "Banana",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			expectedResponse: "*<@U123|craig> (Banana) now has the token*",
		},
		{
			test:             "leave active if already active and no reason specified",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			args:             "",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			expectedResponse: "*<@U123|craig> (Banana) now has the token*",
		},
		{
			test:             "leave active if already active and no reason specified and also lower in the queue",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}, {"U123", "No barge"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			args:             "",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}, {"U123", "No barge"}}),
			expectedResponse: "*<@U123|craig> (Banana) now has the token*",
		},
	})
}
