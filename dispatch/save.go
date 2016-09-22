package dispatch

import (
	"log"
	"sync"

	"github.com/doozr/jot"
	"github.com/doozr/qbot/usercache"
)

// User handles user renaming in the user cache
func User(userCache *usercache.UserCache, userUpdateChan UserChan, waitGroup *sync.WaitGroup) {

	jot.Print("user dispatch started")
	defer func() {
		waitGroup.Done()
		jot.Print("user dispatch done")
	}()

	for u := range userUpdateChan {
		oldName := userCache.GetUserName(u.ID)
		userCache.UpdateUserName(u.ID, u.Name)
		if oldName == "" {
			log.Printf("New user %s cached", u.Name)
		} else if oldName != u.Name {
			log.Printf("User %s renamed to %s", oldName, u.Name)
		}
	}
}
