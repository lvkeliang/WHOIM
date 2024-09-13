package dao

import (
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	user "github.com/lvkeliang/WHOIM/RPC/kitex_gen/user/userservice"
	"github.com/lvkeliang/WHOIM/config"
	"log"
)

// 初始化 Kitex 用户服务客户端
var UserClient user.Client

func InitRPC() {
	// 初始化 etcd 服务发现
	endpoints := config.LoadConfig().ZookeeperServers
	r, err := etcd.NewEtcdResolver(endpoints)
	if err != nil {
		log.Fatalf("Failed to create etcd resolver: %v", err)
	}

	// 创建 Kitex 客户端并使用 etcd 进行服务发现
	UserClient, err = user.NewClient("WHOIM.UserService", client.WithResolver(r))
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}
}
