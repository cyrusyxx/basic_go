package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	ProducerReadEvent(evt ReadEvent) error
}

type ReadEvent struct {
	Aid int64
	Uid int64
}

const (
	TopicReadEvent = "article_read"
)

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{producer: producer}
}

func (p *SaramaSyncProducer) ProducerReadEvent(evt ReadEvent) error {

	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicReadEvent,
		Value: sarama.StringEncoder(val),
	})

	return err
}
