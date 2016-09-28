package main

import (
	"log"

	"github.com/doozr/guac"
	"github.com/doozr/qbot/usercache"
)

// UserChangeHandler handles incoming user change events
type UserChangeHandler func(guac.UserInfo)

// NewUserChangeHandler creates a new user change handler
func createUserChangeHandler(userCache *usercache.UserCache) UserChangeHandler {
	return func(userChange guac.UserInfo) {
		oldName := userCache.GetUserName(userChange.ID)
		userCache.UpdateUserName(userChange.ID, userChange.Name)
		if oldName == "" {
			log.Printf("New user %s cached", userChange.Name)
		} else if oldName != userChange.Name {
			log.Printf("User %s renamed to %s", oldName, userChange.Name)
		}
	}
}
