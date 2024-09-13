package config

import (
	"os"
	"strconv"
)

type Config struct {
	EtcdAddress                []string
	ZookeeperServers           []string
	RocketMQNameSrv            string
	RocketMQBrokerAddress      string
	WebSocketPort              string
	RocketMQProducerGroupName  string
	RocketMQConsumerGroupName  string
	RocketMQSingleMessageTopic string
	RocketMQTag                string
	RocketMQPullMaxNum         int
}

var config Config

func LoadConfig() *Config {
	config = Config{
		EtcdAddress:                []string{"127.0.0.1:2379"},
		ZookeeperServers:           []string{"127.0.0.1:2181"},
		RocketMQNameSrv:            getEnv("ROCKETMQ_NAMESRV", "127.0.0.1:9876"),
		RocketMQBrokerAddress:      getEnv("ROCKETMQ_BROKER", "127.0.0.1:10911"),
		WebSocketPort:              getEnv("WS_PORT", ":8080"),
		RocketMQProducerGroupName:  getEnv("ROCKETMQ_PRODUCER_GROUP", "WHOIMProducerGroup"),
		RocketMQConsumerGroupName:  getEnv("ROCKETMQ_CONSUMER_GROUP", "WHOIMConsumerGroup"),
		RocketMQSingleMessageTopic: getEnv("ROCKETMQ_TOPIC", "WHOIMSingleMessage"),
		RocketMQTag:                getEnv("ROCKETMQ_TAG", "SingleMessage"),
		RocketMQPullMaxNum:         getEnvAsInt("ROCKETMQ_PULL_MAXNUM", 32),
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
