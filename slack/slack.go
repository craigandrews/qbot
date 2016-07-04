package slack

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

type Slack struct {
	Token     string
	WebSocket *websocket.Conn
	Id        string
}

type SlackError struct {
	msg string
}

// New creates a new Slack instance
func New(token string) (slackConn *Slack, err error) {
	wsurl, id, err := getWebsocketUrl(token)
	if err != nil {
		return
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		return
	}

	slackConn = &Slack{token, ws, id}
	return
}

// GetMessage blocks until a message arrives from Slack
func (s *Slack) GetEvent() (m RtmEvent, err error) {
	err = websocket.JSON.Receive(s.WebSocket, &m)
	return
}

var counter uint64

// PostMessage sends a message to a Slack channel
func (s *Slack) PostMessage(channel, text string) error {
	id := atomic.AddUint64(&counter, 1)
	m := RtmMessage{id, "message", channel, "", text}
	m.Id = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(s.WebSocket, m)
}

// GetUserList retrieves a list of user IDs mapped to usernames from Slack
func (s *Slack) GetUserList() (users []UserInfo, err error) {
	body := encodeFormData(map[string]string{
		"token": s.Token,
	})

	resp, err := get("https://slack.com/api/users.list?" + body)
	if err != nil {
		return
	}

	var response ResponseUserList
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return
	}

	if !response.Ok {
		err = fmt.Errorf("Error getting user info: %s", response.Error)
	}

	users = response.Members
	return
}
