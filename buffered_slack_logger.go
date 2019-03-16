package goslacklog

import (
	"strings"
	"sync"
	"time"
)

type action int

const (
	exitAction action = iota
	sendMessageAction
)

type bufferedSlackWriter struct {
	messageSender MessageSender
	queue         chan string
	size          int
	exitCh        chan bool
	messagesWg    *sync.WaitGroup

	ticker *time.Ticker

	closed    bool
	closeOnce *sync.Once

	actionCh chan action
}

func NewBufferedSlackLogger(hookURL string, bufferSize int, sendInBatches int, sendInterval time.Duration) *bufferedSlackWriter {
	return NewBufferedSlackLoggerWithMessageSender(
		&singleMessageSlackClient{
			slackPostMessageUrl: hookURL,
		},
		bufferSize,
		sendInBatches,
		sendInterval,
	)
}

func NewBufferedSlackLoggerWithMessageSender(messageSender MessageSender, bufferSize int, sendInBatches int, sendInterval time.Duration) *bufferedSlackWriter {
	buffered := &bufferedSlackWriter{
		messageSender: messageSender,
		queue:         make(chan string, bufferSize),
		size:          sendInBatches,
		closed:        false,
		exitCh:        make(chan bool),
		messagesWg:    &sync.WaitGroup{},
		closeOnce:     &sync.Once{},
		actionCh:      make(chan action),
		ticker:        time.NewTicker(sendInterval),
	}
	go buffered.sendTickAction()
	go buffered.sender()
	return buffered
}

func (s *bufferedSlackWriter) Write(p []byte) (n int, err error) {
	if s.closed {
		return len(p), nil
	}
	s.messagesWg.Add(1)
	defer s.messagesWg.Done()
	msg := string(p)
	select {
	case s.queue <- msg:
	case <-s.exitCh:
	}
	// it shouldn't matter if it didn't write anything so no error will be sent
	return len(p), nil
}

func (s *bufferedSlackWriter) Close() {
	s.closeOnce.Do(func() {
		// stop accepting new writes
		s.closed = true

		// kill ticker so no new events get sent
		s.ticker.Stop()

		// release all the ones that are currently waiting
		close(s.exitCh)

		// wait for them to exit
		s.messagesWg.Wait()

		// tell main loop to exit
		s.actionCh <- exitAction
		close(s.queue)
	})
}

func (s *bufferedSlackWriter) sendTickAction() {
	for {
		_, ok := <-s.ticker.C
		if !ok {
			return
		}
		s.actionCh <- sendMessageAction
	}
}

func (s *bufferedSlackWriter) sender() {
	sb := strings.Builder{}
	for act := range s.actionCh {
		switch act {
		case exitAction:
			return
		case sendMessageAction:
			sb.Reset()
			for msg := range readN(s.queue, s.size) {
				sb.WriteString("```")
				sb.WriteString(msg)
				sb.WriteString("```\n")
			}
			str := sb.String()
			if str != "" {
				s.messageSender.SendMessage(sb.String())
			}
		}
	}
}

func readN(chIn <-chan string, n int) <-chan string {
	chOut := make(chan string)
	go func() {
		for i := 0; i < n; i++ {
			select {
			case str, ok := <-chIn:
				if !ok {
					break
				}
				chOut <- str
			default:
				break
			}
		}
		close(chOut)
	}()

	return chOut
}
