package service

import (
	"context"
	"log"

	user "github.com/lvkeliang/WHOIM/RPC/kitex_gen/user"
	"github.com/lvkeliang/WHOIM/dao"
)

// Register 调用用户注册服务
func Register(username, password, email string) (bool, error) {

	resp, err := dao.UserClient.Register(context.Background(), username, password, email)
	if err != nil {
		log.Printf("Register error: %v", err)
		return false, err
	}
	return resp, nil
}

// Login 调用用户登录服务
func Login(username, password string) (string, error) {

	resp, err := dao.UserClient.Login(context.Background(), username, password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return "", err
	}
	return resp, nil
}

// ValidateToken 调用令牌验证服务
func ValidateToken(token string) (*user.User, error) {

	resp, err := dao.UserClient.ValidateToken(context.Background(), token)
	if err != nil {
		log.Printf("ValidateToken error: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetUserInfo 调用获取用户信息服务
func GetUserInfo(userID string) (*user.User, error) {

	resp, err := dao.UserClient.GetUserInfo(context.Background(), userID)
	if err != nil {
		log.Printf("GetUserInfo error: %v", err)
		return nil, err
	}
	return resp, nil
}

// SetUserOnline 设置用户设备在线，使用 serverID
func SetUserOnline(userID, deviceID, serverID string) (bool, error) {
	resp, err := dao.UserClient.SetUserOnline(context.Background(), userID, deviceID, serverID)
	if err != nil {
		log.Printf("SetUserOnline error: %v", err)
		return false, err
	}
	return resp, nil
}

// SetUserOffline 设置用户设备离线
func SetUserOffline(userID, deviceID string) (bool, error) {

	resp, err := dao.UserClient.SetUserOffline(context.Background(), userID, deviceID)
	if err != nil {
		log.Printf("SetUserOffline error: %v", err)
		return false, err
	}
	return resp, nil
}

// GetUserDevices 获取用户所有在线设备
func GetUserDevices(userID string) (map[string]*user.UserStatus, error) {

	resp, err := dao.UserClient.GetUserDevices(context.Background(), userID)
	if err != nil {
		log.Printf("GetUserDevices error: %v", err)
		return nil, err
	}
	return resp, nil
}
