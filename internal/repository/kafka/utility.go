package kafka

import "github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/api"

const (
	productsTopicName = "products"
	consumersGroupID  = "products-reader"
)

type TrackedProduct struct {
	ChatID   int64
	Response api.ProductResponse
}
