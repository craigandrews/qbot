package usercache

import (
	"github.com/doozr/qbot/slack"
	"sync"
)

type UserCache struct {
	Mux sync.Mutex
	UserNames map[string]string
}

func New(users []slack.UserInfo) *UserCache {
	uc := UserCache{}
	uc.UserNames = make(map[string]string)
	for _, user := range users {
		uc.UserNames[user.Id] = user.Name
	}
	return &uc
}

func (u *UserCache) GetUserName(id string) (username string) {
	u.Mux.Lock()
	if val, ok := u.UserNames[id]; ok {
		username = val
	}
	u.Mux.Unlock()
	return
}

func (u *UserCache) UpdateUserName(user slack.UserInfo) {
	u.Mux.Lock()
	u.UserNames[user.Id] = user.Name
	u.Mux.Unlock()
}