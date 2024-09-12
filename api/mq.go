package api

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"log"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/lvkeliang/WHOIM/config"
)

func InitConsumer(serverAddress string, cfg config.Config) (rocketmq.PushConsumer, error) {
	pushConsumer, err := consumer.NewPushConsumer(
		consumer.WithGroupName(serverAddress),
		consumer.WithNameServer([]string{cfg.RocketMQNameSrv}),
	)
	if err != nil {
		return nil, err
	}

	err = pushConsumer.Subscribe("MessageTopic", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			// TODO: 实现接收消息逻辑
			log.Printf("Received message: %s\n", msg.Body)
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		return nil, err
	}

	err = pushConsumer.Start()
	if err != nil {
		return nil, err
	}

	return pushConsumer, nil
}

func CleanupConsumerQueue(serverAddress string) error {
	// TODO: 此处为 RocketMQ 消费队列清理逻辑
	log.Printf("Consumer queue for %s cleaned up.", serverAddress)
	return nil
}
