package postgres

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/redis"
	"github.com/MaKcm14/price-service/pkg/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgreSQLRepo defines the logic of the DB interaction.
type PostgreSQLRepo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger

	cache redis.RedisRepo
}

func New(ctx context.Context, dsn string, log *slog.Logger) (PostgreSQLRepo, error) {
	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		log.Error(fmt.Sprintf("error of connecting to the DB: %v", err))
		return PostgreSQLRepo{}, ErrDBConnection
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		log.Error(fmt.Sprintf("error of connecting to the DB: %v", err))
		return PostgreSQLRepo{}, ErrDBConnection
	}

	cache, err := redis.New(ctx, log)

	if err != nil {
		log.Error(fmt.Sprintf("error of starting the redis: %s", err))
		return PostgreSQLRepo{}, ErrCacheConnection
	}

	return PostgreSQLRepo{
		pool:   pool,
		logger: log,
		cache:  cache,
	}, nil
}

// IsUserExists checks that the user with the chatID exists at the DB.
func (p PostgreSQLRepo) IsUserExists(ctx context.Context, chatID int64) (bool, error) {
	const op = "postgres.check-user"

	rows, err := p.pool.Query(ctx, "SELECT id FROM users WHERE telegram_id=$1", chatID)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %v: %v", op, err))
		return false, ErrQueryExec
	}
	defer rows.Close()

	return rows.Next(), nil
}

// GetUserID returns the user ID of the current user.
func (p PostgreSQLRepo) GetUserID(ctx context.Context, chatID int64) (int, error) {
	const op = "postgres.get-user-id"

	var id int

	res, err := p.pool.Query(ctx, "SELECT id FROM users WHERE telegram_id=$1", chatID)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of %s: %s", op, err))
		return -1, ErrQueryExec
	}
	defer res.Close()

	if res.Next() {
		res.Scan(&id)
	}

	return id, nil
}

// AddUser adds the user with the chatID to the DB.
func (p PostgreSQLRepo) AddUser(ctx context.Context, chatID int64) error {
	const op = "postgres.adding-user"

	_, err := p.pool.Exec(ctx, "INSERT INTO users (telegram_id) VALUES ($1)", chatID)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %v: %v", op, err))
		return ErrQueryExec
	}

	return nil
}

// getUserProducts gets the products from the DB for the current chatID.
func (p PostgreSQLRepo) getUserProducts(ctx context.Context, chatID int64) (map[int]entities.Product, error) {
	const op = "postgres.getting-products"

	prods := make(map[int]entities.Product, 100)

	id, err := p.GetUserID(ctx, chatID)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return nil, ErrQueryExec
	}

	res, err := p.pool.Query(ctx, "SELECT f.product_id, f.product_name, f.product_link, f.base_price, f.product_brand, f.supplier "+
		"FROM users as u JOIN favorites as f ON f.user_id=u.id and u.id=$1", id)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return nil, ErrQueryExec
	}
	defer res.Close()

	for res.Next() {
		var product entities.Product
		var productID int

		res.Scan(&productID, &product.Name, &product.Links.URL,
			&product.Price.BasePrice, &product.Brand, &product.Supplier)

		prods[productID] = product
	}

	if res.Err() != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, res.Err()))
		return nil, ErrQueryExec
	}

	return prods, nil
}

// updateCacheInfo updates the cache info about the products in the main-cache.
func (p PostgreSQLRepo) updateCacheInfo(ctx context.Context, chatID int64, prods map[int]entities.Product) error {
	const op = "postgres.update-cache-info"

	err := p.cache.AddUserFavoriteProducts(ctx, chatID, prods)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return err
	}

	err = p.cache.SetTTL(ctx, chatID, time.Duration(time.Hour*5))

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
	}

	return err
}

// GetFavoriteProducts gets the products from the repository.
func (p PostgreSQLRepo) GetFavoriteProducts(ctx context.Context, chatID int64) (map[int]entities.Product, error) {
	const op = "postgres.getting-favorite-products"

	if flagExpired, err := p.cache.IsKeyExpired(ctx, chatID); err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
	} else if !flagExpired {
		products, err := p.cache.GetUserFavoriteProducts(ctx, chatID)

		if err == nil {
			return products, nil
		}
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
	}

	prods, err := p.getUserProducts(ctx, chatID)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return nil, fmt.Errorf("error of the %s: %w", op, err)
	}

	if len(prods) != 0 {
		err = p.updateCacheInfo(ctx, chatID, prods)

		if err != nil {
			p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
			return nil, fmt.Errorf("error of the %s: %w", op, err)
		}
	}

	return prods, nil
}

// AddFavoriteProducts adds the new products to the favorites for the current user.
func (p PostgreSQLRepo) AddFavoriteProducts(ctx context.Context, chatID int64, products []entities.Product) error {
	const op = "postgres.add-favorite-products"
	const query = "INSERT INTO favorites (product_name, product_link, base_price, product_brand, supplier, user_id)\n"

	err := p.cache.FlushUserCache(ctx, chatID)

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
	}

	for _, product := range products {
		id, err := p.GetUserID(ctx, chatID)

		if err != nil {
			p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
			return ErrQueryExec
		}

		_, err = p.pool.Exec(ctx, fmt.Sprintf("%sVALUES ('%s', '%s', %d, '%s', '%s', %d)", query,
			product.Name, product.Links.URL, product.Price.BasePrice,
			product.Brand, product.Supplier, id))

		if err != nil {
			p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
			return ErrQueryExec
		}
	}

	return nil
}

// DeleteFavoriteProducts deletes the products from the favorites of the current user.
func (p PostgreSQLRepo) DeleteFavoriteProducts(ctx context.Context, chatID int64, products []int) error {
	const op = "postgres.delete-favorite-products"

	query := bytes.Buffer{}
	query.WriteString("DELETE FROM favorites WHERE product_id in (")

	if flagExpired, err := p.cache.IsKeyExpired(ctx, chatID); err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
	} else if !flagExpired {
		err := p.cache.DeleteFavoriteProducts(ctx, chatID, products)

		if err != nil {
			p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		}
	}

	for i, id := range products {
		var str = fmt.Sprintf("%d,", id)

		if i == len(products)-1 {
			str = fmt.Sprintf("%d)", id)
		}
		query.WriteString(str)
	}

	_, err := p.pool.Exec(ctx, query.String())

	if err != nil {
		p.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrQueryExec
	}

	return nil
}

// Close defines the releasing the resources of the DB.
func (p PostgreSQLRepo) Close() {
	p.pool.Close()
}
