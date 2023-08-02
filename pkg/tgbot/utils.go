package tgbot

import (
	"fmt"
	"unicode/utf16"
)

func ParseCommand(message Message) (string, error) {
	var command string
	for _, entity := range message.Entities {
		if entity.Type == CommandEntity {
			msgText16 := utf16.Encode([]rune(message.Text))
			substrTo := entity.Offset + entity.Length
			if substrTo > len(msgText16) {
				return "", fmt.Errorf("user %d bad update: text too short", message.Sender.ID)
			}
			command16 := msgText16[entity.Offset:substrTo]
			command = string(utf16.Decode(command16))
		}
	}
	return command, nil
}

// func MakeMessage[T any](content T) Message {
// 	switch v := content(type) {
// 	case string:
// 		break
// 	}
// }
