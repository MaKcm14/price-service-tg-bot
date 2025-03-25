package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/MaKcm14/price-service/pkg/entities"
)

// favoriteProdsRepo defines the actions with the favorite products.
type favoriteProdsRepo struct {
	conf redisConfig
}

func newFavoriteProdsRepo(conf redisConfig) favoriteProdsRepo {
	return favoriteProdsRepo{
		conf: conf,
	}
}

// GetUserFavoriteProducts gets the favorite user's products from the cache (cache-hit)
// and error if there're no products in the cache (cache-miss).
func (f favoriteProdsRepo) GetUserFavoriteProducts(ctx context.Context, chatID int64) (map[int]entities.Product, error) {
	const op = "redis.getting-favorite-products"

	res := make(map[int]entities.Product, 100)

	cacheProds, err := f.conf.conn.HGetAll(ctx, fmt.Sprintf("%s%d", favProductsKeyTemplate, chatID)).Result()

	if err != nil {
		f.conf.log.Error(fmt.Sprintf("error of the %s: %s", op, err))
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
func (f favoriteProdsRepo) AddUserFavoriteProducts(ctx context.Context, chatID int64, products map[int]entities.Product) error {
	const op = "redis.adding-favorite-products"

	addProds := make([]interface{}, 0, len(products))

	for id, product := range products {
		jsonProdView, _ := json.Marshal(product)
		addProds = append(addProds, fmt.Sprint(id), string(jsonProdView))
	}

	res := f.conf.conn.HSet(ctx, fmt.Sprintf("%s%d", favProductsKeyTemplate, chatID), addProds...)

	if res.Err() != nil {
		f.conf.log.Error(fmt.Sprintf("error of the %s: %s", op, res.Err()))
		return ErrOfRedisRequest
	}

	return nil
}

// DeleteUserFavoriteProducts deletes the products from the cache for the current user.
func (f favoriteProdsRepo) DeleteUserFavoriteProducts(ctx context.Context, chatID int64, productIDs []int) error {
	const op = "redis.deleting-favorite-products"

	ids := make([]string, 0, len(productIDs))

	for _, id := range productIDs {
		ids = append(ids, fmt.Sprint(id))
	}

	_, err := f.conf.conn.HDel(ctx, fmt.Sprintf("%s%d", favProductsKeyTemplate, chatID), ids...).Result()

	if err != nil {
		f.conf.log.Error("error of the %s: %s", op, err)
		return ErrOfRedisRequest
	}

	return nil
}
