package api

import (
	"fmt"
	"github.com/lvkeliang/WHOIM/config"
	"log"
	"time"

	"github.com/go-zookeeper/zk"
)

var zkConn *zk.Conn

// InitZookeeper 初始化Zookeeper连接
func InitZookeeper(cfg config.Config) (string, error) {
	var err error
	zkConn, _, err = zk.Connect(cfg.ZookeeperServers, time.Second) // 使用配置的Zookeeper地址
	if err != nil {
		return "", err
	}

	serverAddress := fmt.Sprintf("/servers/server_%d", time.Now().UnixNano())
	_, err = zkConn.Create(serverAddress, []byte("online"), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return "", err
	}

	log.Printf("Registered with Zookeeper as %s", serverAddress)
	return serverAddress, nil
}

// CleanupZookeeper 清理宕机或停止的服务
func CleanupZookeeper(serverAddress string) error {
	if zkConn == nil {
		return fmt.Errorf("Zookeeper connection is not initialized")
	}

	znodePath := serverAddress

	// 检查该节点是否存在
	exists, _, err := zkConn.Exists(znodePath)
	if err != nil {
		return fmt.Errorf("Failed to check node existence in Zookeeper: %v", err)
	}

	// 删除该节点
	if exists {
		err = zkConn.Delete(znodePath, -1)
		if err != nil {
			return fmt.Errorf("Failed to delete node in Zookeeper: %v", err)
		}
		log.Printf("Node %s successfully removed from Zookeeper", znodePath)
	} else {
		log.Printf("Node %s does not exist in Zookeeper", znodePath)
	}

	// RocketMQ 清理逻辑
	err = CleanupConsumerQueue(serverAddress)
	if err != nil {
		return fmt.Errorf("Failed to clean up consumer queue: %v", err)
	}

	log.Printf("Cleanup for serverAddress %s completed", serverAddress)
	return nil
}

// CloseZookeeper 关闭Zookeeper连接
func CloseZookeeper() {
	if zkConn != nil {
		zkConn.Close()
	}
}
