package command_test

import (
	"fmt"
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestDelegate(t *testing.T) {
	cmd := command.New(id, name, userCache)
	testCommand(t, cmd.Delegate, []CommandTest{
		{
			test: "delegate when present",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"}}),
			channel: "C1A2B3C",
			user:    "U456",
			args:    "andrew",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U789", Reason: "Apple"}}),
			expectedResponse: "<@U456|edward> (Apple) has delegated to <@U789|andrew>",
		},
		{
			test: "delegate when present with position",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"},
				{ID: "U456", Reason: "Lemon"}}),
			channel: "C1A2B3C",
			user:    "U456",
			args:    "2 andrew",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U789", Reason: "Apple"},
				{ID: "U456", Reason: "Lemon"}}),
			expectedResponse: "<@U456|edward> (Apple) has delegated to <@U789|andrew>",
		},
		{
			test: "delegate with position that does not belong to you",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"},
				{ID: "U456", Reason: "Lemon"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "2 andrew",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"},
				{ID: "U456", Reason: "Lemon"}}),
			expectedResponse: "<@U123|craig> Not replacing because <@U456|edward> is 2nd in line",
		},
		{
			test: "delegate when active",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "andrew",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U789", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"}}),
			expectedResponse: "<@U123|craig> (Banana) has delegated to <@U789|andrew>\n*<@U789|andrew> (Banana) now has the token*",
		},
		{
			test: "delegate to qbot when inactive",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"}}),
			channel: "C1A2B3C",
			user:    "U456",
			args:    name,
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"}}),
			expectedResponse: "What am I going to do with the token?",
		},
		{
			test: "delegate to qbot when active",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    name,
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"}}),
			expectedResponse: fmt.Sprintf("<@U123|craig> (Banana) has delegated to <@%s|%s>\n*<@%s|%s> (Banana) now has the token*\n:zap: :zap: AT LAST! ULTIMATE POWER! :zap: :zap:\n\nJust kidding ... I don't need the token, you can have it back\n<@U12345|the_bot_name> (Banana) has delegated to <@U123|craig>\n*<@U123|craig> (Banana) now has the token*",
				id, name, id, name),
		},
		{
			test: "delegate when not present",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"}}),
			channel: "C1A2B3C",
			user:    "U789",
			args:    "andrew",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Apple"}}),
			expectedResponse: "<@U789|andrew> You cannot delegate if you are not in the queue",
		},
		{
			test: "delegate to invalid user",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Next up"}}),
			channel: "C1A2B3C",
			user:    "U456",
			args:    "invalid",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"},
				{ID: "U456", Reason: "Next up"}}),
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
