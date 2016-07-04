package slack

import (
	"golang.org/x/net/websocket"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"net/url"
)

type Slack struct {
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

func post(url string, body string) (response []byte, err error) {
	reader := strings.NewReader(body)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", reader)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("API POST '%s' with '%v' failed with code %d", url, body, resp.StatusCode)
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

type IncomingMessage struct {
	Id        uint64 `json:"id"`
	Type      string `json:"type"`
	Channel   string `json:"channel"`
	User      string `json:"user"`
	Text      string `json:"text"`
}

func (s *Slack) GetMessage() (m IncomingMessage, err error) {
	err = websocket.JSON.Receive(s.WebSocket, &m)
	return
}

type PostMessageResponse struct {
	Ok        bool   `json:"ok"`
	Error     string `json:"error"`
}

func (s *Slack) PostMessage(channel, text string) (err error) {
	body := encodeFormData(map[string]string{
		"token": s.Token,
		"channel": channel,
		"text": text,
		"parse": "full",
		"as_user": "true",
	})

	resp, err := post("https://slack.com/api/chat.postMessage", body)
	if err != nil {
		return
	}

	var response PostMessageResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return
	}

	if !response.Ok {
		err = fmt.Errorf("Error posting message: %s", response.Error)
	}
	return
}

type UserInfoResponse struct {
	Ok bool `json:"ok"`
	Error string `json:"error"`
	User UserInfo `json:"user"`
}

type UserInfo struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

func (s *Slack) GetUsername(id string) (username string) {
	body := encodeFormData(map[string]string {
		"token": s.Token,
		"user": id,
	})

	resp, err := get("https://slack.com/api/users.info?" + body)
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

	username = response.User.Name
	return
}