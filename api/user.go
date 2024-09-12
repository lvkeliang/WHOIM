package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 注册路由处理函数
func registerHandler(c *gin.Context) {
	var registerReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 调用 userservice 的 Register 方法
	success, err := userClient.Register(context.Background(), registerReq.Username, registerReq.Password, registerReq.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User registration failed"})
	}
}

// 登录路由处理函数
func loginHandler(c *gin.Context) {
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 解析请求体
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 调用 userservice 的 Login 方法
	token, err := userClient.Login(context.Background(), loginReq.Username, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
