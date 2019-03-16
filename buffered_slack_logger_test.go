package goslacklog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestThatCloseWillWaitUntilMessagesAreSend(t *testing.T) {
	called := false
	messageSender := &mockSender{
		sendMessage: func(msg string) error {
			called = true
			assert.Equal(t, "```a```\n```b```\n", msg)
			return nil
		},
	}

	buffered := NewBufferedSlackLoggerWithMessageSender(messageSender, 5, 2, 1*time.Second)
	defer buffered.Close()
	buffered.Write([]byte("a"))
	buffered.Write([]byte("b"))
	<-time.After(2 * time.Second)
	assert.True(t, called)
}

func TestThatMoreOverflowOfMessagesWillGetDropped(t *testing.T) {
	called := false
	messageSender := &mockSender{
		sendMessage: func(msg string) error {
			called = true
			assert.Equal(t, "```a```\n```b```\n", msg)
			return nil
		},
	}

	buffered := NewBufferedSlackLoggerWithMessageSender(messageSender, 3, 2, 1*time.Second)
	buffered.Write([]byte("a"))
	buffered.Write([]byte("b"))
	buffered.Write([]byte("c")) // this will not get logged
	<-time.After(1*time.Second + 300*time.Millisecond)
	buffered.Close()
	assert.True(t, called)
}

func TestThatCallingWriteAfterCloseDoesNotLogIt(t *testing.T) {
	calledCount := 0
	messageSender := &mockSender{
		sendMessage: func(msg string) error {
			calledCount++
			assert.Equal(t, "```a```\n```b```\n", msg)
			return nil
		},
	}

	buffered := NewBufferedSlackLoggerWithMessageSender(messageSender, 3, 2, 1*time.Second)
	buffered.Write([]byte("a"))
	buffered.Write([]byte("b"))
	<-time.After(1*time.Second + 300*time.Millisecond)
	buffered.Close()
	buffered.Write([]byte("c")) // this will not get logged
	<-time.After(1*time.Second + 300*time.Millisecond)
	assert.Equal(t, 1, calledCount)
}
