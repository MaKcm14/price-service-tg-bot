package tgbot

import (
	"bytes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// start is the action on the /start command.
func (t TgBot) start(update *tgbotapi.Update) {
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
		tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menu),
	))
	var message = tgbotapi.NewMessage(update.Message.Chat.ID, buffer.String())

	message.ParseMode = parseMode
	message.ReplyMarkup = keyboardStart

	t.bot.Send(message)
}

// menu is the action on the /menu command.
func (t TgBot) menu(chatID int64) {
	var menu = []string{
		"*Вот, с чем я могу тебе помочь:*\n\n",
		"✔*Price range*\n",
		"- найти товары по заданному тобою диапазону цен 📊\n\n",
		"✔*Exact price*\n",
		"- найти товары по твоей точной цене 🏷️\n\n",
		"✔*Best price*\n",
		"- найти товары по минимальной цене (самые дешевые) 📉\n\n",
		"✔*Обычный поиск*\n",
		"- осуществить обычный поиск товаров 🔎\n\n",
		"✔*Избранное:*\n",
		"- вести избранные товары ⭐\n\n",
		"✔*Отслеживаемые товары:*\n",
		"- вести отслеживание товаров 🔔\n\n",
		"*Выбери режим, напиши запрос для товара, а я подберу оптимальную цену для него: быстро и выгодно* 💲",
	}

	var keyboardMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Price range 📊", priceRangeModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Exact price 🏷️", exactPriceModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Best price 📉", bestPriceModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Обычный поиск 🔎", searchModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Избранное ⭐", favouriteModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отслеживаемые товары 🔔", trackedModeData)),
	)

	buffer := bytes.Buffer{}

	for _, opt := range menu {
		buffer.WriteString(opt)
	}

	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ParseMode = parseMode
	message.ReplyMarkup = keyboardMenu

	t.bot.Send(message)
}
