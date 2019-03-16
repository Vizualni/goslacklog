package goslacklog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSender struct {
	sendMessage func(string) error
}

func (m *mockSender) SendMessage(msg string) error {
	return m.sendMessage(msg)
}

func TestThatLoggerWillCallSendMessage(t *testing.T) {
	called := false
	logger := Logger{
		messageSender: &mockSender{
			sendMessage: func(msg string) error {
				called = true
				return nil
			},
		},
	}
	logger.Write([]byte("called"))
	assert.True(t, called)
}

func TestThatMessageWillContainBackticks(t *testing.T) {
	called := false
	logger := Logger{
		messageSender: &mockSender{
			sendMessage: func(msg string) error {
				called = true
				assert.Equal(t, msg, "```called```")
				return nil
			},
		},
	}
	logger.Write([]byte("called"))
	assert.True(t, called)

}
