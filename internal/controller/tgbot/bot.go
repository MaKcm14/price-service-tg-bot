package tgbot

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ParseMode = "Markdown"
)

type TgBot struct {
	bot    *tgbotapi.BotAPI
	logger *slog.Logger
}

func New(token string, logger *slog.Logger) (TgBot, error) {
	logger.Info("initializing the bot begun")

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		logger.Error("error of initializing the tg-bot")
		return TgBot{}, err
	}

	return TgBot{
		bot:    bot,
		logger: logger,
	}, nil
}

// Run starts the telegram bot.
func (t TgBot) Run() {
	t.logger.Info("starting the bot begun")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := t.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case "start":
				t.start(&update)

			case "menu":
				t.menu(&update)
			}

		} else if update.CallbackQuery != nil {
		}
	}
}
