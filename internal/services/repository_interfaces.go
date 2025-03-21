package services

import (
	"context"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/price-service/pkg/entities"
)

type (
	// Validator defines the repository check-actions.
	Validator interface {
		IsUserExists(ctx context.Context, tgID int64) (bool, error)
		IsTrackedProductExists(ctx context.Context, tgID int64) (bool, error)
	}

	// Updater defines the repository-storage modify-actions.
	Modificator interface {
		AddUser(ctx context.Context, tgID int64) error
		AddFavoriteProducts(ctx context.Context, tgID int64, products []entities.Product) error
		AddTrackedProduct(ctx context.Context, chatID int64, request dto.ProductRequest) error

		DeleteFavoriteProducts(ctx context.Context, tgID int64, products []int) error
	}

	// Getter defines the repository-storage read operations.
	Getter interface {
		GetFavoriteProducts(ctx context.Context, tgID int64) (map[int]entities.Product, error)
	}

	// Repository defines the repository-storage object abstraction.
	Repository interface {
		Validator
		Modificator
		Getter
	}

	// ApiInteractor defines the repository price-service-api object abstraction.
	ApiInteractor interface {
		GetProductsByBestPrice(request dto.ProductRequest) (map[string]entities.ProductSample, error)
	}
)
