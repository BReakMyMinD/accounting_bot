package main

import( "config" 
		"tgbot"
	  	"log"
		"os" //todo log in file, not stdout
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
	if err != nil{
		errorLog.Println(err)
		return err
	}

	var token string
	var timeout float64

	err = configReader.GetParameter("bot_token", &token)
	if err != nil {
		errorLog.Println(err)
		return err
	}
	err = configReader.GetParameter("long_polling_timeout", &timeout) //json number interprets as float64!
	if err != nil {
		errorLog.Println(err)
		return err
	}
	bot, err := tgbot.NewBot(token, int(timeout))
	if err != nil {
		errorLog.Println(err)
		return err
	}
	go bot.StartGettingUpdates(errorLog)
	for update := range bot.UpdatesChan {
		infoLog.Println(update.UpdateID)
	}
	return nil
}
