package service

import (
	"context"
	"encoding/json"
	"github.com/lvkeliang/WHOIM/dao"
	"github.com/lvkeliang/WHOIM/protocol"
	"log"
)

// HandleSingleMessage 处理单发消息，并将消息发送到 RocketMQ 队列
func HandleSingleMessage(msg *protocol.MessageProtocol) {

	// 将 MessageProtocol 结构体转换为 JSON
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	// 发送消息到 RocketMQ 单发消息队列
	err = dao.SendMessage(messageBytes)

	log.Printf("Message sent to RocketMQ successfully, result: %s", err)
}

func HandleSingleMessageDirect(msg *protocol.MessageProtocol) {
	devices, err := dao.UserClient.GetUserDevices(context.Background(), msg.ReceiverID)
	if devices == nil {
		log.Printf("User %s has no connected devices", msg.ReceiverID)
		return
	}

	// 将 MessageProtocol 结构体转换为 JSON
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	for _, status := range devices {
		// 发送消息到 RocketMQ 单发消息队列
		err = dao.SendMessageDirect(messageBytes, status.ServerAddress)
		log.Printf("Message sent to RocketMQ successfully, result: %s", err)
	}

}
