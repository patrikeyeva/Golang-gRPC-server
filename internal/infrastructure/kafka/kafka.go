//go:generate mockgen -source ./kafka.go -destination=./mocks/kafka.go -package=mock_kafka
package kafka

import "github.com/Shopify/sarama"

type ProducerKafka interface {
	SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error)
}
