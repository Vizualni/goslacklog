package goslacklog

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type MessageSender interface {
	SendMessage(msg string) error
}

type singleMessageSlackClient struct {
	Client              *http.Client
	slackPostMessageUrl string
}

func (s *singleMessageSlackClient) SendMessage(msg string) error {

	slackMsg := map[string]string{
		"text": msg,
	}

	b, _ := json.Marshal(slackMsg)

	req, _ := http.NewRequest(
		"POST",
		s.slackPostMessageUrl,
		bytes.NewReader(b),
	)

	var c *http.Client
	if s.Client == nil {
		c = http.DefaultClient
	}
	c.Do(req)
	return nil
}
