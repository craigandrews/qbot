package command_test

import (
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestDelegate(t *testing.T) {
	cmd := command.New(name, userCache)
	testCommand(t, cmd.Delegate, []CommandTest{
		{
			test:             "delegate when present",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Apple"}}),
			channel:          "C1A2B3C",
			user:             "U456",
			args:             "andrew",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U789", "Apple"}}),
			expectedResponse: "<@U456|edward> (Apple) has delegated to <@U789|andrew>",
		},
		{
			test:             "delegate when present with reason",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Apple"}, {"U456", "Lemon"}}),
			channel:          "C1A2B3C",
			user:             "U456",
			args:             "andrew Apple",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U789", "Apple"}, {"U456", "Lemon"}}),
			expectedResponse: "<@U456|edward> (Apple) has delegated to <@U789|andrew>",
		},
		{
			test:             "delegate when active",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Apple"}}),
			channel:          "C1A2B3C",
			user:             "U123",
			args:             "andrew",
			expectedQueue:    queue.Queue([]queue.Item{{"U789", "Banana"}, {"U456", "Apple"}}),
			expectedResponse: "<@U123|craig> (Banana) has delegated to <@U789|andrew>\n*<@U789|andrew> (Banana) now has the token*",
		},
		{
			test:             "delegate when not present",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Apple"}}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "andrew",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Apple"}}),
			expectedResponse: "<@U789|andrew> You cannot delegate if you are not in the queue",
		},
		{
			test:             "delegate to invalid user",
			startQueue:       queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			channel:          "C1A2B3C",
			user:             "U456",
			args:             "invalid",
			expectedQueue:    queue.Queue([]queue.Item{{"U123", "Banana"}, {"U456", "Next up"}}),
			expectedResponse: "<@U456|edward> You cannot delegate to invalid because they don't exist",
		},
		{
			test:             "delegate when not present",
			startQueue:       queue.Queue([]queue.Item{}),
			channel:          "C1A2B3C",
			user:             "U789",
			args:             "andrew",
			expectedQueue:    queue.Queue([]queue.Item{}),
			expectedResponse: "<@U789|andrew> You cannot delegate if you are not in the queue",
		},
	})
}
