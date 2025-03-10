package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
)

const (
	// markDown sets the Markdown parse mode.
	markDown = "Markdown"

	// the markets for the request.
	megamarket  = "megamarket"
	wildberries = "wildberries"

	// command's name.
	startAction = "start"
	menuAction  = "menu"
	showRequest = "show-request"

	// button's data name.
	trackedModeData = "tracked-products"

	bestPriceModeData = "best-price"

	startSearch      = "start-search"
	marketSetterMode = "markets-setter"
	productSetter    = "product-setter"
	productsIter     = "products-iter"
	nextProduct      = "next-product"

	favoriteModeData      = "favourite-products"
	addToFavorite         = "add-favorite"
	showFavoriteProducts  = "show-favorite-products"
	deleteFavoriteProduct = "delete-favorite-product"
)

// userSampleConfig defines the logic of the current user's products sample interaction and
// iteration.
type userSampleConfig struct {
	lastMarketChoice string
	samplePtr        map[string]int
	sample           map[string]entities.ProductSample
}

func newUserSampleConfig() userSampleConfig {
	return userSampleConfig{
		samplePtr: make(map[string]int),
		sample:    make(map[string]entities.ProductSample),
	}
}

// userFavoritesConfig defines the logic of the current user's favorite products interaction and
// iteration.
type userFavoritesConfig struct {
	favoriteLastProdsID map[int]struct{}
	lastFavoriteProdID  int
}

func newUserFavoritesConfig() userFavoritesConfig {
	return userFavoritesConfig{
		favoriteLastProdsID: make(map[int]struct{}),
	}
}

// userConfig defines the main configs of the bot's user.
type userConfig struct {
	lastAction string
	request    dto.ProductRequest

	sample    userSampleConfig
	favorites userFavoritesConfig
}

func newUserConfig() *userConfig {
	return &userConfig{
		sample:    newUserSampleConfig(),
		favorites: newUserFavoritesConfig(),
	}
}

// tgBotConfigs defines the main logic of the bot and users' configuration.
type tgBotConfigs struct {
	users map[int64]*userConfig
	bot   *tgbotapi.BotAPI
}

func newTgBotConfigs(bot *tgbotapi.BotAPI) *tgBotConfigs {
	return &tgBotConfigs{
		bot:   bot,
		users: make(map[int64]*userConfig),
	}
}
