package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/MaKcm14/price-service/pkg/entities"
	"github.com/redis/go-redis/v9"
)

const (
	keyTemplate = "products:"
)

type RedisRepo struct {
	log *slog.Logger

	conn *redis.Client
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

	return RedisRepo{
		log:  logger,
		conn: client,
	}, nil
}

// IsKeyExpired checks whether the products-cache for the current user was expired.
func (r RedisRepo) IsKeyExpired(ctx context.Context, tgID int64) (bool, error) {
	const op = "redis.check-key-expired"

	if flagExists, err := r.conn.Exists(ctx, fmt.Sprintf("%s%d", keyTemplate, tgID)).Result(); err != nil {
		r.log.Error(fmt.Sprintf("error of redis interaction: %s: %s", op, err))
		return false, ErrOfRedisRequest
	} else if flagExists > 0 {
		return false, nil
	}

	return true, nil
}

// FlushUserCache cleares the user's cache.
func (r RedisRepo) FlushUserCache(ctx context.Context, tgID int64) error {
	const op = "redis.flush-user-cache"

	_, err := r.conn.Del(ctx, fmt.Sprintf("%s%d", keyTemplate, tgID)).Result()

	if err != nil {
		r.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrOfRedisRequest
	}

	return nil
}

// SetTTL sets the TTL for the current user's key.
func (r RedisRepo) SetTTL(ctx context.Context, tgID int64) error {
	const op = "redis.set-ttl"

	_, err := r.conn.Expire(ctx, fmt.Sprintf("%s%d", keyTemplate, tgID), time.Duration(time.Hour*72)).Result()

	if err != nil {
		r.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrOfRedisRequest
	}

	return nil
}

// GetUserFavoriteProducts gets the favorite user's products from cache (cache-hit) and error if
// there're no products in the cache (cache-miss).
func (r RedisRepo) GetUserFavoriteProducts(ctx context.Context, tgID int64) (map[int]entities.Product, error) {
	const op = "redis.getting-favorite-products"

	res := make(map[int]entities.Product, 100)

	cacheProds, err := r.conn.HGetAll(ctx, fmt.Sprintf("%s%d", keyTemplate, tgID)).Result()

	if err != nil {
		r.log.Error(fmt.Sprintf("error of redis interaction: %s: %s", op, err))
		return nil, ErrOfRedisRequest
	}

	for id, prod := range cacheProds {
		var product entities.Product

		json.Unmarshal([]byte(prod), &product)

		productID, _ := strconv.Atoi(id)
		res[productID] = product
	}

	return res, nil
}

// AddUserFavoriteProducts adds the favorite user's products to the hash in the cache.
func (r RedisRepo) AddUserFavoriteProducts(ctx context.Context, tgID int64, products map[int]entities.Product) error {
	const op = "redis.adding-favorite-products"

	addProds := make([]interface{}, 0, len(products))

	for id, product := range products {
		jsonProdView, _ := json.Marshal(product)
		addProds = append(addProds, fmt.Sprint(id), string(jsonProdView))
	}

	res := r.conn.HSet(ctx, fmt.Sprintf("%s%d", keyTemplate, tgID), addProds...)

	if res.Err() != nil {
		r.log.Error(fmt.Sprintf("error of redis interaction: %s: %s", op, res.Err()))
		return ErrOfRedisRequest
	}

	return nil
}

// DeleteFavoriteProducts deletes the products from the cache for the current user.
func (r RedisRepo) DeleteFavoriteProducts(ctx context.Context, tgID int64, productIDs []int) error {
	const op = "redis.deleting-favorite-products"

	ids := make([]string, 0, len(productIDs))

	for _, id := range productIDs {
		ids = append(ids, fmt.Sprint(id))
	}

	_, err := r.conn.HDel(ctx, fmt.Sprintf("%s%d", keyTemplate, tgID), ids...).Result()

	if err != nil {
		r.log.Error("error of redis interaction: %s: %s", op, err)
		return ErrOfRedisRequest
	}

	return nil
}
