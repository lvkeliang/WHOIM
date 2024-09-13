package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lvkeliang/WHOIM/protocol"
	"github.com/lvkeliang/WHOIM/service"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebSocketConnection 处理 WebSocket 连接
func HandleWebSocketConnection(c *gin.Context) {
	// 从请求中获取 token、userID 和 deviceID
	token := c.Query("token")
	userID := c.Query("user_id")
	deviceID := c.Query("device_id")
	if token == "" || userID == "" || deviceID == "" {
		log.Println("Missing token, user_id or device_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token, user_id or device_id"})
		return
	}

	// 验证 token 并获取用户信息
	userInfo, err := service.ValidateToken(token)
	if err != nil || userInfo.Id != userID {
		log.Printf("Invalid token for user: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 获取 serverID
	serverID := service.GetServerID() // 获取 serverID
	onlineSuccess, err := service.SetUserOnline(userID, deviceID, serverID)
	if err != nil || !onlineSuccess {
		log.Printf("Failed to set user %s device %s online: %v", userID, deviceID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set device online"})
		return
	}

	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// 将 WebSocket 连接加入连接池
	service.AddConnection(userID, deviceID, conn)
	defer func() {
		// 设置用户设备离线
		offlineSuccess, err := service.SetUserOffline(userID, deviceID)
		if err != nil || !offlineSuccess {
			log.Printf("Failed to set user %s device %s offline: %v", userID, deviceID, err)
		}
		// 从连接池中移除连接
		service.RemoveConnection(userID, deviceID)
	}()

	// 处理 WebSocket 消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}
		log.Printf("Received message from user %s device %s: %s", userID, deviceID, string(message))

		// 可以在这里处理收到的消息，或者转发给其他设备
		ProcessWebSocketMessage(userID, deviceID, message)
	}
}

func ProcessWebSocketMessage(userID, deviceID string, message []byte) {
	log.Printf("Processing message from user %s device %s: %s", userID, deviceID, message)

	// 解析消息为 MessageProtocol
	var msgProtocol protocol.MessageProtocol
	err := json.Unmarshal(message, &msgProtocol)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	// 判断消息目标类型
	switch msgProtocol.TargetType {
	case protocol.SingleTarget:
		service.HandleSingleMessageDirect(&msgProtocol)
	case protocol.GroupTarget:
		// TODO: 实现群发消息
		log.Printf("Received group message, processing not implemented")
	default:
		log.Printf("Unknown target type: %s", msgProtocol.TargetType)
	}
}
