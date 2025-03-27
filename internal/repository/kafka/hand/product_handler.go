package hand

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/repository/kafka"
)

// ProductsHandler defines the logic of handling the kafka's messages.
type ProductsHandler struct {
	logger *slog.Logger
	prods  chan *kafka.TrackedProduct
}

func NewProductsHandler(log *slog.Logger, prods chan *kafka.TrackedProduct) *ProductsHandler {
	return &ProductsHandler{
		logger: log,
		prods:  prods,
	}
}

func (p *ProductsHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		trackedProduct := &kafka.TrackedProduct{}

		for _, header := range message.Headers {
			if string(header.Key) == "ChatID" {
				chatID, _ := strconv.ParseInt(string(header.Value), 10, 64)
				trackedProduct.ChatID = chatID
			}
		}

		json.Unmarshal(message.Value, &trackedProduct.Response)

		p.prods <- trackedProduct

		session.Commit()
	}
	return nil
}

func (p *ProductsHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (p *ProductsHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
