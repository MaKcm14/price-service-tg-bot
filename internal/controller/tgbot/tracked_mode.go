package tgbot

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// trackedMode defines the logic of the tracked mode user's interaction.
type trackedMode struct {
	botConf *tgBotConfigs
	logger  *slog.Logger
	repo    services.Repository
}

func newTrackedMode(botConf *tgBotConfigs, repo services.Repository, logger *slog.Logger) trackedMode {
	return trackedMode{
		botConf: botConf,
		logger:  logger,
		repo:    repo,
	}
}

// trackedModeMenu defines the logic of user's tracked-prods menu.
func (t trackedMode) trackedModeMenu(chatID int64) {
	if _, flagExist := t.botConf.users[chatID]; !flagExist {
		t.botConf.users[chatID] = newUserConfig()
	}

	t.botConf.users[chatID].lastAction = trackedModeData
	t.botConf.users[chatID].request = dto.NewProductRequest()

	menu := []string{
		"*Ты перешёл в режим Отслеживаемые товары* 🔔\n\n",
		"- Здесь можно найти или установить уведомление на товар\n\n",
		"❓*Как он работает?*\n\n",
		"- Порой, необходимо отслеживать товары постоянно, но на это не хватает времени... 🙁\n\n",
		"🦆 Скрудж поможет тебе установить уведомление на товар, информация о котором будет приходить тебе автоматически раз в сутки!\n\n",
		"❓*Как его использовать?*\n\n",
		"- Необходимо нажать на кнопку *Товар 🔔*, если хочешь посмотреть, на какой товар поставлено уведомление\n\n",
		"- Если хочешь снять уведоление с товара, то нажми на кнопку *Удалить товар* 🗑️\n\n",
		"- Если хочешь поставить уведомление на товар, то нажми кнопку *Добавить товар 🛒*\n\n",
		"- Если хочешь вернуться в меню, нажми *Меню* 📋\n\n",
		"*К товару!* 👇",
	}

	keyboardMenu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Товар 🔔", getTrackedProdMode)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Добавить товар 🛒", addTrackedProductData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Удалить товар 🗑️", deleteTrackedProductData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	buffer := bytes.Buffer{}

	for _, opt := range menu {
		buffer.WriteString(opt)
	}

	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMenu

	t.botConf.bot.Send(message)
}

// mode defines the logic of start the setting the tracked products.
func (t trackedMode) mode(chatID int64) {
	const op = "tgbot.add-tracked-product"
	t.botConf.users[chatID].lastAction = addTrackedProductData

	if flagExist, err := t.repo.IsTrackedProductExists(context.Background(), chatID); err != nil {
		t.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		menu := fmt.Sprint("**Упс... Похоже, произошла ошибка 😞*\n\n*",
			"*Попробуй зайти позже...*! ⏳\n\n")
		message := tgbotapi.NewMessage(chatID, menu)
		message.ParseMode = markDown

		t.botConf.bot.Send(message)
		t.trackedModeMenu(chatID)
		return

	} else if flagExist {
		menu := fmt.Sprint("*Упс... Кажется, ты уже поставил уведомление на товар!\n\n*",
			"*Чтобы переустановить уведомление, сними его*! 🔔\n\n")

		message := tgbotapi.NewMessage(chatID, menu)
		message.ParseMode = markDown

		t.botConf.bot.Send(message)
		t.trackedModeMenu(chatID)
		return
	}

	menu := []string{
		"*Ты перешёл в режим добавления уведомления на товар 🛒*\n\n",
		"❓*Как его использовать?*\n",
		"- Необходимо нажать на кнопки тех маркетов, в которых ты хочешь искать\n\n",
		"- Затем необходимо ввести название товара, который ты хочешь найти\n\n",
		"*P.S. название товара должно быть максимально точным для увеличения точности поиска*\n\n",
		"*Давай поставим уведомление!* 👇",
	}

	buffer := bytes.Buffer{}

	for _, instruct := range menu {
		buffer.WriteString(instruct)
	}

	keyboardMode := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Задать маркеты 🛒", marketSetterMode)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMode

	t.botConf.bot.Send(message)
}

// productSetter defines the logic of setting the tracked product's query.
func (t trackedMode) productSetter(chatID int64) {
	t.botConf.users[chatID].lastAction = productSetter

	message := tgbotapi.NewMessage(chatID,
		"*Введи точное название товара, по которому будет осуществляться поиск* 📦")

	message.ParseMode = markDown
	t.botConf.bot.Send(message)
}

// setRequest sets the product query request for the current ChatID.
func (t trackedMode) setRequest(update *tgbotapi.Update) {
	var chatID = update.Message.Chat.ID

	request := t.botConf.users[chatID].request

	request.Query = update.Message.Text

	t.botConf.users[chatID].request = request
}

// showRequest shows the request for the tracked product to the user.
func (t trackedMode) showRequest(chatID int64) {
	const op = "tgbot.show-request"

	t.botConf.users[chatID].lastAction = showRequest

	err := t.repo.AddTrackedProduct(context.Background(), chatID, t.botConf.users[chatID].request)

	if err != nil {
		t.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		menu := fmt.Sprint("**Упс... Похоже, произошла ошибка 😞*\n\n*",
			"*Попробуй зайти позже...*! ⏳\n\n")

		message := tgbotapi.NewMessage(chatID, menu)

		message.ParseMode = markDown

		t.botConf.bot.Send(message)

		t.trackedModeMenu(chatID)
		return
	}

	request := "✔*Запрос готов! 📝*\n\n*✔Маркеты поиска 🛒*\n"

	for _, market := range t.botConf.users[chatID].request.Markets {
		request += fmt.Sprintf("• %s\n", market)
	}

	request += fmt.Sprintf("\n*Товар: %s* 📦\n", t.botConf.users[chatID].request.Query)

	request += "\n*Диапазон цен:* минимально возможные цены 🎚️\n\n"

	request += "*Если ты заметил, что ошибся в запросе - сними уведомление и собери заново!* 👇"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отслеживаемые товары 🔔", trackedModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, request)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	t.botConf.bot.Send(message)
}

func (t trackedMode) startSearch(chatID int64) {
	// logic of sending the request to the price-service
}
