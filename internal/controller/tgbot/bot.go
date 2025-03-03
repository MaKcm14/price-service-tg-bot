package tgbot

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
)

// TgBot defines the bot's logic.
type TgBot struct {
	logger  *slog.Logger
	botConf *tgBotConfigs

	userInteractor services.UserConfiger

	favorite  favoriteMode
	prodsMode productsMode
	bestPrice bestPriceMode
}

func New(token string, logger *slog.Logger, interactor services.UserConfiger, api services.ApiInteractor, repo services.Repository) (*TgBot, error) {
	logger.Info("initializing the bot begun")

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		logger.Error("error of initializing the tg-bot")
		return &TgBot{}, err
	}

	botConf := newTgBotConfigs(bot)

	return &TgBot{
		botConf:        botConf,
		logger:         logger,
		userInteractor: interactor,
		favorite:       newFavoriteMode(logger, botConf, repo),
		bestPrice:      newBestPriceMode(logger, botConf, api),
	}, nil
}

// Run starts the telegram bot.
func (t *TgBot) Run() {
	t.logger.Info("starting the bot begun")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := t.botConf.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID

			switch update.Message.Command() {
			case startAction:
				t.start(chatID)

			case menuAction:
				t.menu(chatID)

			default:
				if user, flagExist := t.botConf.users[chatID]; flagExist &&
					user.lastAction == productSetter {
					t.prodsMode.setQuery(&update)
					t.prodsMode.showRequest(chatID)
				}
			}

		} else if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID

			switch update.CallbackQuery.Data {
			case menuAction:
				t.menu(chatID)

			case bestPriceModeData:
				t.bestPrice.bestPriceMode(chatID)

			case marketSetterMode:
				t.bestPrice.marketSetterMode(chatID)

			case wildberries, megamarket:
				if user, flagExist := t.botConf.users[chatID]; flagExist &&
					user.lastAction == productsIter {
					t.prodsMode.productsIter(chatID, update.CallbackQuery.Data)
					continue
				}

				t.prodsMode.addMarket(&update)

			case productSetter:
				t.bestPrice.productSetter(chatID)

			case startSearch:
				t.bestPrice.startSearch(chatID)

			case productsIter:
				t.prodsMode.productsIter(chatID, "")

			case favoriteModeData:
				t.favorite.favoriteMode(chatID)

			case addToFavorite:
				t.favorite.addFavoriteProduct(chatID)

			case showFavoriteProducts:
				t.favorite.showFavoriteProducts(chatID)

			case deleteFavoriteProduct:
				t.favorite.deleteFavoriteProduct(chatID)

			}
		}
	}
}
