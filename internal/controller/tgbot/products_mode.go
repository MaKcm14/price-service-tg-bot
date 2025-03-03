package tgbot

import (
	"bytes"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
)

// productsMode defines the main logic of the products mode processing.
type productsMode struct {
	botConf *tgBotConfigs
}

// nextProduct defines the logic of getting the next product.
func (p *productsMode) nextProduct(chatID int64, market string) {
	p.botConf.users[chatID].sample.lastMarketChoice = market

	count := p.botConf.users[chatID].sample.samplePtr[market]

	if sample := p.botConf.users[chatID].sample.sample[market]; len(sample.Products) <= count {
		return
	}

	sample := p.botConf.users[chatID].sample.sample[market]

	p.botConf.users[chatID].sample.samplePtr[market] = count + 1

	productDesc := []string{
		fmt.Sprintf("*✔️ %s* 📦\n\n", sample.Products[count].Name),
		fmt.Sprintf("*⚙️ Производитель:*  %s\n\n", sample.Products[count].Brand),
		fmt.Sprintf("*🏷️ Цена без скидки:*  %d %s\n\n", sample.Products[count].Price.BasePrice, sample.Currency),
		fmt.Sprintf("*🏷️ Цена со скидкой:*  %d %s\n\n", sample.Products[count].Price.DiscountPrice, sample.Currency),
		fmt.Sprintf("*🔖 Скидка:*  %d%%\n\n", sample.Products[count].Price.Discount),
		fmt.Sprintf("*🔗 Поставщик:* %s\n\n", sample.Products[count].Supplier),
		fmt.Sprintf("*🛒 Маркет:* %s\n\n", sample.Market),
		fmt.Sprintf("*📦 Товар:*\n%s\n\n", sample.Products[count].Links.URL),
		fmt.Sprintf("*Выборка товаров:*\n%s\n\n", sample.SampleLink),
	}

	buffer := bytes.Buffer{}

	for _, desc := range productDesc {
		buffer.WriteString(desc)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Megamarket 🛍️", megamarket)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Wildberries 🌸", wildberries)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить товар в избранное ⭐", addToFavorite)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	p.botConf.bot.Send(message)
}

// productsIter defines the logic of iterating the user's products sample.
func (p *productsMode) productsIter(chatID int64, market string) {

	if p.botConf.users[chatID].lastAction != productsIter {
		choiceText := "*Выбери, откуда ты хочешь получить товар* 👇"
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Megamarket 🛍️", megamarket)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Wildberries 🌸", wildberries)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
		)
		message := tgbotapi.NewMessage(chatID, choiceText)

		message.ParseMode = markDown
		message.ReplyMarkup = keyboard

		p.botConf.bot.Send(message)

		p.botConf.users[chatID].lastAction = productsIter

		return
	}

	p.nextProduct(chatID, market)
}

// showRequest shows the finished request that will use to get the products.
func (p *productsMode) showRequest(chatID int64) {
	p.botConf.users[chatID].lastAction = showRequest

	request := "✔*Запрос готов! 📝*\n\n*✔Маркеты поиска 🛒*\n"

	for _, market := range p.botConf.users[chatID].request.Markets {
		request += fmt.Sprintf("• %s\n", market)
	}

	request += fmt.Sprintf("\n*Товар: %s* 📦\n\n", p.botConf.users[chatID].request.Query)

	request += "*Если ты заметил, что ошибся в запросе - собери заново!* 👇"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Запустить поиск 🔎", startSearch)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Собрать заново 🔁", bestPriceModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, request)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	p.botConf.bot.Send(message)
}

// addMarket adds the market to the request for the current ChatID.
func (p *productsMode) addMarket(update *tgbotapi.Update) {
	var chatID = update.CallbackQuery.From.ID

	request := p.botConf.users[chatID].request

	p.botConf.users[chatID].request.Markets[update.CallbackQuery.Data] = update.CallbackQuery.Data

	p.botConf.users[chatID].request = dto.ProductRequest{
		Mode:    request.Mode,
		Markets: request.Markets,
	}
}

// setQuery sets the product query request for the current ChatID.
func (p *productsMode) setQuery(update *tgbotapi.Update) {
	var chatID = update.Message.Chat.ID

	request := p.botConf.users[chatID].request

	request.Query = update.Message.Text

	p.botConf.users[chatID].request = request
}
