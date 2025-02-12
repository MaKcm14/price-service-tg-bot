package tgbot

import (
	"log/slog"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	bot    *tgbotapi.BotAPI
	logger *slog.Logger

	userLastAction map[int64]string
	userInteractor services.UserActions
}

func New(token string, logger *slog.Logger, interactor services.UserActions) (TgBot, error) {
	logger.Info("initializing the bot begun")

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		logger.Error("error of initializing the tg-bot")
		return TgBot{}, err
	}

	return TgBot{
		bot:            bot,
		logger:         logger,
		userLastAction: make(map[int64]string, 10000),
		userInteractor: interactor,
	}, nil
}

// Run starts the telegram bot.
func (t *TgBot) Run() {
	t.logger.Info("starting the bot begun")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := t.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case "start":
				if t.userLastAction[update.Message.Chat.ID] == "start" {
					continue
				}

				t.start(&update)
				t.userLastAction[update.Message.Chat.ID] = "start"

			case menu:
				if t.userLastAction[update.CallbackQuery.From.ID] == menu {
					continue
				}

				t.menu(update.Message.Chat.ID)
				t.userLastAction[update.CallbackQuery.From.ID] = menu
			}

		} else if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case menu:
				if t.userLastAction[update.CallbackQuery.From.ID] == menu {
					continue
				}

				t.menu(update.CallbackQuery.From.ID)
				t.userLastAction[update.CallbackQuery.From.ID] = menu
			}
		}
	}
}
