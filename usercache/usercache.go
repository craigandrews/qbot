package usercache

import (
	"sync"

	"github.com/doozr/goslack/apitypes"
)

type UserCache struct {
	Mux       sync.Mutex
	UserNames map[string]string
}

// New creates an instance of UserCache
func New(users []apitypes.UserInfo) *UserCache {
	uc := UserCache{}
	uc.UserNames = make(map[string]string)
	for _, user := range users {
		uc.UserNames[user.ID] = user.Name
	}
	return &uc
}

// GetUserName looks up the username associated with an ID
func (u *UserCache) GetUserName(id string) (username string) {
	u.Mux.Lock()
	if val, ok := u.UserNames[id]; ok {
		username = val
	}
	u.Mux.Unlock()
	return
}

// UpdateUserName updates the username associated with an ID
func (u *UserCache) UpdateUserName(user apitypes.UserInfo) {
	u.Mux.Lock()
	u.UserNames[user.ID] = user.Name
	u.Mux.Unlock()
}

// GetUserId gets the ID associated with a username
func (u *UserCache) GetUserId(name string) (id string) {
	u.Mux.Lock()
	for k, v := range u.UserNames {
		if v == name {
			id = k
			break
		}
	}
	u.Mux.Unlock()
	return
}
