package tgbot

import (
	"bytes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t TgBot) start(update *tgbotapi.Update) {
	// write here the start logic with greetings and identification.
}

// menu is the action on the /menu command.
func (t TgBot) menu(update *tgbotapi.Update) {
	var menu = []string{
		"*Меня зовут Скрудж, и я очень люблю экономить время людей!* 🦆\n\n",
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
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Price range 📊", "price-range")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Exact price 🏷️", "exact-price")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Best price 📉", "best-price")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Обычный поиск 🔎", "search")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Избранное ⭐", "favourite-products")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отслеживаемые товары 🔔", "tracked-products")),
	)

	buffer := bytes.Buffer{}

	for _, opt := range menu {
		buffer.WriteString(opt)
	}

	message := tgbotapi.NewMessage(update.Message.Chat.ID, buffer.String())

	message.ParseMode = ParseMode
	message.ReplyMarkup = keyboardMenu
	t.bot.Send(message)
}
