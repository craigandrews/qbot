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
			test: "remove last entry if no position provided",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Active"},
				{ID: "U456", Reason: "First"},
				{ID: "U456", Reason: "Last"}}),
			channel: "C1A2B3C",
			user:    "U789",
			args:    "edward",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Active"},
				{ID: "U456", Reason: "First"}}),
			expectedResponse: "<@U789|andrew> booted <@U456|edward> (Last) from the list",
		},
		{
			test: "remove entry at position if provided",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Active"},
				{ID: "U456", Reason: "First"},
				{ID: "U456", Reason: "Fitbit"},
				{ID: "U456", Reason: "Last"}}),
			channel: "C1A2B3C",
			user:    "U789",
			args:    "3 edward",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Active"},
				{ID: "U456", Reason: "First"},
				{ID: "U456", Reason: "Last"}}),
			expectedResponse: "<@U789|andrew> booted <@U456|edward> (Fitbit) from the list",
		},
		{
			test: "advise to use oust if target has the token",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Active"},
				{ID: "U456", Reason: "First"}}),
			channel: "C1A2B3C",
			user:    "U789",
			args:    "craig",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Active"},
				{ID: "U456", Reason: "First"}}),
			expectedResponse: "<@U789|andrew> You must oust the token holder",
		},
		{
			test: "do not boot if target does not own entry",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Active"},
				{ID: "U456", Reason: "First"}}),
			channel: "C1A2B3C",
			user:    "U456",
			args:    "2 andrew",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Active"},
				{ID: "U456", Reason: "First"}}),
			expectedResponse: "<@U456|edward> Not replacing because <@U456|edward> is 2nd in line",
		},
		{
			test: "do not boot if target does not have an entry and no position specified",
			startQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Active"},
				{ID: "U456", Reason: "First"}}),
			channel: "C1A2B3C",
			user:    "U456",
			args:    "andrew",
			expectedQueue: queue.Queue([]queue.Item{
				{ID: "U123", Reason: "Active"},
				{ID: "U456", Reason: "First"}}),
			expectedResponse: "<@U456|edward> No entry for andrew was found",
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
