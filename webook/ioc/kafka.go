package ioc

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"time"
	"webook/webook/internal/events"
	"webook/webook/internal/events/article"
)

func InitSaramaClient() sarama.Client {
	type Config struct {
		Addr []string
	}
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(cfg.Addr, scfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return producer
}

func InitConsumers(c *article.InteractiveReadEventConsumer) []events.Consumer {
	return []events.Consumer{c}
}

func KafkaHealthCheck(brokers []string) error {
	const testTopic = "server_startup_test"

	// 生产者配置
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true

	// 创建生产者
	producer, err := sarama.NewSyncProducer(brokers, producerConfig)
	if err != nil {
		return fmt.Errorf("生产者创建失败: %v", err)
	}
	defer producer.Close()

	// 发送测试消息
	msg := &sarama.ProducerMessage{
		Topic: testTopic,
		Value: sarama.StringEncoder("healthcheck_" + time.Now().String()),
	}
	if _, _, err := producer.SendMessage(msg); err != nil {
		return fmt.Errorf("消息发送失败: %v", err)
	}

	// 消费者配置
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	// 创建消费者
	consumer, err := sarama.NewConsumer(brokers, consumerConfig)
	if err != nil {
		return fmt.Errorf("消费者创建失败: %v", err)
	}
	defer consumer.Close()

	// 订阅主题
	partitionConsumer, err := consumer.ConsumePartition(testTopic, 0, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("分区消费失败: %v", err)
	}
	defer partitionConsumer.Close()

	// 消费验证 (15秒超时)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	select {
	case msg := <-partitionConsumer.Messages():
		if string(msg.Value) != string(msg.Value) {
			return fmt.Errorf("消息内容不匹配")
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("消费超时")
	}
}
