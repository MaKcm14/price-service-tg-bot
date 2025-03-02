package tgbot

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
	"github.com/MaKcm14/price-service/pkg/entities"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type favoriteMode struct {
	botConf tgBotConfigs
	logger  *slog.Logger

	repo services.Repository
}

func newFavoriteMode(log *slog.Logger, bot tgBotConfigs, repo services.Repository) favoriteMode {
	return favoriteMode{
		botConf: bot,
		logger:  log,
		repo:    repo,
	}
}

// addFavoriteProduct adds the product to the favorites.
func (f *favoriteMode) addFavoriteProduct(chatID int64) {
	const op = "tgbot.add-favorite-product"

	var response = "*Товар был успешно добавлен! ⭐*"

	market := f.botConf.usersConfig.sampleConfig.usersLastMarketChoice[chatID]
	count := f.botConf.usersConfig.sampleConfig.usersSamplePtr[chatID][market] - 1
	sample := f.botConf.usersConfig.sampleConfig.usersSample[chatID][market]

	product := sample.Products[count]

	err := f.repo.AddFavoriteProducts(context.Background(), chatID, []entities.Product{product})

	if err != nil {
		f.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		response = "*Упс... Похоже, произошла ошибка 😞*"
	}

	var message = tgbotapi.NewMessage(chatID, response)

	message.ParseMode = markDown

	f.botConf.bot.Send(message)
}

// favoriteMode defines the entry-point to the favorite mode.
func (f *favoriteMode) favoriteMode(chatID int64) {
	f.botConf.usersConfig.usersLastAction[chatID] = favoriteModeData

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

	var message = tgbotapi.NewMessage(chatID, buffer.String())

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	f.botConf.bot.Send(message)
}

// showFavoriteProducts shows the user's favorite products.
func (f *favoriteMode) showFavoriteProducts(chatID int64) {
	const op = "tgbot.show-favorite-products"

	var product entities.Product

	products, err := f.repo.GetFavoriteProducts(context.Background(), chatID)

	if err != nil {
		f.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))

		response := "*Упс... Похоже, произошла ошибка 😞*"
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
		)

		var message = tgbotapi.NewMessage(chatID, response)

		message.ReplyMarkup = keyboard
		message.ParseMode = markDown

		f.botConf.bot.Send(message)
		return
	} else if len(f.botConf.usersConfig.favoriteConfig.usersFavoriteLastProdsID[chatID]) == len(products) {
		response := "*Товаров больше нет 📦*"
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
		)

		var message = tgbotapi.NewMessage(chatID, response)

		message.ReplyMarkup = keyboard
		message.ParseMode = markDown

		f.botConf.bot.Send(message)
		return
	}

	for key, val := range products {
		if _, flagExist := f.botConf.usersConfig.favoriteConfig.usersFavoriteLastProdsID[chatID][key]; !flagExist {
			product = val
			f.botConf.usersConfig.favoriteConfig.lastFavoriteProdID[chatID] = key
			f.botConf.usersConfig.favoriteConfig.usersFavoriteLastProdsID[chatID][key] = struct{}{}
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

	var message = tgbotapi.NewMessage(chatID, buffer.String())

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	f.botConf.bot.Send(message)
}

// deleteFavoriteProduct defines the logic of the deleting the user's favorite product.
func (f *favoriteMode) deleteFavoriteProduct(chatID int64) {
	const op = "tgbot.delete-favorite-product"

	var response = "*Товар был успешно удалён 🗑️*"

	err := f.repo.DeleteFavoriteProducts(context.Background(), chatID, []int{f.botConf.usersConfig.favoriteConfig.lastFavoriteProdID[chatID]})

	if err != nil {
		f.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		response = "*Упс... Похоже, произошла ошибка 😞*"
	} else {
		delete(f.botConf.usersConfig.favoriteConfig.usersFavoriteLastProdsID[chatID],
			f.botConf.usersConfig.favoriteConfig.lastFavoriteProdID[chatID])
	}

	var message = tgbotapi.NewMessage(chatID, response)

	var keyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Следующий товар ➡️", showFavoriteProducts)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	f.botConf.bot.Send(message)
}
