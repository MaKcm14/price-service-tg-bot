package redis

import (
	"context"
	"fmt"
	"time"
)

// userRepo defines the actions with the user's cache.
type userRepo struct {
	favoriteProdsRepo
}

// FlushUserCache cleares the user's cache.
func (u userRepo) FlushUserCache(ctx context.Context, chatID int64) error {
	const op = "redis.flush-user-cache"

	_, err := u.conf.conn.Del(ctx, fmt.Sprintf("%s%d", favProductsKeyTemplate, chatID)).Result()

	if err != nil {
		u.conf.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrOfRedisRequest
	}

	return nil
}

// IsKeyExpired checks whether the products-cache for the current user was expired.
func (u userRepo) IsKeyExpired(ctx context.Context, chatID int64) (bool, error) {
	const op = "redis.check-key-expired"

	if flagExists, err := u.conf.conn.Exists(ctx, fmt.Sprintf("%s%d", favProductsKeyTemplate, chatID)).Result(); err != nil {
		u.conf.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return false, ErrOfRedisRequest
	} else if flagExists > 0 {
		return false, nil
	}

	return true, nil
}

// SetTTL sets the needed TTL for the current user's key.
func (u userRepo) SetTTL(ctx context.Context, chatID int64, ttl time.Duration) error {
	const op = "redis.set-ttl"

	_, err := u.conf.conn.Expire(ctx, fmt.Sprintf("%s%d", favProductsKeyTemplate, chatID), ttl).Result()

	if err != nil {
		u.conf.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrOfRedisRequest
	}

	return nil
}
