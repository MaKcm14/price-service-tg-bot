package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/redis"
	"github.com/MaKcm14/price-service/pkg/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
		p.logger.Warn(fmt.Sprintf("error of the %v: %v", op, err))
		return false, ErrQueryExec
	}
	defer rows.Close()

	return rows.Next(), nil
}

// AddUser adds the user with the chatID to the DB.
func (p PostgreSQLRepo) AddUser(ctx context.Context, chatID int64) error {
	const op = "postgres.adding-user"

	_, err := p.pool.Exec(ctx, "INSERT INTO users (telegram_id) VALUES ($1)", chatID)

	if err != nil {
		p.logger.Warn(fmt.Sprintf("error of the %v: %v", op, err))
		return ErrQueryExec
	}

	return nil
}

// GetFavoriteProducts gets the products from the repository.
func (p PostgreSQLRepo) GetFavoriteProducts(ctx context.Context, chatID int64) (map[int]entities.Product, error) {
	const op = "postgres.getting-favorite-products"

	prods := make(map[int]entities.Product, 100)

	if flagExpired, err := p.cache.IsKeyExpired(ctx, chatID); err != nil {
		p.logger.Error(fmt.Sprintf("%s: error of cache: %s", op, err))
	} else if !flagExpired {
		prods, err := p.cache.GetUserFavoriteProducts(ctx, chatID)

		if err == nil {
			return prods, nil
		}

		p.logger.Error(fmt.Sprintf("%s: error of cache: %s", op, err))
	}

	// Add here the query to the Postgres and setting the products for the user.

	return prods, nil
}

// AddFavoriteProducts adds the new products to the favorites for the current user.
func (p PostgreSQLRepo) AddFavoriteProducts(ctx context.Context, chatID int64, products map[int]entities.Product) error {
	return nil
}

// DeleteFavoriteProducts deletes the products from the favorites of the current user.
func (p PostgreSQLRepo) DeleteFavoriteProducts(ctx context.Context, chatID int64, products map[int]entities.Product) error {
	return nil
}

func (p PostgreSQLRepo) Close() {
	p.pool.Close()
}
