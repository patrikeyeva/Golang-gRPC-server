package reader

import (
	"context"
	"homework8/internal/infrastructure/kafka"
	"homework8/pkg/logger"
	"sync"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

func StartReader(ctx context.Context, clientConsumer sarama.ConsumerGroup) {
	l := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, l.With(zap.String("method", "StartReader")))

	consumer := kafka.NewConsumerGroup()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := clientConsumer.Consume(ctx, []string{"requests"}, &consumer); err != nil {
				logger.Errorf(ctx, "Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.Ready()
	logger.Infof(ctx, "Sarama consumer up and running!...")

	select {
	case <-ctx.Done():
		logger.Infof(ctx, "terminating: context cancelled")
	}

	wg.Wait()
}
