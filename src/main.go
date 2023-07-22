package main

import (
	"config"
	"log"
	"os" //todo log in file, not stdout
	"tgbot"
)

func main() {
	if run() != nil {
		os.Exit(1)
	}
}

func run() error {
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)

	infoLog.Println("starting bot...")

	configReader, err := config.NewConfigReader("../config.json") //for debug, by default "config.json"
	if err != nil {
		errorLog.Println(err)
		return err
	}

	var token string
	var httpTimeout, sessionTimeout float64

	err = configReader.GetParameter("bot_token", &token)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	err = configReader.GetParameter("long_polling_timeout", &httpTimeout) //json number interprets as float64!
	if err != nil {
		errorLog.Println(err)
		return err
	}
	err = configReader.GetParameter("session_timeout", &sessionTimeout)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	bot, err := tgbot.NewBot(token, httpTimeout, sessionTimeout)
	if err != nil {
		errorLog.Println(err)
		return err
	}

	go bot.StartGettingUpdates()
	for err := range bot.ErrorChan {
		errorLog.Println(err)
	}
	infoLog.Println("bot terminated")
	return nil
}
