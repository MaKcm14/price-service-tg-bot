package tgbot

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/kafka"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services/actor"
)

// TgBot defines the bot's logic.
type TgBot struct {
	logger  *slog.Logger
	botConf *tgBotConfigs

	favorite favoriteMode
	track    trackedMode

	prodsMode productsMode
	bestPrice bestPriceMode

	set    setter
	search searcher
	uinter services.UserConfiger
	api    services.Actor
}

func New(
	token string,
	logger *slog.Logger,
	interactor services.UserConfiger,
	api services.ApiInteractor,
	repo services.Repository,
	trackedProducts chan *kafka.TrackedProduct,
	reader services.Reader,
) (*TgBot, error) {
	logger.Info("initializing the bot begun")

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		logger.Error("error of initializing the tg-bot")
		return &TgBot{}, err
	}

	botConf, err := newTgBotConfigs(bot, api)

	if err != nil {
		logger.Error(fmt.Sprintf("error of configuring the tg-bot: %s", err))
		return nil, fmt.Errorf("error of configuring the tg-bot: %s", err)
	}

	return &TgBot{
		botConf:   botConf,
		logger:    logger,
		uinter:    interactor,
		favorite:  newFavoriteMode(logger, botConf, repo),
		prodsMode: newProductsMode(botConf),
		track:     newTrackedMode(botConf, repo, logger, api, trackedProducts, reader),
		bestPrice: newBestPriceMode(logger, botConf, api),
		api:       actor.NewAPI(logger, repo, api),
	}, nil
}

// Run starts the telegram bot.
func (t *TgBot) Run() {
	t.logger.Info("starting the bot begun")

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := t.botConf.bot.GetUpdatesChan(updateConfig)

	go t.track.readTrackedProducts()
	go t.api.SendTrackedProducts(context.Background())

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
					t.set = &t.bestPrice
					t.search = &t.bestPrice

				} else if update.CallbackQuery.Data == addTrackedProductData {
					t.set = t.track
				}

				t.set.mode(chatID)

			case marketSetterMode:
				t.prodsMode.marketSetterMode(chatID)

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

			default:
				for market := range t.botConf.markets {
					if update.CallbackQuery.Data == market {
						if user, flagExist := t.botConf.users[chatID]; flagExist &&
							user.lastAction == productsIter {
							t.prodsMode.productsIter(chatID, strings.ToLower(update.CallbackQuery.Data))
							continue
						}

						t.prodsMode.addMarket(&update)
					}
				}
			}
		}
	}
}

func (t *TgBot) Close() {
	t.track.close()
}
