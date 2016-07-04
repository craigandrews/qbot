package slack

import (
	"golang.org/x/net/websocket"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"net/url"
	"sync/atomic"
)

type Slack struct {
	Name string
	Token string
	WebSocket *websocket.Conn
	Id string
}

type SlackError struct {
	msg string
}

type ResponseSelf struct {
	Id string `json:"id"`
}

type ResponseRtmStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	Url   string       `json:"url"`
	Self  ResponseSelf `json:"self"`
}

func get(url string) (response []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API GET '%s' failed with code %d", url, resp.StatusCode)
		return
	}

	response, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func encodeFormData(fields map[string]string) string {
	a := make([]string, len(fields))
	ix := 0
	for k, v := range fields {
		a[ix] = fmt.Sprintf("%s=%s", k, url.QueryEscape(v))
		ix++
	}
	return strings.Join(a, "&")
}

func getWebsocketUrl(token string) (wsurl string, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	body, err := get(url)
	if err != nil {
		return
	}

	var respObj ResponseRtmStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return
	}

	if !respObj.Ok {
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsurl = respObj.Url
	id = respObj.Self.Id
	return
}

// New creates a new Slack instance
func New(name, token string) (slackConn *Slack, err error) {
	wsurl, id, err := getWebsocketUrl(token)
	if err != nil {
		return
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		return
	}

	slackConn = &Slack{name, token, ws, id}
	return
}

type RtmMessage struct {
	Id        uint64 `json:"id"`
	Type      string `json:"type"`
	Channel   string `json:"channel"`
	User      string `json:"user"`
	Text      string `json:"text"`
}

// GetMessage blocks until a message arrives from Slack
func (s *Slack) GetMessage() (m RtmMessage, err error) {
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

type UserInfoResponse struct {
	Ok bool            `json:"ok"`
	Error string       `json:"error"`
	Members []UserInfo `json:"members"`
}

type UserInfo struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

// GetUsername retrieves the username of a Slack user from their Slack ID
func (s *Slack) GetUserList() (users []UserInfo, err error) {
	body := encodeFormData(map[string]string {
		"token": s.Token,
	})

	resp, err := get("https://slack.com/api/users.list?" + body)
	if err != nil {
		return
	}

	var response UserInfoResponse
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
