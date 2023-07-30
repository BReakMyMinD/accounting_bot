package main

import (
	"fmt"
	"tgbot"
	"unicode/utf16"
)

type InitialState struct {
	ResetStateCommand string
}

func (state InitialState) ValidateInput(update tgbot.Update) error {
	return nil
}

func (state InitialState) Action(bot *tgbot.Tgbot, update tgbot.Update) (tgbot.State, error) {
	//define command
	var command string
	var nextState tgbot.State
	for _, entity := range update.Message.Entities {
		if entity.Type == tgbot.CommandEntity {
			msgText16 := utf16.Encode([]rune(update.Message.Text))
			substrTo := entity.Offset + entity.Length
			if substrTo > len(msgText16) {
				return nil, fmt.Errorf("user %d bad update: text too short", update.Message.Sender.ID)
			}
			command16 := msgText16[entity.Offset:substrTo]
			command = string(utf16.Decode(command16))
		}
	}
	switch command {
	case "/start": //maybe without slash
		//do something
	case "/quit":
		nextState = state
	}
	return nextState, nil
}
