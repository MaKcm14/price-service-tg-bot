package kafka

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
)

// Consumer defines the logic of kafka's reading messages.
type Consumer struct {
	cons   sarama.ConsumerGroup
	logger *slog.Logger
}

func NewConsumer(brokers []string, log *slog.Logger) (Consumer, error) {
	const op = "kafka.new-consumer"

	conf := sarama.NewConfig()

	conf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(brokers, consumersGroupID, conf)

	if err != nil {
		log.Error(fmt.Sprintf("error of the %s: %s: %s", op, err, ErrKafkaConnection))
		return Consumer{}, fmt.Errorf("error of the %s: %w: %s", op, ErrKafkaConnection, err)
	}

	return Consumer{
		logger: log,
		cons:   consumer,
	}, nil
}

// ReadProducts reads the messages from the Kafka cluster.
func (c Consumer) ReadProducts(ctx context.Context, handler sarama.ConsumerGroupHandler) {
	const op = "kafka.read-messages"

	for {
		if err := c.cons.Consume(ctx, []string{productsTopicName}, handler); err != nil {
			c.logger.Warn(fmt.Sprintf("error of %s: %s", op, err))
		}
	}
}

// Close releases the resources of the consumer.
func (c Consumer) Close() {
	c.cons.Close()
}
