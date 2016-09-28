package usercache_test

import (
	"testing"

	"github.com/doozr/guac"
	. "github.com/doozr/qbot/usercache"
)

func TestAddsNewEntry(t *testing.T) {
	cache := New([]guac.UserInfo{})
	cache.UpdateUserName("test", "Mr Test")
	name := cache.GetUserName("test")
	if name != "Mr Test" {
		t.Fatal("Incorrect name ", name)
	}
}

func TestUpdatesExistingEntry(t *testing.T) {
	cache := New([]guac.UserInfo{guac.UserInfo{ID: "test", Name: "Old Testy"}})
	cache.UpdateUserName("test", "Mr Test")
	name := cache.GetUserName("test")
	if name != "Mr Test" {
		t.Fatal("Incorrect name ", name)
	}
}

func GetsIDFromName(t *testing.T) {
	cache := New([]guac.UserInfo{guac.UserInfo{ID: "test", Name: "Old Testy"}})
	id := cache.GetUserID("Old Testy")
	if id != "test" {
		t.Fatal("Incorrect ID ", id)
	}
}

func GetsEmptyNameIfIDNotFound(t *testing.T) {
	cache := New([]guac.UserInfo{})
	name := cache.GetUserName("not there")
	if name != "" {
		t.Fatal("Expected empty name, got ", name)
	}
}

func GetsEmptyIDIfNameNotFound(t *testing.T) {
	cache := New([]guac.UserInfo{})
	id := cache.GetUserID("Mr Nowhere")
	if id != "" {
		t.Fatal("Expected empty ID, got ", id)
	}
}
