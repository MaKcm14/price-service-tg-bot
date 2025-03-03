package tgbot

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/api"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type bestPriceMode struct {
	botConf *tgBotConfigs
	logger  *slog.Logger
	api     services.ApiInteractor
}

func newBestPriceMode(log *slog.Logger, bot *tgBotConfigs, api services.ApiInteractor) bestPriceMode {
	return bestPriceMode{
		botConf: bot,
		logger:  log,
		api:     api,
	}
}

// bestPriceMode is the action on the pressing the best-price button.
func (b *bestPriceMode) bestPriceMode(chatID int64) {
	if _, flagExist := b.botConf.users[chatID]; !flagExist {
		b.botConf.users[chatID] = *newUserConfig()
	}

	var newConfig userConfig = b.botConf.users[chatID]

	newConfig.lastAction = bestPriceModeData
	newConfig.request = dto.NewProductRequest(entities.BestPriceMode)

	b.botConf.users[chatID] = newConfig

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

	var message = tgbotapi.NewMessage(chatID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMode

	b.botConf.bot.Send(message)
}

// marketSetterMode sets the markets.
func (b *bestPriceMode) marketSetterMode(chatID int64) {
	var newConfig userConfig = b.botConf.users[chatID]

	newConfig.lastAction = marketSetterMode

	b.botConf.users[chatID] = newConfig

	var keyboardMarketSetter = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Megamarket 🛍️", megamarket)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Wildberries 🌸", wildberries)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Задать товар 📦", productSetter)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	var message = tgbotapi.NewMessage(chatID, "*Выбери маркеты, в которых будет производиться поиск* 👇")

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMarketSetter

	b.botConf.bot.Send(message)
}

// productSetter defines the logic of setting the product's name.
func (b *bestPriceMode) productSetter(chatID int64) {
	var newConfig userConfig = b.botConf.users[chatID]

	newConfig.lastAction = productSetter

	b.botConf.users[chatID] = newConfig

	var message = tgbotapi.NewMessage(chatID,
		"*Введи точное название товара, по которому будет осуществляться поиск* 📦")

	message.ParseMode = markDown
	b.botConf.bot.Send(message)
}

// startSearch defines the logic of searching the products using the finished request.
func (b *bestPriceMode) startSearch(chatID int64) {
	var newConfig userConfig = b.botConf.users[chatID]

	newConfig.lastAction = startSearch

	b.botConf.users[chatID] = newConfig

	products, err := b.api.GetProductsByBestPrice(b.botConf.users[chatID].request)

	if err != nil {
		var errText = "*Упс... Похоже, произошла ошибка 😞*"

		if errors.Is(err, api.ErrApiInteraction) {
			errText += "\n\n*Что-то не так с парсером... \nПопробуй отключить VPN или попробовать позже ⏳*"
		}

		var message = tgbotapi.NewMessage(chatID, errText)
		var keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
		)

		message.ReplyMarkup = keyboard
		message.ParseMode = markDown

		b.botConf.bot.Send(message)

		return
	}

	newConfig = b.botConf.users[chatID]

	newConfig.sample.sample = products

	b.botConf.users[chatID] = newConfig

	markets := make(map[string]int)

	for _, market := range b.botConf.users[chatID].request.Markets {
		markets[market] = 0
	}

	newConfig = b.botConf.users[chatID]

	newConfig.sample.samplePtr = markets

	b.botConf.users[chatID] = newConfig

	var iterInstrs = []string{
		"*Запрос был обработан успешно!* 😊\n\n",
		"❓*Как использовать поиск?*\n",
		"✔ Нажимай на тот маркет, товар которого хочешь посмотреть\n",
		"✔ Если хочешь добавить товар в Избранное, нажми на кнопку\n",
		"*Давай искать!* 👇",
	}

	buffer := bytes.Buffer{}

	for _, instruct := range iterInstrs {
		buffer.WriteString(instruct)
	}

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Смотреть товары 📦", productsIter)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	var message = tgbotapi.NewMessage(chatID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboard

	b.botConf.bot.Send(message)
}

// productsIter defines the logic of iterating the user's products sample.
func (b *bestPriceMode) productsIter(chatID int64, market string) {
	var newConfig userConfig = b.botConf.users[chatID]

	newConfig.sample.lastMarketChoice = market

	b.botConf.users[chatID] = newConfig

	if b.botConf.users[chatID].lastAction != productsIter {
		var choiceText = "*Выбери, откуда ты хочешь получить товар* 👇"
		var keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Megamarket 🛍️", megamarket)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Wildberries 🌸", wildberries)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
		)
		var message = tgbotapi.NewMessage(chatID, choiceText)

		message.ParseMode = markDown
		message.ReplyMarkup = keyboard

		b.botConf.bot.Send(message)

		newConfig = b.botConf.users[chatID]

		newConfig.lastAction = productsIter

		b.botConf.users[chatID] = newConfig

		return
	}

	count := b.botConf.users[chatID].sample.samplePtr[market]

	if sample := b.botConf.users[chatID].sample.sample[market]; len(sample.Products) <= count {
		return
	}

	sample := b.botConf.users[chatID].sample.sample[market]

	newConfig = b.botConf.users[chatID]

	newConfig.sample.samplePtr[market] = count + 1

	b.botConf.users[chatID] = newConfig

	var productDesc = []string{
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

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Megamarket 🛍️", megamarket)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Wildberries 🌸", wildberries)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить товар в избранное ⭐", addToFavorite)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	var message = tgbotapi.NewMessage(chatID, buffer.String())

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	b.botConf.bot.Send(message)
}

// showRequest shows the finished request that will use to get the products.
func (b *bestPriceMode) showRequest(chatID int64) {
	var newConfig userConfig = b.botConf.users[chatID]

	newConfig.lastAction = showRequest

	b.botConf.users[chatID] = newConfig

	var request = "✔*Запрос готов! 📝*\n\n*✔Маркеты поиска 🛒*\n"

	for _, market := range b.botConf.users[chatID].request.Markets {
		request += fmt.Sprintf("• %s\n", market)
	}

	request += fmt.Sprintf("\n*Товар: %s* 📦\n\n", b.botConf.users[chatID].request.Query)

	request += "*Если ты заметил, что ошибся в запросе - собери заново!* 👇"

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Запустить поиск 🔎", startSearch)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Собрать заново 🔁", bestPriceModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	var message = tgbotapi.NewMessage(chatID, request)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	b.botConf.bot.Send(message)
}

// addMarket adds the market to the request for the current ChatID.
func (b *bestPriceMode) addMarket(update *tgbotapi.Update) {
	var chatID = update.CallbackQuery.From.ID

	request := b.botConf.users[chatID].request

	request.Markets[update.CallbackQuery.Data] = update.CallbackQuery.Data

	newConfig := b.botConf.users[chatID]

	newConfig.request = dto.ProductRequest{
		Mode:    request.Mode,
		Markets: request.Markets,
	}

	b.botConf.users[chatID] = newConfig
}

// setQuery sets the product query request for the current ChatID.
func (b *bestPriceMode) setQuery(update *tgbotapi.Update) {
	var chatID = update.Message.Chat.ID

	request := b.botConf.users[chatID].request

	request.Query = update.Message.Text

	newConfig := b.botConf.users[chatID]

	newConfig.request = request

	b.botConf.users[chatID] = newConfig
}
