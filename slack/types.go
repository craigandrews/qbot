package slack

type UserInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
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

type ResponseUserList struct {
	Ok      bool       `json:"ok"`
	Error   string     `json:"error"`
	Members []UserInfo `json:"members"`
}

type RtmEvent struct {
	Id      uint64      `json:"id"`
	Type    string      `json:"type"`
	Channel string      `json:"channel"`
	User    interface{} `json:"user"`
	Text    string      `json:"text"`
}

type RtmMessage struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Text    string `json:"text"`
}

type RtmUserChange struct {
	Type string   `json:"type"`
	User UserInfo `json:"user"`
}

// ConvertEventToMessage casts an RtmEvent to an RtmMessage
func ConvertEventToMessage(e RtmEvent) (msg RtmMessage) {
	user := ""
	if e.User != nil {
		user = e.User.(string)
	}
	msg = RtmMessage{e.Id, e.Type, e.Channel, user, e.Text}
	return
}

// ConvertEventToUserChange casts an RtmEvent to an RtmUserChange
func ConvertEventToUserChange(e RtmEvent) (uc RtmUserChange) {
	user := UserInfo{}
	if e.User != nil {
		ui := e.User.(map[string]interface{})
		user = UserInfo{ui["id"].(string), ui["name"].(string)}
	}
	uc = RtmUserChange{e.Type, user}
	return
}
