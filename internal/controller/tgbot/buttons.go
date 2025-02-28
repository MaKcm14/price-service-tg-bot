package tgbot

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/api"

	epkg "github.com/MaKcm14/price-service/pkg/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// bestPriceMode is the action on the pressing the best-price button.
func (t *TgBot) bestPriceMode(update *tgbotapi.Update) {
	t.userLastAction[update.CallbackQuery.From.ID] = bestPriceModeData
	t.userRequest[update.CallbackQuery.From.ID] = dto.NewProductRequest(entities.BestPriceMode)

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

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMode

	t.bot.Send(message)
}

// marketSetterMode sets the markets.
func (t *TgBot) marketSetterMode(update *tgbotapi.Update) {
	t.userLastAction[update.CallbackQuery.From.ID] = marketSetterMode

	var keyboardMarketSetter = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Megamarket 🛍️", megamarket)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Wildberries 🌸", wildberries)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Задать товар 📦", productSetter)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, "*Выбери маркеты, в которых будет производиться поиск* 👇")

	message.ParseMode = markDown
	message.ReplyMarkup = keyboardMarketSetter

	t.bot.Send(message)
}

// productSetter defines the logic of setting the product's name.
func (t *TgBot) productSetter(update *tgbotapi.Update) {
	t.userLastAction[update.CallbackQuery.From.ID] = productSetter

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID,
		"*Введи точное название товара, по которому будет осуществляться поиск* 📦")

	message.ParseMode = markDown
	t.bot.Send(message)
}

// startSearch defines the logic of searching the products using the finished request.
func (t *TgBot) startSearch(update *tgbotapi.Update) {
	t.userLastAction[update.CallbackQuery.From.ID] = startSearch

	products, err := t.api.GetProductsByBestPrice(t.userRequest[update.CallbackQuery.From.ID])

	if err != nil {
		var errText = "*Упс... Похоже, произошла ошибка 😞*"

		if errors.Is(err, api.ErrApiInteraction) {
			errText += "\n\n*Что-то не так с парсером... \nПопробуй отключить VPN или попробовать позже ⏳*"
		}

		var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, errText)
		message.ParseMode = markDown

		t.bot.Send(message)

		t.menu(update.CallbackQuery.From.ID)
		return
	}

	t.userSample[update.CallbackQuery.From.ID] = products

	markets := make(map[string]int)

	for _, market := range t.userRequest[update.CallbackQuery.From.ID].Markets {
		markets[market] = 0
	}

	t.userSamplePtr[update.CallbackQuery.From.ID] = markets

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

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, buffer.String())

	message.ParseMode = markDown
	message.ReplyMarkup = keyboard

	t.bot.Send(message)
}

// productsIter defines the logic of iterating the user's products sample.
func (t *TgBot) productsIter(update *tgbotapi.Update, market string) {
	t.userLastMarketChoice[update.CallbackQuery.From.ID] = market

	if t.userLastAction[update.CallbackQuery.From.ID] != productsIter {
		var choiceText = "*Выбери, откуда ты хочешь получить товар* 👇"
		var keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Megamarket 🛍️", megamarket)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Wildberries 🌸", wildberries)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
		)
		var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, choiceText)

		message.ParseMode = markDown
		message.ReplyMarkup = keyboard

		t.bot.Send(message)

		t.userLastAction[update.CallbackQuery.From.ID] = productsIter
		return
	}

	count := t.userSamplePtr[update.CallbackQuery.From.ID][market]

	if sample := t.userSample[update.CallbackQuery.From.ID][market]; len(sample.Products) <= count {
		return
	}

	sample := t.userSample[update.CallbackQuery.From.ID][market]

	t.userSamplePtr[update.CallbackQuery.From.ID][market] = count + 1

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

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, buffer.String())

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	t.bot.Send(message)
}

