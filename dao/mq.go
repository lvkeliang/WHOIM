package dao

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/lvkeliang/WHOIM/config"
	"log"
)

var (
	// Global variables for producer and consumer
	MQProducer rocketmq.Producer
	MQConsumer rocketmq.PushConsumer
)

// InitMQProducer 初始化 RocketMQ 生产者
func InitMQProducer(groupName string, namesrvAddr []string) error {
	var err error
	MQProducer, err = rocketmq.NewProducer(
		producer.WithGroupName(groupName),
		producer.WithNameServer(namesrvAddr),
	)

	if err != nil {
		log.Fatalf("Create RocketMQ producer error: %s", err)
		return err
	}

	err = MQProducer.Start()
	if err != nil {
		log.Fatalf("Start RocketMQ producer error: %s", err)
		return err
	}

	log.Println("RocketMQ producer started successfully")
	return nil
}

type MessageExt = func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)

// InitConsumer initializes the RocketMQ push consumer
func InitConsumer(groupName string, namesrvAddr []string, topic string, tag string, f MessageExt) error {
	var err error
	MQConsumer, err = rocketmq.NewPushConsumer(
		consumer.WithGroupName(groupName),
		consumer.WithNameServer(namesrvAddr),
	)

	if err != nil {
		log.Fatalf("Create consumer error: %s", err.Error())
		return err
	}

	err = MQConsumer.Subscribe(topic, consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: tag,
	}, f)

	if err != nil {
		log.Fatalf("Subscribe topic error: %s", err.Error())
		return err
	}

	err = MQConsumer.Start()
	if err != nil {
		log.Fatalf("Start consumer error: %s", err.Error())
		return err
	}

	fmt.Println("RocketMQ Consumer started successfully")
	return nil
}

// ShutdownProducer shuts down the producer
func ShutdownProducer() {
	if MQProducer != nil {
		err := MQProducer.Shutdown()
		if err != nil {
			log.Printf("Shutdown producer error: %s", err.Error())
		} else {
			fmt.Println("RocketMQ Producer shutdown successfully")
		}
	}
}

// ShutdownConsumer shuts down the consumer
func ShutdownConsumer() {
	if MQConsumer != nil {
		err := MQConsumer.Shutdown()
		if err != nil {
			log.Printf("Shutdown consumer error: %s", err.Error())
		} else {
			fmt.Println("RocketMQ Consumer shutdown successfully")
		}
	}
}

// SendMessageToQueue 发送消息到指定的消息队列
func SendMessageToQueue(topic string, tag string, body string) error {
	// 构建消息
	msg := &primitive.Message{
		Topic: topic,
		Body:  []byte(body),
	}
	msg.WithTag(tag)

	// 发送消息
	res, err := MQProducer.SendSync(context.Background(), msg)
	if err != nil {
		log.Printf("Send message error: %v", err)
		return err
	}

	log.Printf("Send message success: result=%v, messageID=%s", res.Status, res.MsgID)
	return nil
}

// ListenToConsumerQueue 监听消费队列并处理消息
func ListenToConsumerQueue(topic string, tag string, handler MessageExt) error {
	// 初始化消费者（如果没有初始化的话）
	if MQConsumer == nil {
		log.Println("Consumer not initialized, initializing...")
		err := InitConsumer("consumerGroup", []string{"127.0.0.1:9876"}, topic, tag, handler)
		if err != nil {
			return err
		}
	}

	// 订阅并监听消息
	log.Printf("Listening to topic: %s with tag: %s", topic, tag)
	return nil
}

// GetUnconsumedMessages 拉取未消费的消息
func GetUnconsumedMessages(serviceID string, maxNums int) ([]*primitive.MessageExt, error) {
	cfg := config.LoadConfig()

	var unconsumedMsgs []*primitive.MessageExt

	// 创建 PullConsumer 用于拉取未消费的消息
	pullConsumer, err := rocketmq.NewPullConsumer(
		consumer.WithGroupName(cfg.RocketMQGroupName),
		consumer.WithNameServer([]string{cfg.RocketMQNameSrv}),
	)
	if err != nil {
		log.Printf("Failed to create pull consumer: %v", err)
		return nil, err
	}

	// 启动 PullConsumer
	err = pullConsumer.Start()
	if err != nil {
		log.Printf("Failed to start pull consumer: %v", err)
		return nil, err
	}
	defer pullConsumer.Shutdown()

	// 使用 serviceID 作为 topic 进行拉取
	topic := serviceID

	// 拉取指定数量的消息
	resp, err := pullConsumer.Pull(context.Background(), maxNums)
	if err != nil {
		log.Printf("Pull consumer failed: %v", err)
		return nil, err
	}

	// 检查拉取结果
	if len(resp.GetMessageExts()) > 0 {
		unconsumedMsgs = append(unconsumedMsgs, resp.GetMessageExts()...)
	}

	log.Printf("Fetched %d unconsumed messages from topic %s", len(unconsumedMsgs), topic)
	return unconsumedMsgs, nil
}

// ShutdownConsumerQueue 关闭指定 serviceID 的消费者队列
func ShutdownConsumerQueue(serviceID string) error {
	// 关闭消费者
	log.Printf("Shutting down consumer for service %s", serviceID)

	// 调用 RocketMQ 的 API 关闭消费者
	if MQConsumer != nil {
		err := MQConsumer.Shutdown()
		if err != nil {
			log.Printf("Failed to shutdown consumer for service %s: %v", serviceID, err)
			return err
		}
	}

	log.Printf("Consumer queue for service %s shutdown successfully", serviceID)
	return nil
}
