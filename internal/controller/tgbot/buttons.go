package tgbot

import (
	"bytes"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// bestPriceMode is the action on the pressing the best-price button.
func (t *TgBot) bestPriceMode(update *tgbotapi.Update) {
	t.userLastAction[update.CallbackQuery.From.ID] = bestPriceModeData
	t.userRequest[update.CallbackQuery.From.ID] = dto.NewProductRequest(entities.BestPriceMode)

	var priceRangeInstructs = []string{
		"*Ты перешёл в режим поиска Best Price 📊 *\n\n",
		"❓*Как его использовать?*\n",
		"- Необходимо нажать на кнопки тех маркетов, в которых ты хочешь искать\n\n",
		"- Затем необходимо ввести название товара, который ты хочешь найти\n\n",
		"*P.S. название товара должно быть максимально точным для увеличения точности поиска*\n\n",
		"*Давай поищем!* 👇",
	}

	buffer := bytes.Buffer{}

	for _, instruct := range priceRangeInstructs {
		buffer.WriteString(instruct)
	}

	var keyboardMode = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Задать маркеты 🛒", marketSetterMode)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMode

	t.bot.Send(message)
}

// marketSetterMode sets the markets.
func (t *TgBot) marketSetterMode(update *tgbotapi.Update) {
	t.userLastAction[update.CallbackQuery.From.ID] = marketSetterMode

	var keyboardMarketSetter = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Megamarket 🛍️", wildberries)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Wildberries 🌸", megamarket)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Задать товар 📦", productSetter)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, "*Выбери маркеты, в которых будет производиться поиск* 👇")

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMarketSetter

	t.bot.Send(message)
}

// productSetter defines the logic of setting the product's name.
func (t *TgBot) productSetter(update *tgbotapi.Update) {
	t.userLastAction[update.CallbackQuery.From.ID] = productSetter

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID,
		"*Введи точное название товара, по которому будет осуществляться поиск* 📦")

	message.ParseMode = markDown
	t.bot.Send(message)
}

// startSearch defines the logic of searching the products using the finished request.
func (t *TgBot) startSearch(update *tgbotapi.Update) {
	t.userLastAction[update.CallbackQuery.From.ID] = startSearch
}
