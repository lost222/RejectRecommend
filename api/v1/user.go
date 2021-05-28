package v1

import (
	"ginrss/model"
	"ginrss/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
)


var code int
//添加
func AddUser(c *gin.Context) {
	var data model.User
	_ = c.ShouldBindJSON(&data)

	//var msg string
	//var validCode int
	//msg, validCode = validator.Validate(&data)
	//if validCode != errmsg.SUCCSE {
	//	c.JSON(
	//		http.StatusOK, gin.H{
	//			"status":  validCode,
	//			"message": msg,
	//		},
	//	)
	//	c.Abort()
	//	return
	//}

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


//:todo 后台
//查询用户列表
//编辑用户
//删除用户