package server

import (
	"context"
	"encoding/json"
	"fmt"
	"homework8/internal/infrastructure/kafka"
	"homework8/pkg/logger"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type ServerMessage struct {
	CreationTime  time.Time
	MethodType    string
	RequestMethod string
	RequestBody   string
}

type KafkaSender struct {
	producer kafka.ProducerKafka
	topic    string
}

func NewKafkaSender(producer kafka.ProducerKafka, topic string) *KafkaSender {
	return &KafkaSender{
		producer,
		topic,
	}
}

func NewServerMessage(path string, method string, reqBody []byte) (*ServerMessage, error) {
	return &ServerMessage{
		CreationTime:  time.Now(),
		MethodType:    path,
		RequestMethod: method,
		RequestBody:   string(reqBody),
	}, nil

}

func (s *KafkaSender) SendMessage(ctx context.Context, message *sarama.ProducerMessage) error {
	l := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, l.With(zap.String("method", "kafka-SendMessage")))

	partition, offset, err := s.producer.SendSyncMessage(message)
	if err != nil {
		logger.Errorf(ctx, "Send message connector error", err)
		return err
	}

	logger.Infof(ctx, "Send message:\nPartition: %v, Offset: %v, Type: %v", partition, offset, message.Key)
	return nil
}

func BuildMessage(message *ServerMessage, topic string) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(message.MethodType)),
	}, nil
}
