package tgbot

import (
	"log/slog"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
	"github.com/MaKcm14/price-service/pkg/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	bot    *tgbotapi.BotAPI
	logger *slog.Logger

	userLastAction map[int64]string
	userSamplePtr  map[int64]map[string]int
	userSample     map[int64]map[string]entities.ProductSample

	userRequest map[int64]dto.ProductRequest

	userInteractor services.UserActions
	api            services.ApiInteractor
}

func New(token string, logger *slog.Logger, interactor services.UserActions, api services.ApiInteractor) (TgBot, error) {
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
		userRequest:    make(map[int64]dto.ProductRequest, 10000),
		userSamplePtr:  make(map[int64]map[string]int, 10000),
		userSample:     make(map[int64]map[string]entities.ProductSample),
		userInteractor: interactor,
		api:            api,
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
				t.start(&update)

			case menuAction:
				t.menu(update.Message.Chat.ID)

			default:
				if action, flagExist := t.userLastAction[update.Message.Chat.ID]; flagExist &&
					action == productSetter {
					t.setQuery(&update)
					t.showRequest(&update)
				}
			}

		} else if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case menuAction:
				t.menu(update.CallbackQuery.From.ID)

			case bestPriceModeData:
				t.bestPriceMode(&update)

			case marketSetterMode:
				t.marketSetterMode(&update)

			case wildberries, megamarket:
				if action, flagExist := t.userLastAction[update.CallbackQuery.From.ID]; flagExist &&
					action == productsIter {
					t.productsIter(&update, update.CallbackQuery.Data)
					continue
				}

				t.addMarket(&update)

			case productSetter:
				t.productSetter(&update)

			case startSearch:
				t.startSearch(&update)

			case productsIter:
				t.productsIter(&update, "")
			}
		}
	}
}
