package main_test

import (
	"testing"

	"github.com/doozr/guac"
	. "github.com/doozr/qbot"
	"github.com/doozr/qbot/usercache"
)

func getTestUserChangeHandler() (usercache.UserCache, UserChangeHandler) {
	cache := usercache.New([]guac.UserInfo{})
	return cache, CreateUserChangeHandler(cache)
}

func TestAddsNewUser(t *testing.T) {
	userCache, handler := getTestUserChangeHandler()

	handler(guac.UserInfo{
		ID:   "U1234",
		Name: "Mr Test",
	})

	if userCache.GetUserName("U1234") != "Mr Test" {
		t.Fatal("Expected user to be added")
	}
}

func TestUpdatesExistingUser(t *testing.T) {
	userCache, handler := getTestUserChangeHandler()
	userCache.UpdateUserName("U1234", "Mr Oldname")

	handler(guac.UserInfo{
		ID:   "U1234",
		Name: "Mr Test",
	})

	if userCache.GetUserName("U1234") != "Mr Test" {
		t.Fatal("Expected user to be updated")
	}
}
