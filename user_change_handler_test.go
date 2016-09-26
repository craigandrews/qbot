package main

import (
	"testing"

	"github.com/doozr/guac"
	"github.com/doozr/qbot/usercache"
)

func TestAddsNewUser(t *testing.T) {
	cache := usercache.New([]guac.UserInfo{})
	handler := createUserChangeHandler(cache)

	err := handler(guac.UserInfo{
		ID:   "U1234",
		Name: "Mr Test",
	})
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if cache.GetUserName("U1234") != "Mr Test" {
		t.Fatal("Expected user to be added")
	}
}

func TestUpdatesExistingUser(t *testing.T) {
	cache := usercache.New([]guac.UserInfo{
		guac.UserInfo{ID: "U1234", Name: "Mr Oldname"},
	})
	handler := createUserChangeHandler(cache)

	err := handler(guac.UserInfo{
		ID:   "U1234",
		Name: "Mr Test",
	})
	if err != nil {
		t.Fatal("Unexpected error")
	}

	if cache.GetUserName("U1234") != "Mr Test" {
		t.Fatal("Expected user to be updated")
	}
}
