package tgbot

import (
	"bytes"
	"fmt"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// start is the action on the /start command.
func (t *TgBot) start(update *tgbotapi.Update) {
	t.userLastAction[update.Message.Chat.ID] = "start"

	var greets = []string{
		"*Привет, меня зовут Скрудж, и я очень люблю экономить время людей!* 🦆\n\n",
		"📝*Немного обо мне*\nЯ бот, который поможет тебе найти оптимальную цену на нужный товар 🛒\n\n",
		"*Чтобы воспользоваться моими функциями переходи в меню 👇*",
	}

	buffer := bytes.Buffer{}

	for _, greet := range greets {
		buffer.WriteString(greet)
	}

	t.userInteractor.IdentifyUser(update.Message.Chat.ID)

	var keyboardStart = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction),
	))
	var message = tgbotapi.NewMessage(update.Message.Chat.ID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardStart

	t.bot.Send(message)
}

// menu is the action on the /menu command or pressing the menu-button.
func (t *TgBot) menu(chatID int64) {
	t.userLastAction[chatID] = menuAction
	t.userRequest[chatID] = dto.ProductRequest{}

	var menu = []string{
		"*Вот, с чем я могу тебе помочь:*\n\n",
		"✔*Best price*\n",
		"- найти товары по минимальной цене (самые дешевые) 📉\n\n",
		"✔*Избранное:*\n",
		"- вести избранные товары ⭐\n\n",
		"✔*Отслеживаемые товары:*\n",
		"- вести отслеживание товаров 🔔\n\n",
		"*Выбери режим, напиши запрос для товара, а я подберу оптимальную цену для него: быстро и выгодно* 💲",
	}

	var keyboardMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Best price 📉", bestPriceModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Избранное ⭐", favoriteModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отслеживаемые товары 🔔", trackedModeData)),
	)

	buffer := bytes.Buffer{}

	for _, opt := range menu {
		buffer.WriteString(opt)
	}

	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMenu

	t.bot.Send(message)
}

// showRequest shows the finished request that will use to get the products.
func (t *TgBot) showRequest(update *tgbotapi.Update) {
	t.userLastAction[update.Message.Chat.ID] = showRequest

	var request = "✔*Запрос готов! 📝*\n\n*✔Маркеты поиска 🛒*\n"

	for _, market := range t.userRequest[update.Message.Chat.ID].Markets {
		request += fmt.Sprintf("• %s\n", market)
	}

	request += fmt.Sprintf("\n*Товар: %s* 📦\n\n", t.userRequest[update.Message.Chat.ID].Query)

	request += "*Если ты заметил, что ошибся в запросе - собери заново!* 👇"

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Запустить поиск 🔎", startSearch)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Собрать заново 🔁", bestPriceModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	var message = tgbotapi.NewMessage(update.Message.Chat.ID, request)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	t.bot.Send(message)
}
