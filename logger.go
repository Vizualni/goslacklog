package goslacklog

import (
	"strings"
)

type Logger struct {
	messageSender MessageSender
}

func NewSlackLogger(hookURL string) *Logger {
	return &Logger{
		messageSender: &singleMessageSlackClient{
			slackPostMessageUrl: hookURL,
		},
	}
}

func (s *Logger) Write(p []byte) (n int, err error) {
	msg := surroundWithBackticks(string(p))
	s.messageSender.SendMessage(msg)
	return len(p), nil
}

func surroundWithBackticks(msg string) string {
	sb := strings.Builder{}
	sb.WriteString("```")
	sb.WriteString(msg)
	sb.WriteString("```")
	return sb.String()
}
