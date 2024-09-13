package service

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/lvkeliang/WHOIM/config"
	"github.com/lvkeliang/WHOIM/dao"
	"log"
)

// StartListeningForService 启动指定 serviceID 的消息队列监听
func StartListeningForService(serviceID string) error {
	cfg := config.LoadConfig()
	topic := serviceID // 使用 serviceID 作为队列的 topic 名称
	tag := cfg.RocketMQTag

	// 调用 dao 层的 ListenToConsumerQueue 来启动消息监听
	err := dao.ListenToConsumerQueue(topic, tag, messageHandler)
	if err != nil {
		log.Printf("Failed to listen to consumer queue for service %s: %v", serviceID, err)
		return err
	}

	log.Printf("Started listening to message queue for service %s successfully", serviceID)
	return nil
}

// SendMessageBackToQueue 将未消费的消息发送回单发消息队列
func SendMessageBackToQueue(message *primitive.MessageExt) error {
	cfg := config.LoadConfig()
	// 将消息重新发送到单发消息队列
	rmqMessage := &primitive.Message{
		Topic: cfg.RocketMQTopic, // 发送回单发消息队列
		Body:  message.Body,
	}

	_, err := dao.MQProducer.SendSync(context.Background(), rmqMessage)
	if err != nil {
		log.Printf("Failed to send message back to RocketMQ: %v", err)
		return err
	}

	log.Printf("Message sent back to single message queue successfully")
	return nil
}

// GetUnconsumedMessages 获取未消费的消息
func GetUnconsumedMessages(serviceID string, maxNums int) ([]*primitive.MessageExt, error) {
	// 调用 dao 层拉取未消费的消息
	unconsumedMsgs, err := dao.GetUnconsumedMessages(serviceID, maxNums)
	if err != nil {
		log.Printf("Failed to fetch unconsumed messages for service %s: %v", serviceID, err)
		return nil, err
	}

	log.Printf("Fetched %d unconsumed messages for service %s", len(unconsumedMsgs), serviceID)
	return unconsumedMsgs, nil
}

// ShutdownServiceQueue 关闭当前进程的消息队列，并将未消费的消息送回单发消息队列
func ShutdownServiceQueue(serviceID string, maxNums int) error {
	// 获取未消费的消息
	unconsumedMsgs, err := GetUnconsumedMessages(serviceID, maxNums)
	if err != nil {
		log.Printf("Failed to fetch unconsumed messages for service %s: %v", serviceID, err)
		return err
	}

	// 将未消费的消息重新发送到单发消息队列
	for _, msg := range unconsumedMsgs {
		err := SendMessageBackToQueue(msg)
		if err != nil {
			log.Printf("Failed to send message back to single message queue: %v", err)
		}
	}

	// 注销消息队列
	err = dao.ShutdownConsumerQueue(serviceID)
	if err != nil {
		log.Printf("Failed to shutdown consumer queue for service %s: %v", serviceID, err)
		return err
	} else {
		log.Printf("Successfully shutdown consumer queue for service %s", serviceID)
	}

	return nil
}
