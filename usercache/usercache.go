package usercache

import (
	"sync"

	"github.com/doozr/guac"
)

// UserCache is a simple cache of usernames and their IDs
type UserCache interface {
	GetUserName(string) string
	GetUserID(string) string
	UpdateUserName(string, string)
}

// userCache contains a mutex controlled list of user info objects keyed on ID
type userCache struct {
	Mux       sync.Mutex
	UserNames map[string]string
}

// New creates an instance of UserCache
func New(users []guac.UserInfo) UserCache {
	uc := &userCache{}
	uc.UserNames = make(map[string]string)
	for _, user := range users {
		uc.UserNames[user.ID] = user.Name
	}
	return uc
}

// GetUserName looks up the username associated with an ID
func (u *userCache) GetUserName(id string) (username string) {
	u.Mux.Lock()
	if val, ok := u.UserNames[id]; ok {
		username = val
	}
	u.Mux.Unlock()
	return
}

// UpdateUserName updates the username associated with an ID
func (u *userCache) UpdateUserName(id string, name string) {
	u.Mux.Lock()
	u.UserNames[id] = name
	u.Mux.Unlock()
}

// GetUserID gets the ID associated with a username
func (u *userCache) GetUserID(name string) (id string) {
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
