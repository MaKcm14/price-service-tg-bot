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
	api            services.ApiInteractor

	favorite  favoriteMode
	prodsMode productsMode
	set       setter
	search    searcher

	track trackedMode
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
		prodsMode:      newProductsMode(botConf),
		track:          newTrackedMode(botConf, repo, logger, api),
		api:            api,
	}, nil
}

// Run starts the telegram bot.
func (t *TgBot) Run() {
	t.logger.Info("starting the bot begun")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := t.botConf.bot.GetUpdatesChan(updateConfig)

	go t.track.sendAsyncRequests()

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
					t.set.setRequest(&update)
					t.set.showRequest(chatID)
				}
			}

		} else if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID

			switch update.CallbackQuery.Data {
			case menuAction:
				t.menu(chatID)

			case bestPriceModeData, addTrackedProductData:

				if update.CallbackQuery.Data == bestPriceModeData {
					bestPrice := newBestPriceMode(t.logger, t.botConf, t.api)
					t.set = bestPrice
					t.search = bestPrice

				} else if update.CallbackQuery.Data == addTrackedProductData {
					t.set = newTrackedMode(t.botConf, t.favorite.repo, t.logger, t.api)
				}

				t.set.mode(chatID)

			case marketSetterMode:
				t.prodsMode.marketSetterMode(chatID)

			case wildberries, megamarket:
				if user, flagExist := t.botConf.users[chatID]; flagExist &&
					user.lastAction == productsIter {
					t.prodsMode.productsIter(chatID, update.CallbackQuery.Data)
					continue
				}

				t.prodsMode.addMarket(&update)

			case productSetter:
				t.set.productSetter(chatID)

			case startSearch:
				t.search.startSearch(chatID)

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

			case trackedModeData:
				t.track.trackedModeMenu(chatID)

			case deleteTrackedProductData:
				t.track.deleteTrackedProduct(chatID)

			case getTrackedProdMode:
				t.track.getTrackedProduct(chatID)

			}
		}
	}
}
