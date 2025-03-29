package postgres

import (
	"context"
	"fmt"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/entities/dto"
)

// trackedProdsRepo defines the actions with the tracked products.
type trackedProdsRepo struct {
	conf postgresConfig
}

// isTrackedProductExists checks that the user with the current chatID already set the tracked product.
func (p trackedProdsRepo) isTrackedProductExists(ctx context.Context, id int64) (bool, error) {
	const op = "postgres.is-tracked-product-exists"

	rows, err := p.conf.pool.Query(ctx, "SELECT tracked_id FROM tracked_products WHERE user_id=$1", id)

	if err != nil {
		p.conf.logger.Error(fmt.Sprintf("error of the %v: %v", op, err))
		return false, ErrQueryExec
	}
	defer rows.Close()

	return rows.Next(), nil
}

// addTrackedProduct adds the tracked product to the DB.
func (p trackedProdsRepo) addTrackedProduct(ctx context.Context, id int64, request dto.ProductRequest) error {
	const op = "postgres.add-tracked-product"

	markets := make([]string, 0, 10)

	for _, market := range request.Markets {
		markets = append(markets, market)
	}

	query := fmt.Sprintf("INSERT INTO tracked_products "+
		"(product_name, market_filter, user_id)\n"+
		"VALUES ('%s', $1, %d)", request.Query, id)

	_, err := p.conf.pool.Exec(ctx, query, markets)

	if err != nil {
		p.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrQueryExec
	}

	return nil
}

// deleteTrackedProduct deletes the tracked product of the current user_id.
func (p trackedProdsRepo) deleteTrackedProduct(ctx context.Context, id int64) error {
	const op = "postgres.delete-tracked-product"

	_, err := p.conf.pool.Exec(ctx, "DELETE FROM tracked_products WHERE user_id=$1", id)

	if err != nil {
		p.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrQueryExec
	}

	return nil
}

// getTrackedProduct gets the tracked product for the current user_id.
func (p trackedProdsRepo) getTrackedProduct(ctx context.Context, id int64) (dto.ProductRequest, bool, error) {
	const op = "postgres.get-tracked-product"

	res := dto.NewProductRequest()

	rows, err := p.conf.pool.Query(ctx, "SELECT product_name, market_filter FROM tracked_products WHERE user_id=$1", id)

	if err != nil {
		p.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return dto.ProductRequest{}, false, ErrQueryExec
	}
	defer rows.Close()

	if rows.Next() {
		markets := make([]string, 0, 10)

		rows.Scan(&res.Query, &markets)

		for _, market := range markets {
			res.Markets[market] = market
		}
	} else {
		return dto.ProductRequest{}, false, nil
	}

	if rows.Err() != nil {
		p.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return dto.ProductRequest{}, false, ErrQueryExec
	}

	return res, true, nil
}

// GetUsersTrackedProducts defines the logic of getting the tracked products of every user from the DB.
func (p trackedProdsRepo) GetUsersTrackedProducts(ctx context.Context) (map[int64]dto.ProductRequest, error) {
	const op = "postgres.get-users-tracked-products-async"

	query := fmt.Sprint("SELECT telegram_id, product_name, market_filter ",
		"FROM users AS u JOIN tracked_products AS t ON u.id=t.user_id")

	rows, err := p.conf.pool.Query(ctx, query)

	if err != nil {
		p.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return nil, ErrQueryExec
	}
	defer rows.Close()

	res := make(map[int64]dto.ProductRequest)

	for rows.Next() {
		chatID := int64(0)
		markets := make([]string, 0, 200)
		product := dto.NewProductRequest()

		rows.Scan(&chatID, &product.Query, &markets)

		for _, market := range markets {
			product.Markets[market] = market
		}

		res[chatID] = product
	}

	if rows.Err() != nil {
		p.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return nil, ErrQueryExec
	}

	return res, nil
}
