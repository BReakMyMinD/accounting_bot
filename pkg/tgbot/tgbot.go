package tgbot

import (
	"fmt"
	"net/http"
	"strings"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"
)

const apiMethodTemplate string = "https://api.telegram.org/bot<TOKEN>/<METHOD>"
const apiFileTemplate string = "https://api.telegram.org/file/bot<TOKEN>/<PATH>"
const updatesLimit int = 100

type tgbot struct {
	token string
	updatesOffset int
	timeout int
	UpdatesChan chan Update
}

func NewBot(token string, timeout int) (*tgbot, error) {
	bot := tgbot{
		token, 
		0,
		timeout,
		make(chan Update, updatesLimit),
	}
	_, err := bot.makeApiRequest(bot.prepareApiUrl("getMe", ""),
								 "GET",
								 "",
								 nil)
	if err != nil {
		return nil, fmt.Errorf("bot start failed: %s", err.Error())
	} 
	return &bot, nil
}

func (bot *tgbot) StartGettingUpdates(log *log.Logger) {
	for {
		requestBody := RequestUpdates{
			bot.updatesOffset,
			updatesLimit,
			bot.timeout,
			[]string{"message", "callback_query"},
		}
		requestBodyJson, err := json.Marshal(requestBody)
		if err != nil {
			log.Println(err)
			continue
		}
		apiResponse, err := bot.makeApiRequest(bot.prepareApiUrl("getUpdates", ""),
												    "POST",
												    "application/json", 
													requestBodyJson)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 3)
			continue
		}
		var updates []Update
		err = json.Unmarshal(apiResponse.Result, &updates)
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

func (bot *tgbot) makeApiRequest(url string, httpMethod string, contentType string, body []byte) (*ApiResponse, error) {
	request, err := http.NewRequest(httpMethod, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with http status code %d", response.StatusCode)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var apiResponse ApiResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		return nil, err
	}
	if !apiResponse.Ok {
		return nil, fmt.Errorf("request failed with telegram error code %d", apiResponse.ErrorCode)
	}
	return &apiResponse, nil
}

func (bot *tgbot) prepareApiUrl(apiMethod string, filePath string) string {
	var url string
	if apiMethod != "" {
		url = strings.Replace(apiMethodTemplate, "<METHOD>", apiMethod, 1)
	} else {
		url = strings.Replace(apiFileTemplate, "<PATH>", filePath, 1)
	}
	url = strings.Replace(url, "<TOKEN>", bot.token, 1)
	return url
}