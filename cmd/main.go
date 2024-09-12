package main

import (
	"github.com/lvkeliang/WHOIM/api"
	"github.com/lvkeliang/WHOIM/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	api.InitRPC()
	go api.InitRouter()

	// 初始化 Zookeeper，获取 serverAddress
	serverAddress, err := api.InitZookeeper(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Zookeeper: %v", err)
	}

	// 初始化 RocketMQ 消费者
	consumer, err := api.InitConsumer(serverAddress, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize RocketMQ consumer: %v", err)
	}
	defer func() {
		// 清理 Zookeeper 注册信息
		err := api.CleanupZookeeper(serverAddress)
		if err != nil {
			log.Fatalf("Failed to clean up Zookeeper: %v", err)
		}
		log.Println("Zookeeper registration cleaned up.")
		api.CloseZookeeper()

		if err := consumer.Shutdown(); err != nil {
			log.Fatalf("Failed to shutdown consumer: %v", err)
		}
		log.Println("Consumer shut down gracefully.")
	}()

	// 监听系统信号，优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
