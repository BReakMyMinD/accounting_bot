package main

import( "config" 
		"tgbot"
	  	"log"
		"os"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)

	infoLog.Println("starting bot...")

	configReader, err := config.NewConfigReader("../config.json") //for debug, by default "config.json"
	if err != nil{
		errorLog.Fatalln(err)
	}

	
	bot, err := tgbot.NewBot(configReader)
	if err != nil {
		errorLog.Fatalln(err)
	}
	go bot.StartGettingUpdates(errorLog)
	for update := range bot.UpdatesChan {
		infoLog.Println(update.UpdateID)
	}
}
