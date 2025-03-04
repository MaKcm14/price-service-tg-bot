package tgbot

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/services"
	"github.com/MaKcm14/price-service/pkg/entities"
)

// favoriteMode defines the logic of the favorite mode processing.
type favoriteMode struct {
	botConf *tgBotConfigs
	logger  *slog.Logger

	repo services.Repository
}

func newFavoriteMode(log *slog.Logger, bot *tgBotConfigs, repo services.Repository) favoriteMode {
	return favoriteMode{
		botConf: bot,
		logger:  log,
		repo:    repo,
	}
}

// addFavoriteProduct adds the product to the favorites.
func (f *favoriteMode) addFavoriteProduct(chatID int64) {
	const op = "tgbot.add-favorite-product"

	response := "*Товар был успешно добавлен! ⭐*"

	market := f.botConf.users[chatID].sample.lastMarketChoice
	count := f.botConf.users[chatID].sample.samplePtr[market] - 1
	sample := f.botConf.users[chatID].sample.sample[market]

	if count < 0 {
		return
	}

	product := sample.Products[count]

	err := f.repo.AddFavoriteProducts(context.Background(), chatID, []entities.Product{product})

	if err != nil {
		f.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		response = "*Упс... Похоже, произошла ошибка 😞*"
	}

	message := tgbotapi.NewMessage(chatID, response)

	message.ParseMode = markDown

	f.botConf.bot.Send(message)
}

// favoriteMode defines the entry-point to the favorite mode.
func (f *favoriteMode) favoriteMode(chatID int64) {
	if _, flagExist := f.botConf.users[chatID]; !flagExist {
		f.botConf.users[chatID] = newUserConfig()
	}

	f.botConf.users[chatID].lastAction = favoriteModeData

	favoriteModeInstruct := []string{
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

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Смотреть товары 📦", showFavoriteProducts)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	f.botConf.bot.Send(message)
}

// showProduct shows the favorite products.
func (f *favoriteMode) showProduct(chatID int64, products map[int]entities.Product) {
	var product entities.Product

	for key, val := range products {
		if _, flagExist := f.botConf.users[chatID].favorites.favoriteLastProdsID[key]; !flagExist {
			product = val
			f.botConf.users[chatID].favorites.lastFavoriteProdID = key
			f.botConf.users[chatID].favorites.favoriteLastProdsID[key] = struct{}{}
			break
		}
	}

	productDesc := []string{
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

	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Следующий товар ➡️", showFavoriteProducts)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Удалить товар 🗑️", deleteFavoriteProduct)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, buffer.String())

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	f.botConf.bot.Send(message)
}

// showProductModeGettingErrHandler defines the logic of processing
// the error of getting the favorite products.
func (f *favoriteMode) showProductModeGettingErrHandler(chatID int64, err error) {
	const op = "tgbot.show-favorite-products"

	f.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))

	response := "*Упс... Похоже, произошла ошибка 😞*"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, response)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	f.botConf.bot.Send(message)
}

// showProductModeNoProdsHandler defines the logic of processing the
// empty favorites products array getting.
func (f *favoriteMode) showProductModeNoProdsHandler(chatID int64) {
	response := "*Товаров больше нет 📦*"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message := tgbotapi.NewMessage(chatID, response)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	f.botConf.bot.Send(message)
}

// showFavoriteProducts defines the logic of handling the favorite products request.
func (f *favoriteMode) showFavoriteProducts(chatID int64) {
	f.botConf.users[chatID].lastAction = showFavoriteProducts

	products, err := f.repo.GetFavoriteProducts(context.Background(), chatID)

	if err != nil {
		f.showProductModeGettingErrHandler(chatID, err)
		return

	} else if len(f.botConf.users[chatID].favorites.favoriteLastProdsID) == len(products) {
		f.showProductModeNoProdsHandler(chatID)
		return
	}

	f.showProduct(chatID, products)
}

// deleteFavoriteProduct defines the logic of the deleting the user's favorite product.
func (f *favoriteMode) deleteFavoriteProduct(chatID int64) {
	const op = "tgbot.delete-favorite-product"

	f.botConf.users[chatID].lastAction = deleteFavoriteProduct

	response := "*Товар был успешно удалён 🗑️*"

	err := f.repo.DeleteFavoriteProducts(context.Background(), chatID, []int{f.botConf.users[chatID].favorites.lastFavoriteProdID})

	if err != nil {
		f.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		response = "*Упс... Похоже, произошла ошибка 😞*"
	} else {
		delete(f.botConf.users[chatID].favorites.favoriteLastProdsID,
			f.botConf.users[chatID].favorites.lastFavoriteProdID)
	}

	message := tgbotapi.NewMessage(chatID, response)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Следующий товар ➡️", showFavoriteProducts)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Меню 📋", menuAction)),
	)

	message.ReplyMarkup = keyboard
	message.ParseMode = markDown

	f.botConf.bot.Send(message)
}
