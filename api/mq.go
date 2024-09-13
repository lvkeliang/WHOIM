package api

import (
	"github.com/lvkeliang/WHOIM/service"
	"log"
)

// StartServiceQueueListener 启动消息队列监听
func StartServiceQueueListener(serviceID string) error {
	// 调用 service 层启动消息监听，基于 serviceID 创建独立队列
	err := service.StartListeningForService(serviceID)
	if err != nil {
		log.Printf("Failed to start message queue listener for service %s: %v", serviceID, err)
		return err
	}

	log.Printf("Message queue listener for service %s started successfully", serviceID)
	return nil
}

//// ShutdownServiceQueue 关闭当前进程的消息队列，并将未消费的消息送回单发消息队列
//func ShutdownServiceQueue(serviceID string) {
//	// 获取未消费的消息（假设可以从 RocketMQ 中拉取未处理的消息）
//	unconsumedMsgs, err := service.GetUnconsumedMessages(serviceID, 500)
//	if err != nil {
//		log.Printf("Failed to fetch unconsumed messages for service %s: %v", serviceID, err)
//	}
//
//	// 将未消费的消息重新发送到单发消息队列
//	for _, msg := range unconsumedMsgs {
//		err := service.SendMessageBackToQueue(msg)
//		if err != nil {
//			log.Printf("Failed to send message back to single message queue: %v", err)
//		}
//	}
//
//	// 注销消息队列
//	err = service.ShutdownServiceQueue(serviceID, 500)
//	if err != nil {
//		log.Printf("Failed to shutdown consumer queue for service %s: %v", serviceID, err)
//	} else {
//		log.Printf("Successfully shutdown consumer queue for service %s", serviceID)
//	}
//}
