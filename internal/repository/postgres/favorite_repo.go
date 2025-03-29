package postgres

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/redis"
	"github.com/MaKcm14/price-service/pkg/entities"
)

// favoriteProdsRepo defines the actions with the favorite products.
type favoriteProdsRepo struct {
	conf  postgresConfig
	cache redis.RedisRepo
}

func newFavoriteProdsRepo(conf postgresConfig, cache redis.RedisRepo) favoriteProdsRepo {
	return favoriteProdsRepo{
		conf:  conf,
		cache: cache,
	}
}

// getUserProducts gets the products from the DB for the current user_id.
func (f favoriteProdsRepo) getUserProducts(ctx context.Context, id userID) (map[int]entities.Product, error) {
	const op = "postgres.getting-products"

	prods := make(map[int]entities.Product, 100)

	res, err := f.conf.pool.Query(ctx, "SELECT f.product_id, f.product_name, f.product_link, f.base_price, f.product_brand, f.supplier "+
		"FROM users as u JOIN favorites as f ON f.user_id=u.id and u.id=$1", id.userID)

	if err != nil {
		f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
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
		f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, res.Err()))
		return nil, ErrQueryExec
	}

	return prods, nil
}

// updateCacheInfo updates the cache info about the products in the main-cache.
func (f favoriteProdsRepo) updateCacheInfo(ctx context.Context, id userID, prods map[int]entities.Product) error {
	const op = "postgres.update-cache-info"

	err := f.cache.AddUserFavoriteProducts(ctx, id.chatID, prods)

	if err != nil {
		f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return err
	}

	err = f.cache.SetTTL(ctx, id.chatID, time.Duration(time.Hour*5))

	if err != nil {
		f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
	}

	return err
}

// getFavoriteProducts gets the products from the repository for the current user_id.
func (f favoriteProdsRepo) getFavoriteProducts(ctx context.Context, id userID) (map[int]entities.Product, error) {
	const op = "postgres.get-favorite-products"

	if flagExpired, err := f.cache.IsKeyExpired(ctx, id.chatID); err != nil {
		f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
	} else if !flagExpired {
		products, err := f.cache.GetUserFavoriteProducts(ctx, id.chatID)

		if err == nil {
			return products, nil
		}
		f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
	}

	prods, err := f.getUserProducts(ctx, id)

	if err != nil {
		f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return nil, fmt.Errorf("error of the %s: %w", op, err)
	}

	if len(prods) != 0 {
		err = f.updateCacheInfo(ctx, id, prods)

		if err != nil {
			f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
			return nil, fmt.Errorf("error of the %s: %w", op, err)
		}
	}

	return prods, nil
}

// addFavoriteProducts adds the new products to the favorites for the current user_id.
func (f favoriteProdsRepo) addFavoriteProducts(ctx context.Context, id userID, products []entities.Product) error {
	const op = "postgres.add-favorite-products"
	const query = "INSERT INTO favorites (product_name, product_link, base_price, product_brand, supplier, user_id)\n"

	err := f.cache.FlushUserCache(ctx, id.chatID)

	if err != nil {
		f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
	}

	for _, product := range products {

		_, err = f.conf.pool.Exec(ctx, fmt.Sprintf("%sVALUES ('%s', '%s', %d, '%s', '%s', %d)", query,
			product.Name, product.Links.URL, product.Price.BasePrice,
			product.Brand, product.Supplier, id.userID))

		if err != nil {
			f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
			return ErrQueryExec
		}
	}

	return nil
}

// DeleteFavoriteProducts deletes the products from the favorites of the current user.
func (f favoriteProdsRepo) DeleteFavoriteProducts(ctx context.Context, chatID int64, products []int) error {
	const op = "postgres.delete-favorite-products"

	query := bytes.Buffer{}
	query.WriteString("DELETE FROM favorites WHERE product_id in (")

	if flagExpired, err := f.cache.IsKeyExpired(ctx, chatID); err != nil {
		f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
	} else if !flagExpired {
		err := f.cache.DeleteUserFavoriteProducts(ctx, chatID, products)

		if err != nil {
			f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		}
	}

	for i, id := range products {
		var str = fmt.Sprintf("%d,", id)

		if i == len(products)-1 {
			str = fmt.Sprintf("%d)", id)
		}
		query.WriteString(str)
	}

	_, err := f.conf.pool.Exec(ctx, query.String())

	if err != nil {
		f.conf.logger.Error(fmt.Sprintf("error of the %s: %s", op, err))
		return ErrQueryExec
	}

	return nil
}
