package tgbot

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/kafka"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/kafka/hand"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// trackedMode defines the logic of the tracked mode user's interaction.
type trackedMode struct {
	botConf      *tgBotConfigs
	logger       *slog.Logger
	repo         services.Repository
	api          services.ApiInteractor
	trackedProds chan *kafka.TrackedProduct
	reader       services.Reader
}

func newTrackedMode(
	botConf *tgBotConfigs,
	repo services.Repository,
	logger *slog.Logger,
	api services.ApiInteractor,
	prods chan *kafka.TrackedProduct,
	reader services.Reader,
) trackedMode {
	return trackedMode{
		botConf:      botConf,
		logger:       logger,
		repo:         repo,
		api:          api,
		trackedProds: prods,
		reader:       reader,
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

// modeErrHandler defines the logic of handling the errors.
func (t trackedMode) modeErrHandler(chatID int64, menu string) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отслеживаемые товары 🔔", trackedModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, menu)

	message.ParseMode = markDown
	message.ReplyMarkup = keyboard

	t.botConf.bot.Send(message)
}

// mode defines the logic of start the setting the tracked products.
func (t trackedMode) mode(chatID int64) {
	const op = "tgbot.add-tracked-product"

	if _, flagExist := t.botConf.users[chatID]; !flagExist {
		t.botConf.users[chatID] = newUserConfig()
	}
	t.botConf.users[chatID].lastAction = addTrackedProductData

	if flagExist, err := t.repo.IsTrackedProductExists(context.Background(), chatID); err != nil {
		t.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		t.modeErrHandler(chatID, fmt.Sprint("*Упс... Похоже, произошла ошибка 😞*\n\n",
			"*Попробуй зайти позже...*! ⏳\n\n"))
		return

	} else if flagExist {
		t.modeErrHandler(chatID, fmt.Sprint("*Упс... Кажется, ты уже поставил уведомление на товар!\n\n*",
			"*Чтобы переустановить уведомление, сними его*! 🔔\n\n"))
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
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отслеживаемые товары 🔔", trackedModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMode

	t.botConf.bot.Send(message)
}

// productSetter defines the logic of setting the tracked product's query.
func (t trackedMode) productSetter(chatID int64) {
	if len(t.botConf.users[chatID].request.Markets) == 0 {
		message := tgbotapi.NewMessage(chatID, fmt.Sprint("*Упс... Кажется, ты не задал ни один маркет поиска 🛒*\n\n",
			"*Задай сначала их, а затем товар 📦*",
		))
		message.ParseMode = markDown
		message.ReplyMarkup = t.botConf.getKeyBoardWithMarkets(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Задать товар 📦", productSetter)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
		)

		t.botConf.bot.Send(message)
		return
	}

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
		t.modeErrHandler(chatID, fmt.Sprint("*Упс... Похоже, произошла ошибка 😞*\n\n",
			"*Попробуй зайти позже...*! ⏳\n\n"))
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

// getTrackedProduct defines the logic of getting the tracked product for the current chatID.
func (t trackedMode) getTrackedProduct(chatID int64) {
	const op = "tgbot.get-tracked-product"

	if _, flagExist := t.botConf.users[chatID]; !flagExist {
		t.botConf.users[chatID] = newUserConfig()
	}

	t.botConf.users[chatID].lastAction = getTrackedProdMode
	product, flagExist, err := t.repo.GetTrackedProduct(context.Background(), chatID)

	if err != nil {
		t.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))

		t.modeErrHandler(chatID, fmt.Sprint("*Упс... Похоже, произошла ошибка 😞*\n\n",
			"*Попробуй зайти позже...*! ⏳\n\n"))
		return

	} else if !flagExist {
		t.modeErrHandler(chatID,
			fmt.Sprint("*Упс... Похоже, у тебя еще нет отслеживаемого товара 🔔*\n\n",
				"*Давай установим его*! 👇\n\n"),
		)
		return
	}

	request := "*Твой текущий отслеживаемый товар 🔔*\n"

	for _, market := range product.Markets {
		request += fmt.Sprintf("• %s\n", market)
	}
	request += fmt.Sprintf("\n*Товар: %s* 📦\n", product.Query)
	request += "\n*Диапазон цен:* минимально возможные цены 🎚️\n\n"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отслеживаемые товары 🔔", trackedModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, request)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	t.botConf.bot.Send(message)
}

// deleteTrackedProduct defines the logic of deleting the tracked mode.
func (t trackedMode) deleteTrackedProduct(chatID int64) {
	const op = "tgbot.delete-tracked-product"

	if _, flagExist := t.botConf.users[chatID]; !flagExist {
		t.botConf.users[chatID] = newUserConfig()
	}

	t.botConf.users[chatID].lastAction = deleteTrackedProductData

	err := t.repo.DeleteTrackedProduct(context.Background(), chatID)

	if err != nil {
		t.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		t.modeErrHandler(chatID, fmt.Sprint("*Упс... Похоже, произошла ошибка 😞*\n\n",
			"*Попробуй зайти позже...*! ⏳\n\n"))
		return
	}

	menu := "*Уведомление с товара было снято успешно! 🔔*\n\n"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Отслеживаемые товары 🔔", trackedModeData)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, menu)

	message.ParseMode = markDown
	message.ReplyMarkup = keyboard

	t.botConf.bot.Send(message)
}

// showTrackedProduct defines the logic of showing the tracked products.
func (t trackedMode) showTrackedProduct(chatID int64) {
	iterInstrs := []string{
		"*Я вернулся с хорошими новостями!* 😊\n\n",
		"*Твой отслеживамый товар получен!*\n\n",
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

	t.botConf.bot.Send(message)
}

// readTrackedProducts reads the tracked products from the chan connected with the kafka's consumer.
func (t trackedMode) readTrackedProducts() {
	go t.reader.ReadProducts(context.Background(), hand.NewProductsHandler(t.logger, t.trackedProds))

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

// close releases the resources of the tracked products.
func (t trackedMode) close() {
	t.reader.Close()
	close(t.trackedProds)
}
