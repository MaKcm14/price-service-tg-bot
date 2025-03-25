package redis

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

const (
	// favProductsKeyTemplate is the template for every user's key in the cache.
	favProductsKeyTemplate = "products:"
)

// redisConfig defines the main redis configuration.
type redisConfig struct {
	log  *slog.Logger
	conn *redis.Client
}

func newRedisConfig(log *slog.Logger, conn *redis.Client) redisConfig {
	return redisConfig{
		log:  log,
		conn: conn,
	}
}
