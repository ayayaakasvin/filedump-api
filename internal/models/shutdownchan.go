package models

import "fmt"

type ShutdownChannel chan string

const ShutdownMessage = "%s:%s" // where first one is origin and second error msg

func NewShutdownChannel() ShutdownChannel {
	return make(ShutdownChannel, 1)
}

func (s *ShutdownChannel) Value() string {
	return <-*s
}

func (s *ShutdownChannel) Send(msg string, args ...any) {
	select {
	case *s <- fmt.Sprintf(msg, args...):
	default:
	}
}