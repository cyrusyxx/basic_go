package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"webook/webook/pkg/logger"
)

type Handler[T any] struct {
	l  logger.Logger
	fn func(msg *sarama.ConsumerMessage, event T) error
}

func NewHandler[T any](l logger.Logger, fn func(msg *sarama.ConsumerMessage, event T) error) *Handler[T] {

	return &Handler[T]{l: l, fn: fn}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		// Unmarshal the message
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

		// Process the message
		err = h.fn(msg, t)
		if err != nil {
			h.l.Error("Failed to process message",
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.Error(err),
			)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
