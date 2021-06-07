package v1

import (
	"ginrss/middleware"
	"ginrss/model"
	Redismoon "ginrss/redismoon"
	"ginrss/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)


var code int
//添加
func AddUser(c *gin.Context) {
	var data model.User
	_ = c.ShouldBindJSON(&data)

	code = model.CheckUser(data.Username)
	if code == errmsg.SUCCSE {
		model.CreateUser(&data)
	}
	if code == errmsg.ERROR_USERNAME_USED {
		code = errmsg.ERROR_USERNAME_USED
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
//查询用户(1)
func GetAllUsers(c *gin.Context)  {
	Users := model.GetAllUsers()
	count := len(Users)
	code = errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"users":  Users,
			"total" : count,
			"message": errmsg.GetErrMsg(code),
		},
	)
}





//:todo 后台
//查询用户列表
//编辑用户
//删除用户
func DeleteUser(c *gin.Context){
	tokenClaim, _ := c.Get("tokenUser")
	tClaim := tokenClaim.(*middleware.MyClaims)
	userName := tClaim.Username
	userID := model.SearchUserId(userName)
	id, _ := strconv.Atoi(c.Param("id"))
	//试图删除自己
	if userID == uint(id) {
		code = errmsg.ERROR
	}else {
		code = model.DeleteUser(id)
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":code,
			"message":errmsg.GetErrMsg(code),
		})
}

func GetActiveUsers(c *gin.Context)  {
	activeUsers := Redismoon.Getactiveusr()
	c.JSON(
		http.StatusOK, gin.H{
			"status":code,
			"activeUsers":activeUsers,
			"message":errmsg.GetErrMsg(code),
		})
}