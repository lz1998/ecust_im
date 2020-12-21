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
		group.POST("/register", handler.Register)
		group.POST("/login", handler.Login)
		group.POST("/setUserInfo", handler.SetUserInfo)
		group.POST("/getUserInfo", handler.GetUserInfo)
	}
	{
		group := router.Group("/im")
		//group.Use(middleware.CheckLogin)// 登陆检测，从token获取信息，放在context
		group.POST("/createGroup", handler.CreateGroup)
		group.POST("/getFriends", handler.GetFriends)
		group.POST("/getGroups", handler.GetGroups)

		// websocket
		group.GET("/ws", func(c *gin.Context) {
			handler.WsHandler(c.Writer, c.Request)
		})
	}

	if err := router.Run(":9999"); err != nil {
		panic(err)
	}
}
