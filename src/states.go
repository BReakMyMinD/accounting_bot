package main

import (
	"tgbot"
)

const quitCommand string = "/quit"

type InitialState struct{}

func (state InitialState) BeforeInput(bot *tgbot.Tgbot, dataBuffer interface{}, update tgbot.Update) error {
	return nil
}

func (state InitialState) AfterInput(bot *tgbot.Tgbot, dataBuffer interface{}, update tgbot.Update) (tgbot.State, error) {
	//define command
	var nextState tgbot.State
	//var err error
	command, err := tgbot.ParseCommand(*(update.Message))
	switch command {
	case "/start": //maybe without slash
		//do something
	case "/quit":
		nextState = InitialState{}
	case "/new":
		nextState = NewPayState{}
	default:
		nextState = InitialState{}
	}
	return nextState, err
}

type NewPayState struct{}

func (state NewPayState) BeforeInput(bot *tgbot.Tgbot, dataBuffer interface{}, update tgbot.Update) error {
	//inform user about input format
	return nil
}

func (state NewPayState) AfterInput(bot *tgbot.Tgbot, dataBuffer interface{}, update tgbot.Update) (tgbot.State, error) {
	var err error
	if command, _ := tgbot.ParseCommand(*(update.Message)); command != "" {
		nextState, err := InitialState{}.AfterInput(bot, dataBuffer, update)
		return nextState, err
	}
	if true { //if input is fully validated
		//do some shit
		//return some_new_state nil
	}
	//err = bot.SendMessage("bullshit input")
	return nil, err
}
