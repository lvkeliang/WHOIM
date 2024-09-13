package service

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/lvkeliang/WHOIM/protocol"
	"log"
)

// messageHandler 处理来自队列的消息
func messageHandler(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		log.Printf("Received message: %s", string(msg.Body))

		// 解码消息体为 MessageProtocol 结构体
		var message protocol.MessageProtocol
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return consumer.ConsumeRetryLater, err
		}

		// 打印解码后的消息
		log.Printf("Decoded message: %+v", message)

		// 业务逻辑处理
		err = ProcessMessage(&message)
		if err != nil {
			log.Printf("Failed to process message: %v", err)
			return consumer.ConsumeRetryLater, err
		}
	}

	// 返回成功，表示消息已成功处理
	return consumer.ConsumeSuccess, nil
}

// ProcessMessage 处理从队列接收到的消息,将消息广播给接收用户的所有连接
func ProcessMessage(message *protocol.MessageProtocol) error {
	// 从消息中提取接收方用户 ID
	receiverID := message.ReceiverID

	// 将消息编码为 JSON 格式
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return err
	}

	// 广播消息给接收方用户的所有设备
	err = SendMessageToUserDevices(receiverID, messageBytes)
	if err != nil {
		log.Printf("Failed to send message to user %s devices: %v", receiverID, err)
		return err
	}

	log.Printf("Message sent to all devices of user %s", receiverID)
	return nil
}
