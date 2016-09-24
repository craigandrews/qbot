package main

import (
	"fmt"
	"os"
	"testing"

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

	persist := createPersister(writeFile, "output.json", queue.Queue{})
	err := persist(queue.Queue{
		queue.Item{ID: "U12345", Reason: "A reason"},
		queue.Item{ID: "U67890", Reason: "Another reason"},
	})
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}

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
	writeFile := func(f string, c []byte, p os.FileMode) error {
		t.Fatal("Unexpected call to writeFile")
		return nil
	}

	q := queue.Queue{
		queue.Item{ID: "U12345", Reason: "A reason"},
		queue.Item{ID: "U67890", Reason: "Another reason"},
	}

	persist := createPersister(writeFile, "output.json", q)
	err := persist(q)
	if err != nil {
		t.Fatal("Unexpected error ", err)
	}
}

func TestReturnsErrorIfWriteFails(t *testing.T) {
	writeFile := func(f string, c []byte, p os.FileMode) error {
		return fmt.Errorf("Error!")
	}

	persist := createPersister(writeFile, "output.json", queue.Queue{})
	err := persist(queue.Queue{queue.Item{ID: "U1234", Reason: "A reason"}})
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestSavesIfFirstAttemptFails(t *testing.T) {
	calls := 0
	var contentWritten []byte
	writeFile := func(f string, c []byte, p os.FileMode) error {
		if calls == 0 {
			calls++
			return fmt.Errorf("Error!")
		}
		calls++
		contentWritten = c
		return nil
	}

	q := queue.Queue{
		queue.Item{ID: "U12345", Reason: "A reason"},
		queue.Item{ID: "U67890", Reason: "Another reason"},
	}

	persist := createPersister(writeFile, "output.json", queue.Queue{})
	err := persist(q)
	if err == nil {
		t.Fatal("Expected error on first call")
	}

	if contentWritten != nil {
		t.Fatal("Expected no content written, got ", string(contentWritten))
	}

	err = persist(q)
	if err != nil {
		t.Fatal("Unexpected error on second call ", err)
	}

	if calls != 2 {
		t.Fatal("Expected 2 calls, got ", calls)
	}

	if string(contentWritten) != `[{"ID":"U12345","Reason":"A reason"},{"ID":"U67890","Reason":"Another reason"}]` {
		t.Fatal("Incorrect content written: ", contentWritten)
	}
}
