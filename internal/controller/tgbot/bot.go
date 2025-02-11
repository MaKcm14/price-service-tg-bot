package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type TgBot struct {
	bot *tgbotapi.BotAPI
}

func New() TgBot {
	return TgBot{}
}
