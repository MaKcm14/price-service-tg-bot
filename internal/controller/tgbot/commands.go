package tgbot

import (
	"bytes"
	"time"

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

	t.userInteractor.IdentifyUser(chatID)

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

func (t *TgBot) showTrackedProduct(chatID int64) {
	iterInstrs := []string{
		"*🦆 Я вернулся с хорошими новостями!* 😊\n\n",
		"*Твой отслеживамый товар получен!*",
		"❓*Как использовать поиск?*\n",
		"✔ Нажимай на тот маркет, товар которого хочешь посмотреть\n",
		"✔ Если хочешь добавить товар в Избранное, нажми на кнопку\n",
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

	t.botConf.bot.Send(message)
}

// readTrackedProducts reads the tracked products from the chan connected with the kafka's consumer.
func (t *TgBot) readTrackedProducts() {
	for products := range t.trackedProds {
		if _, flagExist := t.botConf.users[products.ChatID]; !flagExist {
			t.botConf.users[products.ChatID] = newUserConfig()
		}

		for t.botConf.users[products.ChatID].lastAction == showRequest {
			continue
		}

		t.botConf.users[products.ChatID].sample.sample = products.Response.Sample

		markets := make(map[string]int)

		for _, market := range t.botConf.users[products.ChatID].request.Markets {
			markets[market] = 0
		}

		t.botConf.users[products.ChatID].sample.samplePtr = markets

		t.showTrackedProduct(products.ChatID)

		time.Sleep(time.Minute * 1)
	}
}
