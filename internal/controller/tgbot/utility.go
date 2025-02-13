package tgbot

import (
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	// markDown sets the Markdown parse mode.
	markDown = "Markdown"

	// Consts set the markets for the request
	megamarket  = "megamarket"
	wildberries = "wildberries"

	// Command's name.
	menuAction  = "menu"
	showRequest = "show-request"

	// Button's data name.
	bestPriceModeData = "best-price"
	startSearch       = "start-search"
	favouriteModeData = "favourite-products"
	trackedModeData   = "tracked-products"
	marketSetterMode  = "markets-setter"
	productSetter     = "product-setter"
)

// addMarket adds the market to the request for the current ChatID.
func (t *TgBot) addMarket(update *tgbotapi.Update) {
	request := t.userRequest[update.CallbackQuery.From.ID]

	request.Markets[update.CallbackQuery.Data] = update.CallbackQuery.Data

	t.userRequest[update.CallbackQuery.From.ID] = dto.ProductRequest{
		Mode:    request.Mode,
		Markets: request.Markets,
	}
}

// setQuery sets the product query request for the current ChatID.
func (t *TgBot) setQuery(update *tgbotapi.Update) {
	request := t.userRequest[update.Message.Chat.ID]

	request.Query = update.Message.Text

	t.userRequest[update.Message.Chat.ID] = request
}
