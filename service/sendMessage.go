package service

import (
	"context"
	"encoding/json"
	"github.com/lvkeliang/WHOIM/config"
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

	cfg := config.LoadConfig()
	// 用一个 map 来记录已发送的 ServerAddress，防止重复发送
	sentServers := make(map[string]bool)

	// 遍历所有设备状态
	for _, status := range devices {
		// 如果该 ServerAddress 尚未发送
		if _, exists := sentServers[status.ServerAddress]; !exists {
			// 发送消息到 RocketMQ 单发消息队列
			err := dao.SendMessageDirect(messageBytes, cfg.RocketMQConsumerGroupName+"_"+status.ServerAddress)
			if err != nil {
				log.Printf("Failed to send message to %s: %v", status.ServerAddress, err)
			} else {
				log.Printf("Message sent to RocketMQ successfully for ServerAddress: %s", status.ServerAddress)
			}
			// 标记这个 ServerAddress 已发送
			sentServers[status.ServerAddress] = true
		}
	}
}
