package command_test

import (
	"testing"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/queue"
)

func TestReplace(t *testing.T) {
	cmd := command.New(id, name, userCache)
	testCommand(t, cmd.Replace, []CommandTest{
		{
			test: "replace when active and position owned by self",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Already here"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "1 Banana",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Banana"}}),
			expectedResponse: "*<@U123|craig> (Banana) now has the token*",
		},
		{
			test: "replace when next and position owned by self",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "2 Banana",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Banana"}}),
			expectedResponse: "<@U123|craig> (Banana) is now next in line",
		},
		{
			test: "replace when not active and position owned by self",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U456", Reason: "And I'm second"},
				{ID: "U123", Reason: "Already here"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "3 Banana",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U456", Reason: "And I'm second"},
				{ID: "U123", Reason: "Banana"}}),
			expectedResponse: "<@U123|craig> (Banana) is now 3rd in line",
		},
		{
			test: "replace when position not owned by self",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "1 Banana",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			expectedResponse: "<@U123|craig> Not replacing because <@U456|edward> is 1st in line",
		},
		{
			test: "replace when position lower than 1",
			startQueue: queue.Queue([]queue.Item{{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "0 Banana",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			expectedResponse: "<@U123|craig> That's not a valid position in the queue",
		},
		{
			test: "replace when position greater than queue length",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "3 Banana",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			expectedResponse: "<@U123|craig> That's not a valid position in the queue",
		},
		{
			test: "replace when position not an integer",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "Banana",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			expectedResponse: "<@U123|craig> That's not a valid position in the queue",
		},
		{
			test: "replace when no new reason text specified",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "2",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "I have it"},
				{ID: "U123", Reason: "Already here"}}),
			expectedResponse: "<@U123|craig> You must provide a new reason",
		},
	})
}
