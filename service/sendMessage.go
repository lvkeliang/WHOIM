package service

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/lvkeliang/WHOIM/config"
	"github.com/lvkeliang/WHOIM/dao"
	"github.com/lvkeliang/WHOIM/protocol"
	"log"
)

// HandleSingleMessage 处理单发消息，并将消息发送到 RocketMQ 队列
func HandleSingleMessage(msg *protocol.MessageProtocol) {
	cfg := config.LoadConfig()

	// 将 MessageProtocol 结构体转换为 JSON
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	// 创建 RocketMQ 消息
	rmqMessage := &primitive.Message{
		Topic: cfg.RocketMQTopic, // 设置单发消息的队列名称
		Body:  messageBytes,
	}

	// 发送消息到 RocketMQ 单发消息队列
	result, err := dao.MQProducer.SendSync(context.Background(), rmqMessage)
	if err != nil {
		log.Printf("Failed to send message to RocketMQ: %v", err)
		return
	}

	log.Printf("Message sent to RocketMQ successfully, result: %s", result.String())
}
