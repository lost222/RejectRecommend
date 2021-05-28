package routes

import (
	v1 "ginrss/api/v1"
	"ginrss/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter()  {
	gin.SetMode(utils.AppMode)
	r := gin.Default()

	router := r.Group("api/v1")
	{
		// 用户信息模块
		router.POST("user/add", v1.AddUser)
		//router.GET("user/:id", v1.GetUserInfo)
		//返回用户订阅的所有Feed名称
		router.GET("feeds",v1.GetUserFeeds)

		//get single feed info
		router.GET("feed/info", v1.GetFeedInfo)

		//:todo 后台管理系统
		//router.GET("users", v1.GetUsers)

	}

	r.Run(utils.HttpPort)
}