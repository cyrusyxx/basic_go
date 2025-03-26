package saramax

import (
	"context"
	"encoding/json"
	"time"
	"webook/webook/pkg/logger"

	"github.com/IBM/sarama"
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

func (h *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10

	for {
		// 在每次循环开始时重置切片
		batch := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		var done bool

		// 收集一批消息
		for i := 0; i < batchSize && !done; i++ {
			select {
			case msg, ok := <-msgs:
				if !ok {
					// 通道关闭，处理完当前批次后返回
					done = true
					cancel()
					if len(batch) > 0 {
						err := h.fn(batch, ts)
						if err != nil {
							h.l.Error("Failed to process final batch",
								logger.Error(err),
							)
						}
						for _, m := range batch {
							session.MarkMessage(m, "")
						}
					}
					return nil
				}
				batch = append(batch, msg)
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					h.l.Error("Failed to unmarshal message",
						logger.String("topic", msg.Topic),
						logger.Int32("partition", msg.Partition),
						logger.Int64("offset", msg.Offset),
						logger.Error(err),
					)
					session.MarkMessage(msg, "") // 标记错误消息为已处理
					continue
				}
				ts = append(ts, t)
			case <-ctx.Done():
				// 超时，处理已收集的消息
				done = true
			}
		}
		cancel()

		// 如果收集到了消息，就处理它们
		if len(batch) > 0 {
			err := h.fn(batch, ts)
			if err != nil {
				h.l.Error("Failed to process batch",
					logger.Error(err),
				)
			}
			// 无论处理是否成功，都标记消息为已处理
			for _, msg := range batch {
				session.MarkMessage(msg, "")
			}
		}

		// 如果通道已关闭，退出循环
		if done {
			return nil
		}
	}
}
