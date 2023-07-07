package main

import( "config" 
	  	"log"
		"os"
	  )

func main() {
	infoLog := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)

	infoLog.Println("starting bot...")

	configReader, err := config.NewConfigReader("./../config.json")
	if err != nil{
		errorLog.Fatalln(err)
	}

	botToken, err := configReader.GetString("token")
	if err != nil{
		errorLog.Fatalln(err)
	}

	infoLog.Printf("token %q", botToken)
	
}
