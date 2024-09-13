package service

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/lvkeliang/WHOIM/dao"
	"log"
	"sync"
)

var (
	// 使用全局 map 来存储每个 serviceID 对应的消费者
	globalConsumers = make(map[string]rocketmq.PushConsumer)
	mu              sync.Mutex // 使用 mutex 来保护并发访问
)

// StartListeningForService 启动指定 serviceID 的消息队列监听
func StartListeningForService(serviceID string) error {
	// 使用 serviceID 作为队列的 topic 名称
	topic := serviceID

	// 调用 InitMQConsumer 来启动消息监听，传入 messageHandler 作为处理函数
	consumer, err := dao.InitMQConsumer(topic, messageHandler)
	if err != nil {
		log.Printf("Failed to listen to consumer queue for service %s: %v", serviceID, err)
		return err
	}

	// 将 consumer 存入全局 map 中，使用 mutex 保护并发访问
	mu.Lock()
	globalConsumers[serviceID] = consumer
	mu.Unlock()

	log.Printf("Started listening to message queue for service %s successfully", serviceID)
	return nil
}

// StopListeningForService 停止指定 serviceID 的消息队列监听
func StopListeningForService(serviceID string) {
	mu.Lock()
	defer mu.Unlock()

	if consumer, exists := globalConsumers[serviceID]; exists {
		// 调用 ShutdownConsumer 来关闭消费者
		dao.ShutdownConsumer(consumer)

		// 从 map 中移除该 serviceID 的 consumer
		delete(globalConsumers, serviceID)
		log.Printf("Stopped listening to message queue for service %s", serviceID)
	} else {
		log.Printf("No consumer found for service %s", serviceID)
	}
}

//// SendMessageBackToQueue 将未消费的消息发送回单发消息队列
//func SendMessageBackToQueue(message *primitive.MessageExt) error {
//	cfg := config.LoadConfig()
//	// 将消息重新发送到单发消息队列
//	rmqMessage := &primitive.Message{
//		Topic: cfg.RocketMQTopic, // 发送回单发消息队列
//		Body:  message.Body,
//	}
//
//	_, err := dao.MQProducer.SendSync(context.Background(), rmqMessage)
//	if err != nil {
//		log.Printf("Failed to send message back to RocketMQ: %v", err)
//		return err
//	}
//
//	log.Printf("Message sent back to single message queue successfully")
//	return nil
//}
//
//// GetUnconsumedMessages 获取未消费的消息
//func GetUnconsumedMessages(serviceID string, maxNums int) ([]*primitive.MessageExt, error) {
//	// 调用 dao 层拉取未消费的消息
//	unconsumedMsgs, err := dao.GetUnconsumedMessages(serviceID, maxNums)
//	if err != nil {
//		log.Printf("Failed to fetch unconsumed messages for service %s: %v", serviceID, err)
//		return nil, err
//	}
//
//	log.Printf("Fetched %d unconsumed messages for service %s", len(unconsumedMsgs), serviceID)
//	return unconsumedMsgs, nil
//}
//
//// ShutdownServiceQueue 关闭当前进程的消息队列，并将未消费的消息送回单发消息队列
//func ShutdownServiceQueue(serviceID string, maxNums int) error {
//	// 获取未消费的消息
//	unconsumedMsgs, err := GetUnconsumedMessages(serviceID, maxNums)
//	if err != nil {
//		log.Printf("Failed to fetch unconsumed messages for service %s: %v", serviceID, err)
//		return err
//	}
//
//	// 将未消费的消息重新发送到单发消息队列
//	for _, msg := range unconsumedMsgs {
//		err := SendMessageBackToQueue(msg)
//		if err != nil {
//			log.Printf("Failed to send message back to single message queue: %v", err)
//		}
//	}
//
//	// 注销消息队列
//	err = dao.ShutdownConsumerQueue(serviceID)
//	if err != nil {
//		log.Printf("Failed to shutdown consumer queue for service %s: %v", serviceID, err)
//		return err
//	} else {
//		log.Printf("Successfully shutdown consumer queue for service %s", serviceID)
//	}
//
//	return nil
//}
