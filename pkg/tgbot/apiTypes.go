package tgbot

import "encoding/json"

type ApiResponse struct {
	Ok bool `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	ErrorCode int `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

type Update struct {
	UpdateID int `json:"update_id"`
	Message *Message `json:"message,omitempty"`
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
}

type Message struct {
	MessageID int `json:"message_id"`
	Sender *User `json:"from,omitempty"`
	Date int `json:"date"`
	Text string `json:"text,omitempty"`
	Entities []MessageEntity `json:"entities,omitempty"`
}

type User struct {
	ID int64 `json:"id"`
	UserName string `json:"username,omitempty"`
}

type MessageEntity struct {
	Type string `json:"type"`
	Offset int `json:"offset"`
	Length int `json:"length"`
}

type CallbackQuery struct {
	ID string `json:"id"`
	Sender *User `json:"from"`
	InlineMessageID string `json:"inline_message_id,omitempty"`
}

type RequestUpdates struct {
	Offset int `json:"offset,omitempty"`
	Limit int `json:"limit,omitempty"`
	Timeout int `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}