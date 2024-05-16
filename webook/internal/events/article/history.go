package article

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository"
	"webook/webook/pkg/logger"
	"webook/webook/pkg/saramax"
)

type HistoryConsumer struct {
	repo   repository.HistoryRepository
	client sarama.Client
	l      logger.Logger
}

func NewHistoryConsumer(repo repository.HistoryRepository,
	client sarama.Client, l logger.Logger) *HistoryConsumer {
	return &HistoryConsumer{repo: repo, client: client, l: l}
}

func (c *HistoryConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(TopicReadEvent, c.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{TopicReadEvent},
			saramax.NewHandler[ReadEvent](c.l, c.Consume),
		)
		if er != nil {
			c.l.Error("Error consuming history", logger.Error(er))
		}
	}()
	return err
}

func (c *HistoryConsumer) Consume(msg *sarama.ConsumerMessage, event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return c.repo.AddRecord(ctx, domain.HistoryRecord{
		Biz:   "article",
		BizId: event.Aid,
		Uid:   event.Uid,
	})
}
