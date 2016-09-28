package main_test

import (
	"sync"
	"testing"
	"time"

	. "github.com/doozr/qbot"
)

func TestSendsMultiplePingsUntilDone(t *testing.T) {
	expectedCalls := 3
	calls := 0
	done := make(DoneChan)
	waitGroup := sync.WaitGroup{}

	pinger := func() error {
		calls++
		return nil
	}

	timeChan := make(chan time.Time)
	close(timeChan)
	after := func(d time.Duration) <-chan time.Time {
		if calls < expectedCalls {
			return timeChan
		}
		close(done)
		return nil
	}

	StartKeepAlive(pinger, after, done, &waitGroup)
	waitGroup.Wait()

	if calls != expectedCalls {
		t.Fatalf("Expected %d calls, got %d", expectedCalls, calls)
	}
}

func TestSendsNoPingsIfDoneBeforeTimeout(t *testing.T) {
	done := make(DoneChan)
	close(done)
	waitGroup := sync.WaitGroup{}

	pinger := func() error {
		t.Fatal("Unexpected call to ping")
		return nil
	}

	after := func(d time.Duration) <-chan time.Time {
		return nil
	}

	StartKeepAlive(pinger, after, done, &waitGroup)
	waitGroup.Wait()
}
