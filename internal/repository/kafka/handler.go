package kafka

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/MaKcm14/best-price-service/price-service-tg-bot/internal/controller/tgbot"
)

type Handler struct {
	logger *slog.Logger
	prods  chan *tgbot.TrackedProduct
}

func NewHandler(log *slog.Logger, prods chan *tgbot.TrackedProduct) *Handler {
	return &Handler{
		logger: log,
		prods:  prods,
	}
}

func (h *Handler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		trackedProduct := &tgbot.TrackedProduct{}

		for _, header := range message.Headers {
			if string(header.Key) == "ChatID" {
				chatID, _ := strconv.ParseInt(string(header.Value), 10, 64)
				trackedProduct.ChatID = chatID
			}
		}

		json.Unmarshal(message.Value, &trackedProduct.Response)

		h.prods <- trackedProduct

		session.Commit()
	}
	return nil
}
