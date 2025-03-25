package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// RedisRepo defines the logic of the cache interaction.
type RedisRepo struct {
	conf redisConfig
	userRepo
}

func New(ctx context.Context, logger *slog.Logger) (RedisRepo, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if res := client.Ping(ctx); res.Err() != nil {
		logger.Error(fmt.Sprintf("error of ping the redis: %v", res.Err()))
		return RedisRepo{}, ErrConnToRedis
	}

	conf := newRedisConfig(logger, client)

	return RedisRepo{
		conf: conf,
		userRepo: userRepo{
			favoriteProdsRepo: newFavoriteProdsRepo(conf),
		},
	}, nil
}
