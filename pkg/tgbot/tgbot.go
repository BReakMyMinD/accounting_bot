package tgbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Tgbot struct {
	token          string
	updatesOffset  int
	httpTimeout    int     //seconds
	sessionTimeout float64 //seconds
	active         bool
	initialState   State
	sessions       map[int64]*userSession
	errorChan      chan error
}

func NewBot(token string, httpTimeout float64, sessionTimeout float64, initialState State) (*Tgbot, error) {
	if initialState == nil {
		return nil, fmt.Errorf("initial step is not set")
	}
	bot := Tgbot{
		token,
		0,
		int(httpTimeout),
		sessionTimeout,
		true,
		nil,
		make(map[int64]*userSession),
		make(chan error, errorLimit),
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

func (bot *Tgbot) StartGettingUpdates() {
	for {
		if !bot.active {
			//if we stop the bot, we must ensure that all session handlers had finished
			for sessionId := range bot.sessions {
				bot.closeSession(sessionId)
			}
			close(bot.errorChan)
			break
		}
		requestBody := RequestUpdates{
			bot.updatesOffset,
			updatesLimit,
			bot.httpTimeout,
			[]string{"message", "callback_query"},
		}
		requestBodyJson, err := json.Marshal(requestBody)
		if err != nil {
			bot.errorChan <- err
			continue
		}
		apiResponse, err := bot.makeApiRequest(bot.prepareApiUrl("getUpdates", ""),
			"POST",
			"application/json",
			requestBodyJson)
		if err != nil {
			bot.errorChan <- err
			time.Sleep(time.Second * 3)
			continue
		}
		var updates []Update
		err = json.Unmarshal(apiResponse.Result, &updates)
		if err != nil {
			bot.errorChan <- err
			continue
		}
		for _, update := range updates {
			if update.UpdateID >= bot.updatesOffset {
				bot.updatesOffset = update.UpdateID + 1
				err = bot.mapSession(update)
				if err != nil {
					bot.errorChan <- err
				}
			}
		}
		for sessionId, session := range bot.sessions {
			//if the session did not receive any updates in last "userSession.Timeout" seconds, we close it
			if session.timeoutExceeded() {
				bot.closeSession(sessionId)
			}
		}
	}
}

func (bot *Tgbot) mapSession(update Update) error {
	var sessionId int64
	switch {
	case update.Message.Sender.ID != 0 && update.Message != nil:
		sessionId = update.Message.Sender.ID
	case update.CallbackQuery.Sender.ID != 0 && update.CallbackQuery != nil:
		sessionId = update.CallbackQuery.Sender.ID
	default:
		return fmt.Errorf("unable to recognize sender of update %d", update.UpdateID)
	}
	if session, ok := bot.sessions[sessionId]; ok {
		session.updatesChan <- update
		session.lastUpdate = time.Now()
	} else {
		session = newUserSession(bot.sessionTimeout)
		session.updatesChan <- update
		bot.sessions[sessionId] = session
		go func() {
			for update := range session.updatesChan {
				var err error
				var nextState State

				if nextState, err = session.state.AfterInput(bot, session.dataBuffer, update); err == nil && nextState != nil {
					session.state = nextState
					err = session.state.BeforeInput(bot, session.dataBuffer, update)
				}

				if err != nil {
					bot.errorChan <- err
					session.state = bot.initialState
					session.dataBuffer = nil
				}
			}
			session.closeChan <- true
		}()
	}
	return nil
}

func (bot *Tgbot) closeSession(sessionId int64) {
	if session, ok := bot.sessions[sessionId]; ok {
		close(session.updatesChan)
		if <-session.closeChan {
			//wait for session goroutine has finished
			delete(bot.sessions, sessionId)
		}
	}
}

func (bot *Tgbot) StopGettingUpdates() {
	bot.active = false
}

func (bot *Tgbot) makeApiRequest(url string, httpMethod string, contentType string, body []byte) (*ApiResponse, error) {
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

func (bot *Tgbot) prepareApiUrl(apiMethod string, filePath string) string {
	var url string
	if apiMethod != "" {
		url = strings.Replace(apiMethodTemplate, "<METHOD>", apiMethod, 1)
	} else {
		url = strings.Replace(apiFileTemplate, "<PATH>", filePath, 1)
	}
	url = strings.Replace(url, "<TOKEN>", bot.token, 1)
	return url
}

func (bot *Tgbot) SendMessage(message Message) error {
	//do generic invalid input response
	return nil
}

func (bot *Tgbot) GetErrorChan() <-chan error {
	return bot.errorChan
}
