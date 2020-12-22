package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lz1998/ecust_im/handler"
)

func main() {
	router := gin.Default()
	router.GET("/ping", handler.Ping) // 测试用

	{
		group := router.Group("/account")
		group.POST("/register", handler.Register) // 注册
		group.POST("/login", handler.Login) // 登陆
		group.POST("/setUserInfo", handler.SetUserInfo) // 设置用户信息
		group.POST("/getUserInfo", handler.GetUserInfo) // 获取用户信息
	}
	{
		group := router.Group("/im")
		group.Use(handler.CheckLogin)
		//group.Use(middleware.CheckLogin)// 登陆检测，从token获取信息，放在context
		group.POST("/createGroup", handler.CreateGroup) // 创建群
		group.POST("/getFriends", handler.GetFriends) // 获取好友列表
		group.POST("/getGroups", handler.GetGroups) // 获取群列表
		group.POST("/processAdd", handler.ProcessAdd) // 处理 加好友/群 请求

		// websocket
		group.GET("/ws", handler.WsHandler)
	}

	if err := router.Run(":9999"); err != nil {
		panic(err)
	}
}
