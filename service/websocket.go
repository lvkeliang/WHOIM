package service

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// connectionPool 存储每个用户的所有设备连接
var connectionPool = struct {
	sync.RWMutex
	connections map[string]map[string]*websocket.Conn // userID -> deviceID -> *websocket.Conn
}{
	connections: make(map[string]map[string]*websocket.Conn),
}

// AddConnection 将新的 WebSocket 连接加入连接池
func AddConnection(userID, deviceID string, conn *websocket.Conn) {
	connectionPool.Lock()
	defer connectionPool.Unlock()

	// 如果用户还没有任何设备连接，初始化设备映射
	if connectionPool.connections[userID] == nil {
		connectionPool.connections[userID] = make(map[string]*websocket.Conn)
	}

	// 添加设备连接
	connectionPool.connections[userID][deviceID] = conn
	log.Printf("User %s device %s connected", userID, deviceID)
}

// RemoveConnection 从连接池中移除 WebSocket 连接
func RemoveConnection(userID, deviceID string) {
	connectionPool.Lock()
	defer connectionPool.Unlock()

	if devices, ok := connectionPool.connections[userID]; ok {
		// 删除特定设备的连接
		delete(devices, deviceID)
		log.Printf("User %s device %s disconnected", userID, deviceID)

		// 如果用户已没有设备连接，删除用户记录
		if len(devices) == 0 {
			delete(connectionPool.connections, userID)
		}
	}
}

// GetUserConnections 获取指定用户的所有设备连接
func GetUserConnections(userID string) map[string]*websocket.Conn {
	connectionPool.RLock()
	defer connectionPool.RUnlock()

	if devices, ok := connectionPool.connections[userID]; ok {
		// 返回用户的所有设备连接
		return devices
	}
	return nil
}

// SendMessageToUserDevices 向指定用户的所有设备发送消息
func SendMessageToUserDevices(userID string, message []byte) error {
	devices := GetUserConnections(userID)
	if devices == nil {
		log.Printf("User %s has no connected devices", userID)
		return nil
	}

	for deviceID, conn := range devices {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Failed to send message to user %s device %s: %v", userID, deviceID, err)
			conn.Close()
			RemoveConnection(userID, deviceID) // 移除断开的连接
		} else {
			log.Printf("Message sent to user %s device %s", userID, deviceID)
		}
	}
	return nil
}

// SendMessageToDevice 向指定用户的特定设备发送消息
func SendMessageToDevice(userID, deviceID string, message []byte) error {
	connectionPool.RLock()
	defer connectionPool.RUnlock()

	if devices, ok := connectionPool.connections[userID]; ok {
		if conn, exists := devices[deviceID]; exists {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Failed to send message to user %s device %s: %v", userID, deviceID, err)
				conn.Close()
				RemoveConnection(userID, deviceID)
				return err
			}
			log.Printf("Message sent to user %s device %s", userID, deviceID)
			return nil
		}
	}
	log.Printf("User %s device %s is not connected", userID, deviceID)
	return nil
}
