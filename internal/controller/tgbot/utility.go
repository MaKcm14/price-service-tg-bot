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

// usersSampleConfig defines the logic of the current users' products sample interaction and
// iteration.
type usersSampleConfig struct {
	usersSamplePtr        map[int64]map[string]int
	usersLastMarketChoice map[int64]string
	usersSample           map[int64]map[string]entities.ProductSample
}

func newUsersSampleConfig() usersSampleConfig {
	return usersSampleConfig{
		usersSamplePtr:        make(map[int64]map[string]int),
		usersLastMarketChoice: make(map[int64]string),
		usersSample:           make(map[int64]map[string]entities.ProductSample),
	}
}

// usersFavoritesConfig defines the logic of the current users' favorite products interaction and
// iteration.
type usersFavoritesConfig struct {
	usersFavoriteLastProdsID map[int64]map[int]struct{}
	lastFavoriteProdID       map[int64]int
}

func newUsersFavoritesConfig() usersFavoritesConfig {
	return usersFavoritesConfig{
		usersFavoriteLastProdsID: make(map[int64]map[int]struct{}),
		lastFavoriteProdID:       make(map[int64]int),
	}
}

// usersConfigs defines the logic of the user's using configurations.
type usersConfigs struct {
	usersLastAction map[int64]string
	usersRequest    map[int64]dto.ProductRequest

	sampleConfig   usersSampleConfig
	favoriteConfig usersFavoritesConfig
}

func newUsersConfigs() *usersConfigs {
	return &usersConfigs{
		usersLastAction: make(map[int64]string),
		usersRequest:    make(map[int64]dto.ProductRequest),
		sampleConfig:    newUsersSampleConfig(),
		favoriteConfig:  newUsersFavoritesConfig(),
	}
}
