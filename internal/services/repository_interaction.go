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
	}

	Repository interface {
		RepoValidator
		RepoUpdater
	}

	ApiInteractor interface {
		GetProductsByBestPrice(request dto.ProductRequest) (map[string]entities.ProductSample, error)
	}
)
