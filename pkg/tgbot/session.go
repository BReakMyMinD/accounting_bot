package tgbot

import (
	"time"
)

type userSession struct {
	CreatedAt   time.Time
	LastUpdate  time.Time
	Timeout     float64
	UpdatesChan chan Update
	CloseChan   chan bool
}

func NewUserSession(updatesBufSize int, timeout float64) *userSession {
	shn := userSession{
		time.Now(),
		time.Now(),
		timeout,
		make(chan Update, updatesBufSize),
		make(chan bool),
	}
	return &shn
}

func (session *userSession) handleUpdate(update Update) error {
	//todo
	return nil
}

func (session *userSession) timeoutExceeded() bool {
	if time.Now().Sub(session.LastUpdate).Seconds() >= session.Timeout {
		return true
	}
	return false
}
