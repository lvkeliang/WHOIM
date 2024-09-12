package main

import (
	"context"
	user "github.com/lvkeliang/WHOIM/RPC/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, username string, password string, email string) (resp bool, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, username string, password string) (resp string, err error) {
	// TODO: Your code here...
	return
}

// ValidateToken implements the UserServiceImpl interface.
func (s *UserServiceImpl) ValidateToken(ctx context.Context, token string) (resp *user.User, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, id string) (resp *user.User, err error) {
	// TODO: Your code here...
	return
}

// SetUserOnline implements the UserServiceImpl interface.
func (s *UserServiceImpl) SetUserOnline(ctx context.Context, id string, deviceID string, serverAddress string) (resp bool, err error) {
	// TODO: Your code here...
	return
}

// SetUserOffline implements the UserServiceImpl interface.
func (s *UserServiceImpl) SetUserOffline(ctx context.Context, id string, deviceID string) (resp bool, err error) {
	// TODO: Your code here...
	return
}

// GetUserDevices implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserDevices(ctx context.Context, id string) (resp map[string]*user.UserStatus, err error) {
	// TODO: Your code here...
	return
}
