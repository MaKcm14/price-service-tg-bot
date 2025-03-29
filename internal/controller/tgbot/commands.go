package tgbot

import (
	"bytes"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
)

// start is the action on the /start command.
func (t *TgBot) start(chatID int64) {
	if _, flagExist := t.botConf.users[chatID]; !flagExist {
		t.botConf.users[chatID] = newUserConfig()
	}

	t.botConf.users[chatID].lastAction = startAction

	greets := []string{
		"*Привет, меня зовут Скрудж, и я очень люблю экономить время людей!* 🦆\n\n",
		"📝*Немного обо мне*\nЯ бот, который поможет тебе найти оптимальную цену на нужный товар 🛒\n\n",
		"*Чтобы воспользоваться моими функциями переходи в меню 👇*",
	}

	buffer := bytes.Buffer{}

	for _, greet := range greets {
		buffer.WriteString(greet)
	}

	t.uinter.IdentifyUser(chatID)

	keyboardStart := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction),
	))
	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardStart

	t.botConf.bot.Send(message)
}

// menu is the action on the /menu command or pressing the menu-button.
func (t *TgBot) menu(chatID int64) {
	if _, flagExist := t.botConf.users[chatID]; !flagExist {
		t.botConf.users[chatID] = newUserConfig()
	}

	t.botConf.users[chatID].lastAction = menuAction
	t.botConf.users[chatID].request = dto.ProductRequest{}

	menu := []string{
		"*Вот, с чем я могу тебе помочь:*\n\n",
		"✔*Best price*\n",
		"- найти товары по минимальной цене (самые дешевые) 📉\n\n",
		"✔*Избранное:*\n",
		"- вести избранные товары ⭐\n\n",
		"✔*Отслеживаемые товары:*\n",
		"- вести отслеживание товаров 🔔\n\n",
		"*Выбери режим, напиши запрос для товара, а я подберу оптимальную цену для него: быстро и выгодно* 💲",
	}

	keyboardMenu := tgbotapi.NewInlineKeyboardMarkup(
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

	t.botConf.bot.Send(message)

	t.botConf.users[chatID].favorites = newUserFavoritesConfig()
}
