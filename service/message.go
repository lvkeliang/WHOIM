package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

// 处理 WebSocket 消息事件
func HandleMessageEvent(messageType int, message []byte, conn *websocket.Conn) {
	var event map[string]interface{}
	err := json.Unmarshal(message, &event)
	if err != nil {
		log.Println("Failed to parse event:", err)
		return
	}

	// 根据事件类型处理不同的逻辑
	switch event["type"] {
	case "send_message":
		handleSendMessage(event, conn)
	case "receive_message":
		handleReceiveMessage(event, conn)
	case "update_status":
		handleUpdateStatus(event, conn)
	default:
		log.Println("Unknown event type:", event["type"])
	}
}

func handleSendMessage(event map[string]interface{}, conn *websocket.Conn) {
	// 处理发送消息逻辑
}

func handleReceiveMessage(event map[string]interface{}, conn *websocket.Conn) {
	// 处理接收消息逻辑
}

func handleUpdateStatus(event map[string]interface{}, conn *websocket.Conn) {
	// 处理用户状态更新逻辑
}
