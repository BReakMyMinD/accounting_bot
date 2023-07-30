package tgbot

import (
	"time"
)

type State interface {
	Action(bot *Tgbot, update Update) (State, error)
	ValidateInput(update Update) error
}

type userSession struct {
	createdAt   time.Time
	lastUpdate  time.Time
	timeout     float64
	updatesChan chan Update
	closeChan   chan bool
	state       State
	dataBuffer  interface{}
}

func (session *userSession) timeoutExceeded() bool {
	if time.Now().Sub(session.lastUpdate).Seconds() > session.timeout {
		return true
	}
	return false
}

func newUserSession(timeout float64) *userSession {
	session := userSession{
		time.Now(),
		time.Now(),
		timeout,
		make(chan Update, updatesLimit),
		make(chan bool),
		nil,
		nil,
	}
	return &session
}
