package tgbot

import (
	"fmt"
	"strings"
	//"net/http"
	"encoding/json"
	"bytes"
	"log"
	"io"
	"time"
	"config"
)

type tgbot struct {
	token string
	apiTemplate string
	updatesOffset int
	timeout int
	UpdatesChan chan Update
}

func NewBot(configReader *config.ConfigReader) (*tgbot, error) {
	var token, apiTemplate string
	var timeout float64

	err := configReader.GetParameter("bot_token", &token)
	if err != nil {
		return nil, err
	}
	err = configReader.GetParameter("api_template", &apiTemplate) 
	if err != nil {
		return nil, err
	}
	err = configReader.GetParameter("long_polling_timeout", &timeout) //json number interprets as float64!
	if err != nil {
		return nil, err
	}

	bot := tgbot{
		token, 
		apiTemplate,
		0,
		int(timeout),
		make(chan Update, 100),
	}
	_, err = bot.makeApiRequest("GET", "getMe", "", nil)
	if err != nil {
		return nil, fmt.Errorf("bot start failed: %s", err.Error())
	} 
	return &bot, nil
}


func (bot *tgbot) StartGettingUpdates(log *log.Logger) {
	for {
		requestBody := RequestUpdates{
			bot.updatesOffset,
			100,
			30,
			[]string{"message", "callback_query"},
		}
		requestBodyJson, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatalln(err)
		}
		responseBody, err := bot.makeApiRequest("POST", "getUpdates", "application/json", bytes.NewReader(requestBodyJson))
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 3)
			continue
		}
		var updates []Update
		err = json.Unmarshal(responseBody, &updates)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, update := range updates {
			if update.UpdateID >= bot.updatesOffset {
				bot.updatesOffset = update.UpdateID + 1
				bot.UpdatesChan <- update
			}
		}
	}
}

func (bot *tgbot) makeApiRequest(httpMethod string, apiMethod string, contentType string, body io.Reader) ([]byte, error) {
	url := strings.Replace(bot.apiTemplate, "<TOKEN>", bot.token, 1)
	url = strings.Replace(url, "<METHOD>", apiMethod, 1)
	//todo
	return []byte{}, nil
}