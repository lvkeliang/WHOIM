package main

import (
	"fmt"
	"github.com/lvkeliang/WHOIM/api"
	"github.com/lvkeliang/WHOIM/config"
	"github.com/lvkeliang/WHOIM/dao"
	"github.com/lvkeliang/WHOIM/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	dao.InitRPC()

	// 初始化 Zookeeper
	dao.InitZookeeper()

	// 向 Zookeeper 注册当前进程
	serverID := service.GetServerID()
	err := service.RegisterProcessToZookeeper("127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Failed to register process to Zookeeper: %v", err)
	}

	// 初始化 RocketMQ 生产者
	err = dao.InitMQProducer(cfg.RocketMQGroupName, []string{cfg.RocketMQNameSrv})
	if err != nil {
		log.Fatalf("Failed to initialize RocketMQ producer: %v", err)
	}

	// 启动消息队列监听
	go func() {
		err = api.StartServiceQueueListener(serverID)
		if err != nil {
			log.Fatalf("Failed to start message queue listener: %v", err)
		}
	}()

	// 启动 WebSocket 路由
	go api.InitRouter()

	// 捕捉信号并优雅关闭服务
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	sig := <-sigCh
	fmt.Printf("Received signal: %v, shutting down...\n", sig)

	// 关闭消息队列
	api.ShutdownServiceQueue(serverID)
	dao.ShutdownProducer()
}
