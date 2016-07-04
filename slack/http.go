package slack

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"encoding/json"
)

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

