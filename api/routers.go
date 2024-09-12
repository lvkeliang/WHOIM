package api

import "github.com/gin-gonic/gin"

func InitRouter() {
	router := gin.Default()

	// 定义 whoim 路由组
	whoimGroup := router.Group("/whoim")
	{
		// 用户路由组
		userGroup := whoimGroup.Group("/user")
		{
			userGroup.POST("/register", registerHandler)
			userGroup.POST("/login", loginHandler)
		}

		// WebSocket 路由组
		wsGroup := whoimGroup.Group("/ws")
		{
			wsGroup.GET("/connect", websocketHandler)
		}
	}

	// 启动服务器
	router.Run(":8080")

}
