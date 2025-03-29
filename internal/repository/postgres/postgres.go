package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/redis"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgreSQLRepo defines the logic of the DB interaction.
type PostgreSQLRepo struct {
	conf postgresConfig
	userRepo
}

func New(ctx context.Context, dsn string, redisConf redis.RedisInitConf, log *slog.Logger) (PostgreSQLRepo, error) {
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

	cache, err := redis.New(ctx, log, redisConf)

	if err != nil {
		log.Error(fmt.Sprintf("error of starting the redis: %s", err))
		return PostgreSQLRepo{}, ErrCacheConnection
	}

	conf := newPostgresConfig(pool, log)

	return PostgreSQLRepo{
		conf:     conf,
		userRepo: newUserRepo(conf, cache),
	}, nil
}

// Close defines the releasing the resources of the DB.
func (p PostgreSQLRepo) Close() {
	p.conf.pool.Close()
}
