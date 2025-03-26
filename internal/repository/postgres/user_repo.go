package postgres

import (
	"context"
	"fmt"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/redis"
	"github.com/MaKcm14/price-service/pkg/entities"
)

// userRepo defines the actions with the users.
type userRepo struct {
	conf postgresConfig
	favoriteProdsRepo
	trackedProdsRepo
}

func newUserRepo(conf postgresConfig, cache redis.RedisRepo) userRepo {
	return userRepo{
		conf:              conf,
		favoriteProdsRepo: newFavoriteProdsRepo(conf, cache),
		trackedProdsRepo:  trackedProdsRepo{conf},
	}
}

// IsUserExists checks that the user with the chatID exists at the DB.
func (u userRepo) IsUserExists(ctx context.Context, chatID int64) (bool, error) {
	const op = "postgres.check-user"

	rows, err := u.conf.pool.Query(ctx, "SELECT id FROM users WHERE telegram_id=$1", chatID)

	if err != nil {
		u.conf.logger.Error(fmt.Sprintf("error of the %v: %v", op, err))
		return false, ErrQueryExec
	}
	defer rows.Close()

	return rows.Next(), nil
}

// AddUser adds the user with the chatID to the DB.
func (u userRepo) AddUser(ctx context.Context, chatID int64) error {
	const op = "postgres.adding-user"

	_, err := u.conf.pool.Exec(ctx, "INSERT INTO users (telegram_id) VALUES ($1)", chatID)

	if err != nil {
		u.conf.logger.Error(fmt.Sprintf("error of the %v: %v", op, err))
		return ErrQueryExec
	}

	return nil
}

// GetUserID returns the user ID of the current user.
func (u userRepo) GetUserID(ctx context.Context, chatID int64) (int64, error) {
	const op = "postgres.get-user-id"

	var id int64

	res, err := u.conf.pool.Query(ctx, "SELECT id FROM users WHERE telegram_id=$1", chatID)

	if err != nil {
		u.conf.logger.Error(fmt.Sprintf("error of %s: %s", op, err))
		return -1, ErrQueryExec
	}
	defer res.Close()

	if res.Next() {
		res.Scan(&id)
	}

	return id, nil
}

// IsTrackedProductExists checks that the user with the current chatID already set the tracked product.
func (u userRepo) IsTrackedProductExists(ctx context.Context, chatID int64) (bool, error) {
	const op = "postgres.is-tracked-product-exists"

	id, err := u.GetUserID(ctx, chatID)

	if err != nil {
		u.conf.logger.Error(fmt.Sprintf("error of the %v: %v", op, err))
		return false, ErrQueryExec
	}

	return u.isTrackedProductExists(ctx, id)
}

// AddTrackedProduct adds the tracked product to the DB.
func (u userRepo) AddTrackedProduct(ctx context.Context, chatID int64, request dto.ProductRequest) error {
	const op = "postgres.add-tracked-product"

	id, err := u.GetUserID(ctx, chatID)

	if err != nil {
		u.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrQueryExec
	}

	return u.addTrackedProduct(ctx, id, request)
}

// DeleteTrackedProduct deletes the tracked product of the current chatID.
func (u userRepo) DeleteTrackedProduct(ctx context.Context, chatID int64) error {
	const op = "postgres.delete-tracked-product"

	id, err := u.GetUserID(ctx, chatID)
	if err != nil {
		u.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrQueryExec
	}

	return u.deleteTrackedProduct(ctx, id)
}

// GetTrackedProduct gets the tracked product for the current chatID.
func (u userRepo) GetTrackedProduct(ctx context.Context, chatID int64) (dto.ProductRequest, bool, error) {
	const op = "postgres.get-tracked-product"

	id, err := u.GetUserID(ctx, chatID)
	if err != nil {
		u.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return dto.ProductRequest{}, false, ErrQueryExec
	}

	return u.getTrackedProduct(ctx, id)
}

// GetFavoriteProducts gets the products from the repository.
func (u userRepo) GetFavoriteProducts(ctx context.Context, chatID int64) (map[int]entities.Product, error) {
	const op = "postgres.getting-favorite-products"

	id, err := u.GetUserID(ctx, chatID)

	if err != nil {
		u.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return nil, ErrQueryExec
	}

	return u.getFavoriteProducts(ctx, userID{
		userID: id,
		chatID: chatID,
	})
}

// AddFavoriteProducts adds the new products to the favorites for the current user.
func (u userRepo) AddFavoriteProducts(ctx context.Context, chatID int64, products []entities.Product) error {
	const op = "postgres.add-favorite-products"

	id, err := u.GetUserID(ctx, chatID)

	if err != nil {
		u.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrQueryExec
	}

	return u.addFavoriteProducts(ctx, userID{
		userID: id,
		chatID: chatID,
	}, products)
}
