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

	auth := r.Group("api/v1")
	auth.Use(middleware.JwtToken())
	{
		//// 用户模块的路由接口
		//auth.GET("admin/users", v1.GetUsers)
		//auth.PUT("user/:id", v1.EditUser)
		//auth.DELETE("user/:id", v1.DeleteUser)
		////修改密码
		//auth.PUT("admin/changepw/:id", v1.ChangeUserPassword)
		// 收藏夹模块的路由接口
		//auth.GET("admin/category", v1.GetCate)
		//auth.POST("category/add", v1.AddCategory)
		//auth.PUT("category/:id", v1.EditCate)
		//auth.DELETE("category/:id", v1.DeleteCate)
		//// 订阅模块的路由接口

		auth.GET("feeds",v1.GetUserFeeds)
		auth.GET("favs",v1.GetFavList)
		//auth.GET("feeds/info/:id", v1.GetArtInfo)
		auth.GET("feed/info/", v1.GetFeedInfo)
		auth.GET("feed/list/", v1.GetFavFeed)
		//auth.GET("feeds", v1.GetArt)
		auth.POST("subscribe/add", v1.AddRecord)
		auth.DELETE("subscribe/:id", v1.DeleteRecord)


	}

	router := r.Group("api/v1")
	{
		// 用户信息模块
		router.POST("user/add", v1.AddUser)
		//router.GET("user/:id", v1.GetUserInfo)
		//返回用户订阅的所有Feed名称
		//router.GET("feeds",v1.GetUserFeeds)

		//get single feed info
		router.GET("feed/info", v1.GetFeedInfo)
		router.DELETE("user/:id", v1.DeleteUser)

		//订阅
		//router.POST("subscribe/add", v1.AddRecord)
		//router.DELETE("subscribe/:id", v1.DeleteRecord)

		router.POST("login", v1.LoginFront)

		//:todo 后台管理系统
		//router.GET("users", v1.GetUsers)

	}

	r.Run(utils.HttpPort)
}

