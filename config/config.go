package config

import (
	"os"
	"strconv"
)

type Config struct {
	ZookeeperServers   []string
	RocketMQNameSrv    string
	WebSocketPort      string
	RocketMQGroupName  string
	RocketMQTopic      string
	RocketMQTag        string
	RocketMQPullMaxNum int
}

var config Config

func LoadConfig() *Config {
	config = Config{
		ZookeeperServers:   []string{"127.0.0.1:2181"},
		RocketMQNameSrv:    getEnv("ROCKETMQ_NAMESRV", "127.0.0.1:9876"),
		WebSocketPort:      getEnv("WS_PORT", ":8080"),
		RocketMQGroupName:  getEnv("ROCKETMQ_GROUP", "consumerGroup"),
		RocketMQTopic:      getEnv("ROCKETMQ_TOPIC", "SingleMessageQueue"),
		RocketMQTag:        getEnv("ROCKETMQ_TAG", "messageTag"),
		RocketMQPullMaxNum: getEnvAsInt("ROCKETMQ_PULL_MAXNUM", 32),
	}

	return &config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(name string, fallback int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