// addFavoriteProduct adds the product to the favorites.
func (t *TgBot) addFavoriteProduct(update *tgbotapi.Update) {
	const op = "tgbot.add-favorite-product"

	var response = "*Товар был успешно добавлен! ⭐*"

	market := t.userLastMarketChoice[update.CallbackQuery.From.ID]
	count := t.userSamplePtr[update.CallbackQuery.From.ID][market] - 1
	sample := t.userSample[update.CallbackQuery.From.ID][market]

	product := sample.Products[count]

	err := t.repo.AddFavoriteProducts(context.Background(), update.CallbackQuery.From.ID, []epkg.Product{product})

	if err != nil {
		t.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		response = "*Упс... Похоже, произошла ошибка 😞*"
	}

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, response)

	message.ParseMode = markDown

	t.bot.Send(message)
}

// favoriteMode defines the entry-point to the favorite mode.
func (t *TgBot) favoriteMode(update *tgbotapi.Update) {
	t.userLastAction[update.CallbackQuery.From.ID] = favoriteModeData

	var favoriteModeInstruct = []string{
		"*Ты перешёл в режим Избранное* ⭐\n\n",
		"- Здесь можно найти все товары, которые тебе когда-то понравились ❤️\n\n",
		"❓*Как его использовать?*\n\n",
		"- Необходимо нажать на кнопку *Следующий товар ➡️*, если хочешь перейти к следующему товару\n\n",
		"- Если хочешь удалить товар из выборки, то нажми на кнопку *Удалить товар* 🗑️\n\n",
		"- Если хочешь вернуться в меню, нажми *Меню* 📋\n\n",
		"*К товарам!* 👇",
	}

	buffer := bytes.Buffer{}

	for _, instruct := range favoriteModeInstruct {
		buffer.WriteString(instruct)
	}

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Смотреть товары 📦", showFavoriteProducts)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, buffer.String())

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	t.bot.Send(message)
}

// showFavoriteProducts shows the user's favorite products.
func (t *TgBot) showFavoriteProducts(update *tgbotapi.Update) {
	const op = "tgbot.show-favorite-products"

	var product epkg.Product

	products, err := t.repo.GetFavoriteProducts(context.Background(), update.CallbackQuery.From.ID)

	if err != nil {
		t.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))

		response := "*Упс... Похоже, произошла ошибка 😞*"
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
		)

		var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, response)

		message.ReplyMarkup = keyboard
		message.ParseMode = markDown

		t.bot.Send(message)
		return
	} else if len(t.userFavoriteLastProds) == len(products) {
		response := "*Товаров больше нет 📦*"
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
		)

		var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, response)

		message.ReplyMarkup = keyboard
		message.ParseMode = markDown

		t.bot.Send(message)
		return
	}

	for key, val := range products {
		if _, flagExist := t.userFavoriteLastProds[key]; !flagExist {
			product = val
			t.lastFavoriteProd = key
			t.userFavoriteLastProds[key] = struct{}{}
			break
		}
	}

	var productDesc = []string{
		fmt.Sprintf("*✔️ %s* 📦\n\n", product.Name),
		fmt.Sprintf("*⚙️ Производитель:*  %s\n\n", product.Brand),
		fmt.Sprintf("*🏷️ Цена без скидки:*  %d\n\n", product.Price.BasePrice),
		fmt.Sprintf("*🔗 Поставщик:* %s\n\n", product.Supplier),
		fmt.Sprintf("*📦 Товар:*\n%s\n\n", product.Links.URL),
	}

	buffer := bytes.Buffer{}

	for _, desc := range productDesc {
		buffer.WriteString(desc)
	}

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Следующий товар ➡️", showFavoriteProducts)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Удалить товар 🗑️", deleteFavoriteProduct)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, buffer.String())

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	t.bot.Send(message)
}

// deleteFavoriteProduct defines the logic of the deleting the user's favorite product.
func (t *TgBot) deleteFavoriteProduct(update *tgbotapi.Update) {
	const op = "tgbot.delete-favorite-product"

	var response = "*Товар был успешно удалён 🗑️*"

	err := t.repo.DeleteFavoriteProducts(context.Background(), update.CallbackQuery.From.ID, []int{t.lastFavoriteProd})

	if err != nil {
		t.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		response = "*Упс... Похоже, произошла ошибка 😞*"
	} else {
		delete(t.userFavoriteLastProds, t.lastFavoriteProd)
	}

	var message = tgbotapi.NewMessage(update.CallbackQuery.From.ID, response)

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Следующий товар ➡️", showFavoriteProducts)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	t.bot.Send(message)
}
