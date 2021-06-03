package routes

import (
	v1 "ginrss/api/v1"
	"ginrss/middleware"
	"ginrss/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter()  {
	gin.SetMode(utils.AppMode)
	r := gin.Default()
	r.Use(middleware.Cors())

	//普通用户
	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken())
	{

		auth.GET("feeds",v1.GetUserFeeds)
		auth.GET("favs",v1.GetFavList)

		auth.GET("feed/info/", v1.GetFeedInfo)
		auth.GET("feed/list/", v1.GetFavFeed)

		auth.POST("subscribe/add", v1.AddRecord)
		auth.POST("subscribe/del", v1.DeleteRecord)

	}

	superauth := r.Group("api/v1")
	superauth.Use(middleware.JwtTokenBackend())
	{
		//todo handle函数实现
		superauth.GET("users", v1.GetAllUsers)
		superauth.POST("user/add",v1.AddUser)
		superauth.DELETE("user/:id",v1.DeleteUser)
		superauth.GET("user/active",v1.GetActiveUsers)
		superauth.GET("pushservice",v1.AddUser)
	}

	router := r.Group("api/v1")
	{


		router.POST("login", v1.LoginFront)
		router.POST("backend/login", v1.Login)


	}

	r.Run(utils.HttpPort)
}

