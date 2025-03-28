package tgbot

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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

// mode is the action on the pressing the best-price button.
func (b bestPriceMode) mode(chatID int64) {
	if _, flagExist := b.botConf.users[chatID]; !flagExist {
		b.botConf.users[chatID] = newUserConfig()
	}

	b.botConf.users[chatID].lastAction = bestPriceModeData
	b.botConf.users[chatID].request = dto.NewProductRequest()

	instructs := []string{
		"*Ты перешёл в режим поиска Best Price 📊 *\n\n",
		"❓*Как его использовать?*\n",
		"- Необходимо нажать на кнопки тех маркетов, в которых ты хочешь искать\n\n",
		"- Затем необходимо ввести название товара, который ты хочешь найти\n\n",
		"*P.S. название товара должно быть максимально точным для увеличения точности поиска*\n\n",
		"*Давай поищем!* 👇",
	}

	buffer := bytes.Buffer{}

	for _, instruct := range instructs {
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

// productSetter defines the logic of setting the product's name.
func (b bestPriceMode) productSetter(chatID int64) {
	if len(b.botConf.users[chatID].request.Markets) == 0 {
		message := tgbotapi.NewMessage(chatID, fmt.Sprint("*Упс... Кажется, ты не задал ни один маркет поиска 🛒*\n\n",
			"*Задай сначала их, а затем товар 📦*",
		))

		message.ParseMode = markDown

		message.ReplyMarkup = b.botConf.getKeyBoardWithMarkets(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Задать товар 📦", productSetter)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
		)

		b.botConf.bot.Send(message)
		return
	}

	b.botConf.users[chatID].lastAction = productSetter

	message := tgbotapi.NewMessage(chatID,
		"*Введи точное название товара, по которому будет осуществляться поиск* 📦")

	message.ParseMode = markDown
	b.botConf.bot.Send(message)
}

// modeErrHandler the logic of searching's error processing.
func (b bestPriceMode) modeErrHandler(chatID int64, response string) {
	message := tgbotapi.NewMessage(chatID, response)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	b.botConf.bot.Send(message)
}

// searchModeReply defines the logic of searching's reply.
func (b bestPriceMode) searchModeReply(chatID int64) {
	iterInstrs := []string{
		"*Запрос был обработан успешно!* 😊\n\n",
		"❓*Как использовать поиск?*\n",
		"✔ Нажимай на тот маркет, товар которого хочешь посмотреть\n",
		"✔ Если хочешь добавить товар в Избранное, нажми на кнопку *Добавить*\n\n",
		"*Давай смотреть!* 👇",
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
func (b bestPriceMode) startSearch(chatID int64) {
	const op = "tgbot.best-price-search"

	b.botConf.users[chatID].lastAction = startSearch
	products, err := b.api.GetProductsByBestPrice(b.botConf.users[chatID].request)

	if err != nil {
		b.logger.Warn(fmt.Sprintf("error of the %s: %s", op, err))
		response := "*Упс... Похоже, произошла ошибка 😞*"

		if errors.Is(err, api.ErrApiInteraction) {
			response += "\n\n*Что-то не так с парсером... \nПопробуй отключить VPN или попробовать позже ⏳*"
		}

		b.modeErrHandler(chatID, response)

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

// showRequest shows the finished request that will use to get the products.
func (p bestPriceMode) showRequest(chatID int64) {
	p.botConf.users[chatID].lastAction = showRequest

	request := "✔*Запрос готов! 📝*\n\n*✔Маркеты поиска 🛒*\n"

	for _, market := range p.botConf.users[chatID].request.Markets {
		request += fmt.Sprintf("• %s\n", market)
	}

	request += fmt.Sprintf("\n*Товар: %s* 📦\n", p.botConf.users[chatID].request.Query)
	request += "\n*Диапазон цен:* минимально возможные цены 🎚️\n\n"
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

// setRequest sets the product query request for the current ChatID.
func (p bestPriceMode) setRequest(update *tgbotapi.Update) {
	var chatID = update.Message.Chat.ID

	request := p.botConf.users[chatID].request

	request.Query = update.Message.Text
	p.botConf.users[chatID].request = request
}
