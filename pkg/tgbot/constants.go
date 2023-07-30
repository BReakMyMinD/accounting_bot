package tgbot

const (
	CommandEntity     string = "bot_command"
	apiMethodTemplate string = "https://api.telegram.org/bot<TOKEN>/<METHOD>"
	apiFileTemplate   string = "https://api.telegram.org/file/bot<TOKEN>/<PATH>"
	updatesLimit      int    = 100
	errorLimit        int    = 1
)
