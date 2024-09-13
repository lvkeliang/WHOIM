package dao

import (
	"log"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/lvkeliang/WHOIM/config"
)

var ZkConn *zk.Conn

var serverID string // 保存 serverID

// GetServerID 返回注册到 Zookeeper 的 serverID
func GetServerID() string {
	return serverID
}

// InitZookeeper 初始化 Zookeeper 连接
func InitZookeeper() {
	// 从配置中获取 Zookeeper 地址
	cfg := config.LoadConfig()

	// 连接 Zookeeper
	var err error
	ZkConn, _, err = zk.Connect(cfg.ZookeeperServers, time.Second*10)
	if err != nil {
		log.Fatalf("Failed to connect to Zookeeper: %v", err)
	}
	log.Println("Zookeeper connection established")
}

// RegisterToZookeeper 向 Zookeeper 注册服务节点
func RegisterToZookeeper(serverID, serverAddress string) error {
	// 父节点路径
	parentPath := "/servers"
	// 服务节点路径
	path := parentPath + "/" + serverID
	data := []byte(serverAddress)

	// Zookeeper 节点权限
	acls := zk.WorldACL(zk.PermAll)

	// 检查父节点 /servers 是否存在
	exists, _, err := ZkConn.Exists(parentPath)
	if err != nil {
		log.Printf("Failed to check existence of parent node %s: %v", parentPath, err)
		return err
	}

	// 如果父节点不存在，先创建父节点
	if !exists {
		_, err := ZkConn.Create(parentPath, []byte{}, 0, acls)
		if err != nil {
			log.Printf("Failed to create parent node %s: %v", parentPath, err)
			return err
		}
		log.Printf("Parent node %s created successfully", parentPath)
	}

	// 创建临时节点（临时节点会在断开连接时自动删除）
	_, err = ZkConn.Create(path, data, zk.FlagEphemeral, acls)
	if err != nil {
		log.Printf("Failed to register server to Zookeeper: %v", err)
		return err
	}
	log.Printf("Server registered to Zookeeper with path: %s", path)
	return nil
}
