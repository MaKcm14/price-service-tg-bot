package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// setter defines the logic of the part of setting the request for various modes.
type setter interface {
	// mode defines the logic of the mode's instructions and actions.
	mode(chatID int64)

	// productSetter defines the logic of the setting the products for the search.
	productSetter(chatID int64)

	// setRequest defines the logic of setting the request.
	setRequest(update *tgbotapi.Update)

	// showRequest defines the view of the products' request for the current user.
	showRequest(chatID int64)
}

// searcher defines the logic of searching part for various modes.
type searcher interface {
	// startSearch defines the logic of the search the products with the set request's params.
	startSearch(chatID int64)
}
