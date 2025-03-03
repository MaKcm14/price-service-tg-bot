package tgbot

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/api"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
)

// bestPriceMode defines the logic of best price mode processing.
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
		b.botConf.users[chatID] = newUserConfig()
	}

	b.botConf.users[chatID].lastAction = bestPriceModeData
	b.botConf.users[chatID].request = dto.NewProductRequest(entities.BestPriceMode)

	priceRangeInstructs := []string{
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

	keyboardMode := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Задать маркеты 🛒", marketSetterMode)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMode

	b.botConf.bot.Send(message)
}

// marketSetterMode sets the markets.
func (b *bestPriceMode) marketSetterMode(chatID int64) {
	b.botConf.users[chatID].lastAction = marketSetterMode

	keyboardMarketSetter := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Megamarket 🛍️", megamarket)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Wildberries 🌸", wildberries)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Задать товар 📦", productSetter)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, "*Выбери маркеты, в которых будет производиться поиск* 👇")

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMarketSetter

	b.botConf.bot.Send(message)
}

// productSetter defines the logic of setting the product's name.
func (b *bestPriceMode) productSetter(chatID int64) {
	b.botConf.users[chatID].lastAction = productSetter

	message := tgbotapi.NewMessage(chatID,
		"*Введи точное название товара, по которому будет осуществляться поиск* 📦")

	message.ParseMode = markDown
	b.botConf.bot.Send(message)
}

// errorOfSearch defines the logic of searching's error processing.
func (b *bestPriceMode) errorOfSearchMode(chatID int64, err error) {
	var errText = "*Упс... Похоже, произошла ошибка 😞*"

	if errors.Is(err, api.ErrApiInteraction) {
		errText += "\n\n*Что-то не так с парсером... \nПопробуй отключить VPN или попробовать позже ⏳*"
	}

	message := tgbotapi.NewMessage(chatID, errText)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	b.botConf.bot.Send(message)
}

// searchReply defines the logic of searching's reply.
func (b *bestPriceMode) searchModeReply(chatID int64) {
	iterInstrs := []string{
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

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Смотреть товары 📦", productsIter)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboard

	b.botConf.bot.Send(message)
}

// startSearch defines the logic of searching the products using the finished request.
func (b *bestPriceMode) startSearch(chatID int64) {
	const op = "tgbot.best-price-search"

	b.botConf.users[chatID].lastAction = startSearch

	products, err := b.api.GetProductsByBestPrice(b.botConf.users[chatID].request)

	if err != nil {
		b.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		b.errorOfSearchMode(chatID, err)
		return
	}

	b.botConf.users[chatID].sample.sample = products

	markets := make(map[string]int)

	for _, market := range b.botConf.users[chatID].request.Markets {
		markets[market] = 0
	}

	b.botConf.users[chatID].sample.samplePtr = markets

	b.searchModeReply(chatID)
}
