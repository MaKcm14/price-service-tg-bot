package services

import (
	"context"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
)

type (
	// RepoValidator defines the check-actions.
	RepoValidator interface {
		IsUserExists(ctx context.Context, chatID int64) (bool, error)
	}

	// RepoUpdater defines the modify-actions.
	RepoUpdater interface {
		AddUser(ctx context.Context, chatID int64) error
		AddFavoriteProducts(ctx context.Context, chatID int64, products []entities.Product) error
		DeleteFavoriteProducts(ctx context.Context, chatID int64, products []int) error
	}

	RepoAdder interface {
		GetFavoriteProducts(ctx context.Context, chatID int64) (map[int]entities.Product, error)
	}

	Repository interface {
		RepoValidator
		RepoUpdater
		RepoAdder
	}

	ApiInteractor interface {
		GetProductsByBestPrice(request dto.ProductRequest) (map[string]entities.ProductSample, error)
	}
)
