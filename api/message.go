package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lvkeliang/WHOIM/service"
	"log"
	"net/http"
)

// WebSocket Upgrader 配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有跨域连接
		return true
	},
}

// WebSocket 连接路由处理函数
func websocketHandler(c *gin.Context) {
	// 从请求头中获取 Token
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	// 验证 Token
	userInfo, err := userClient.ValidateToken(c, token)
	if err != nil {
		log.Println("Failed to validate token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// 提取用户 ID 和设备 ID（假设从请求中提取设备信息）
	userID := userInfo.Id
	deviceID := c.Query("device_id")
	if deviceID == "" {
		deviceID = "default_device"
	}

	// 将 HTTP 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to websocket:", err)
		return
	}
	defer conn.Close()

	// 使用微服务将用户设备状态设置为在线
	// TODO: 使用ETCD集群注册服务器, 使用服务器在ETCD内的ID, 而不是直接使用服务器地址
	serverAddress := "current-server-address" // 假设这是当前服务器的地址
	isSuccess, err := userClient.SetUserOnline(c, userID, deviceID, serverAddress)
	if err != nil || !isSuccess {
		log.Println("Failed to set user device online:", err)
		return
	}
	log.Printf("User %s device %s connected", userID, deviceID)

	// 处理 WebSocket 消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// 处理接收到的消息
		service.HandleMessageEvent(messageType, message, conn)
	}

	// 在连接关闭时将设备状态设置为下线
	isSuccess, err = userClient.SetUserOffline(c, userID, deviceID)
	if err != nil || !isSuccess {
		log.Println("Failed to set user device offline:", err)
	}
	log.Printf("User %s device %s disconnected", userID, deviceID)
}
