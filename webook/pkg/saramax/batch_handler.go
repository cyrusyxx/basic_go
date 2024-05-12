package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
	"webook/webook/pkg/logger"
)

type BatchHandler[T any] struct {
	l  logger.Logger
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
}

func NewBatchHandler[T any](l logger.Logger,
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{l: l, fn: fn}
}

func (h *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	ts := make([]T, 0, 10)
	const batchSize = 10

	for {
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		var done = false

		for i := 0; i < batchSize && !done; i++ {
			select {
			case msg, ok := <-msgs:
				// If the channel is closed, we are done
				// And append the message to the batch
				if !ok {
					done = true
					break
				}
				batch = append(batch, msg)

				// Unmarshal the message
				// And append the message to the batch
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					h.l.Error("Failed to unmarshal message",
						logger.String("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset),
						logger.Error(err),
					)
					continue
				}
				ts = append(ts, t)
			case <-ctx.Done():
				// If the context is done, we are done
				done = true
			}
		}
		cancel()

		// Process the batch
		err := h.fn(batch, ts)
		if err != nil {
			h.l.Error("Failed to process batch",
				logger.Error(err),
			)
		}
		for _, msg := range batch {
			session.MarkMessage(msg, "")
		}
	}
}
