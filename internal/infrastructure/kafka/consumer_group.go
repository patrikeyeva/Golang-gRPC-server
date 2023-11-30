package kafka

import (
	"encoding/json"
	"homework8/pkg/logger"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type serverMessage struct {
	CreationTime  time.Time
	MethodType    string
	RequestMethod string
	RequestBody   string
}

type ConsumerGroup struct {
	ready chan bool
}

func CreateConfigForConsumerGroup() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.ResetInvalidOffsets = true

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}

	return config

}

func NewConsumerGroup() ConsumerGroup {
	return ConsumerGroup{
		ready: make(chan bool),
	}
}

func (consumer *ConsumerGroup) Ready() <-chan bool {
	return consumer.ready
}

func (consumer *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)

	return nil
}

func (consumer *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	l := logger.FromContext(session.Context())
	ctx := logger.ToContext(session.Context(), l.With(zap.String("method", "ConsumeClaim")))

	for {
		select {
		case message := <-claim.Messages():

			pm := serverMessage{}
			err := json.Unmarshal(message.Value, &pm)
			if err != nil {
				logger.Errorf(ctx, "Consumer group error", err)
			}

			logger.Infof(ctx, "Received message:value = %v,timestamp = %v,topic = %s",
				pm,
				message.Timestamp,
				message.Topic)

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
