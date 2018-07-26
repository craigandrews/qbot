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
			test: "do nothing if not present",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "Banana",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"}}),
			expectedResponse: "<@U123|craig> No entry was found",
		},
		{
			test: "remove last instance of user if no position specified",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"},
				{ID: "U123", Reason: "First"},
				{ID: "U123", Reason: "Last"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"},
				{ID: "U123", Reason: "First"}}),
			expectedResponse: "<@U123|craig> (Last) has left the queue",
		},
		{
			test: "remove instance at position if specified",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"},
				{ID: "U123", Reason: "Fitbit"},
				{ID: "U123", Reason: "First"},
				{ID: "U123", Reason: "Last"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "3",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"},
				{ID: "U123", Reason: "Fitbit"},
				{ID: "U123", Reason: "Last"}}),
			expectedResponse: "<@U123|craig> (First) has left the queue",
		},
		{
			test: "do nothing if position not found",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"},
				{ID: "U123", Reason: "First"},
				{ID: "U123", Reason: "Last"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "4",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"},
				{ID: "U123", Reason: "First"},
				{ID: "U123", Reason: "Last"}}),
			expectedResponse: "<@U123|craig> That's not a valid position in the queue",
		},
		{
			test: "do nothing if position not owned",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"},
				{ID: "U789", Reason: "First"},
				{ID: "U123", Reason: "Last"}}),
			channel: "C1A2B3C",
			user:    "U123",
			args:    "2",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"},
				{ID: "U789", Reason: "First"},
				{ID: "U123", Reason: "Last"}}),
			expectedResponse: "<@U123|craig> Not replacing because <@U789|andrew> is 2nd in line",
		},
		{
			test: "warns if entry to leave is active",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"}}),
			channel: "C1A2B3C",
			user:    "U456",
			args:    "",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U456", Reason: "Already here"}}),
			expectedResponse: "<@U456|edward> You have the token, did you mean `done` or `drop`?",
		},
	})
}
