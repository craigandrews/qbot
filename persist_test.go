package qbot_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/doozr/qbot"
	"github.com/doozr/qbot/queue"
)

func TestDifferentQueueIsSaved(t *testing.T) {
	var fileWritten string
	var contentWritten []byte
	var permsWritten os.FileMode
	writeFile := func(f string, c []byte, p os.FileMode) error {
		fileWritten = f
		contentWritten = c
		permsWritten = p
		return nil
	}

	persist := CreatePersister(writeFile, "output.json", queue.Queue{})
	persist(queue.Queue([]queue.Item{queue.Item{ID: "U12345", Reason: "A reason"}, queue.Item{ID: "U67890", Reason: "Another reason"}}))

	if fileWritten != "output.json" {
		t.Fatal("Incorrect file written: ", fileWritten)
	}

	if string(contentWritten) != `[{"ID":"U12345","Reason":"A reason"},{"ID":"U67890","Reason":"Another reason"}]` {
		t.Fatal("Incorrect content written: ", contentWritten)
	}

	if permsWritten != 0644 {
		t.Fatal("Incorrect file perms: ", permsWritten)
	}
}

func TestIdenticalQueueIsNotWritten(t *testing.T) {
	calls := 0
	writeFile := func(f string, c []byte, p os.FileMode) error {
		calls++
		return nil
	}

	oq := queue.Queue([]queue.Item{
		queue.Item{ID: "U12345", Reason: "A reason"},
		queue.Item{ID: "U67890", Reason: "Another reason"},
	})
	nq := queue.Queue([]queue.Item{
		queue.Item{ID: "U12345", Reason: "A reason"},
		queue.Item{ID: "U67890", Reason: "Another reason"},
	})

	persist := CreatePersister(writeFile, "output.json", oq)
	persist(oq)
	persist(nq)

	if calls > 1 {
		t.Fatal("Expected 1 call to persist, got ", calls)
	}
}

func TestReturnsErrorIfWriteFails(t *testing.T) {
	writeFile := func(f string, c []byte, p os.FileMode) error {
		return fmt.Errorf("Error!")
	}

	persist := CreatePersister(writeFile, "output.json", queue.Queue{})
	err := persist(queue.Queue([]queue.Item{queue.Item{ID: "U1234", Reason: "A reason"}}))
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestStillSavesAfterFailure(t *testing.T) {
	called := false
	var contentWritten []byte
	writeFile := func(f string, c []byte, p os.FileMode) error {
		if !called {
			called = true
			return fmt.Errorf("Error!")
		}
		contentWritten = c
		return nil
	}

	q := queue.Queue([]queue.Item{
		queue.Item{ID: "U12345", Reason: "A reason"},
		queue.Item{ID: "U67890", Reason: "Another reason"},
	})

	persist := CreatePersister(writeFile, "output.json", queue.Queue{})
	persist(q)
	persist(q)

	if string(contentWritten) != `[{"ID":"U12345","Reason":"A reason"},{"ID":"U67890","Reason":"Another reason"}]` {
		t.Fatal("Incorrect content written: ", contentWritten)
	}
}
