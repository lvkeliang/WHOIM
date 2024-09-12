package config

import (
	"os"
)

type Config struct {
	ZookeeperServers []string
	RocketMQNameSrv  string
}

func LoadConfig() Config {
	return Config{
		ZookeeperServers: []string{"127.0.0.1:2181"}, // 可以从环境变量或配置文件中读取
		RocketMQNameSrv:  os.Getenv("ROCKETMQ_NAMESRV"),
	}
}
