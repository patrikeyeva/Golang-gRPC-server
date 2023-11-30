package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

type Producer struct {
	brokers      []string
	syncProducer sarama.SyncProducer
}

func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()

	syncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll
	syncProducerConfig.Producer.Idempotent = true
	syncProducerConfig.Net.MaxOpenRequests = 1
	syncProducerConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault

	syncProducerConfig.Producer.Return.Successes = true
	syncProducerConfig.Producer.Return.Errors = true

	syncProducerConfig.Producer.Compression = sarama.CompressionGZIP

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)

	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka-producer")
	}

	return syncProducer, nil
}

func NewProducer(brokers []string) (*Producer, error) {
	syncProducer, err := newSyncProducer(brokers)
	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka-producer")
	}

	producer := &Producer{
		brokers:      brokers,
		syncProducer: syncProducer,
	}

	return producer, nil
}

func (p *Producer) SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return p.syncProducer.SendMessage(message)
}

func (p *Producer) Close() error {
	err := p.syncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "kafka.Connector.Close")
	}

	return nil
}
