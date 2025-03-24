package kafka

import (
	"encoding/json"
	"fmt"
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
	//DEBUG:
	h.logger.Info(fmt.Sprintf("DEBUG: TopicName: %s; Partitions: %d; Offset: %d", claim.Topic(), claim.Partition(), claim.InitialOffset()))
	//TODO: delete

	for message := range claim.Messages() {
		trackedProduct := &tgbot.TrackedProduct{}

		for _, header := range message.Headers {
			//DEBUG:
			h.logger.Info(fmt.Sprintf("DEBUG: header key: %s, header val: %s", string(header.Key), string(header.Value)))
			//TODO: delete

			if string(header.Key) == "ChatID" {
				chatID, _ := strconv.ParseInt(string(header.Value), 10, 64)
				trackedProduct.ChatID = chatID
			}
		}

		json.Unmarshal(message.Value, &trackedProduct.Response)

		//DEBUG:
		h.logger.Info(fmt.Sprintf("DEBUG: product: %v\n\n", trackedProduct))
		//TODO: delete

		h.prods <- trackedProduct

		session.Commit()
	}
	return nil
}
