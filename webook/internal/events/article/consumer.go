package article

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"webook/webook/internal/repository"
	"webook/webook/pkg/logger"
	"webook/webook/pkg/saramax"
)

type InteractiveReadEventConsumer struct {
	repo   repository.InteractiveRepository
	client sarama.Client

	l logger.Logger
}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository,
	client sarama.Client, l logger.Logger) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		repo:   repo,
		client: client,
		l:      l,
	}
}

func (i *InteractiveReadEventConsumer) Start() error {
	consumer, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}

	go func() {
		er := consumer.Consume(context.Background(),
			[]string{"article_read"},
			saramax.NewBatchHandler[ReadEvent](i.l, i.BatchConsume),
		)
		if er != nil {
			i.l.Error("Failed to consume",
				logger.Error(er),
			)
		}
	}()
	return err
}

func (i *InteractiveReadEventConsumer) Start_V_Sigle() error {
	consumer, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}

	go func() {
		er := consumer.Consume(context.Background(),
			[]string{"article_read"},
			saramax.NewHandler[ReadEvent](i.l, i.Consume),
		)
		if er != nil {
			i.l.Error("Failed to consume",
				logger.Error(er),
			)
		}
	}()
	return err
}

func (i *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage,
	event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return i.repo.IncreaseViewCount(ctx, "article", event.Aid)
}

func (i *InteractiveReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage,
	ts []ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	bizs := make([]string, 0, len(ts))
	bizIds := make([]int64, 0, len(ts))
	for _, t := range ts {
		bizs = append(bizs, "article")
		bizIds = append(bizIds, t.Aid)
	}
	return i.repo.IncreaseViewCountBatch(ctx, bizs, bizIds)
}
