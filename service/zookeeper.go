package service

import (
	"log"
	"time"

	"github.com/lvkeliang/WHOIM/dao"
)

// GetServerID 返回注册到 Zookeeper 的 serverID
func GetServerID() string {
	return dao.GetServerID()
}

// RegisterProcessToZookeeper 向 Zookeeper 注册当前进程
func RegisterProcessToZookeeper(serverAddress string) error {
	// 生成唯一的 serverID （可以使用时间戳或其他方式确保唯一性）
	serverID := "server_" + time.Now().Format("20060102150405")

	// 调用 dao 层的注册方法
	err := dao.RegisterToZookeeper(serverID, serverAddress)
	if err != nil {
		log.Printf("Failed to register process to Zookeeper: %v", err)
		return err
	}
	log.Printf("Process registered to Zookeeper with serverID: %s", serverID)

	return nil
}
