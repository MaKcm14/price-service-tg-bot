package kafka

import (
	"log/slog"

	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/controller/tgbot"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/kafka/hand"
)

const (
	productsTopicName = "products"
	consumersGroupID  = "products-reader"
)

type ConfigHandler func(*Consumer)

// ConfigProductHandler defines the logic of configuring the product's handler of the consumer.
func ConfigProductHandler(log *slog.Logger, prods chan *tgbot.TrackedProduct) ConfigHandler {
	return func(consumer *Consumer) {
		consumer.prodsHandler = hand.NewProductsHandler(log, prods)
	}
}
