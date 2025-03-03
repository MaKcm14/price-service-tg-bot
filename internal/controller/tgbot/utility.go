package tgbot

import (
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
)

const (
	// markDown sets the Markdown parse mode.
	markDown = "Markdown"

	// Consts set the markets for the request
	megamarket  = "megamarket"
	wildberries = "wildberries"

	// Command's name.
	startAction = "start"
	menuAction  = "menu"
	showRequest = "show-request"

	// Button's data name.
	bestPriceModeData     = "best-price"
	startSearch           = "start-search"
	favoriteModeData      = "favourite-products"
	trackedModeData       = "tracked-products"
	marketSetterMode      = "markets-setter"
	productSetter         = "product-setter"
	productsIter          = "products-iter"
	addToFavorite         = "add-favorite"
	showFavoriteProducts  = "show-favorite-products"
	deleteFavoriteProduct = "delete-favorite-product"
	nextProduct           = "next-product"
)

// userSampleConfig defines the logic of the current user's products sample interaction and
// iteration.
type userSampleConfig struct {
	samplePtr        map[string]int
	lastMarketChoice string
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
