package v1

import (
	"ginrss/Cro"
	"ginrss/middleware"
	"ginrss/model"
	Redismoon "ginrss/redismoon"
	"ginrss/utils/errmsg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// Login 后台登陆
func Login(c *gin.Context) {
	var formData model.User
	_ = c.ShouldBindJSON(&formData)
	var token string
	var code int

	//加盐哈希
	formData.Password = Cro.SaltCro(formData.Password)
	formData, code = model.CheckLogin(formData.Username, formData.Password)

	if code == errmsg.SUCCSE {
		//更新活跃用户列表
		Redismoon.SetActUser(formData.Username)
		setToken(c, formData)
	}else {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    formData.Username,
			"id":      formData.ID,
			"message": errmsg.GetErrMsg(code),
			"token":   token,
		})
	}

}

// LoginFront 前台登录
func LoginFront(c *gin.Context) {
	var formData model.User
	_ = c.ShouldBindJSON(&formData)
	var token string
	var code int

	formData.Password = Cro.SaltCro(formData.Password)
	formData, code = model.CheckLoginFront(formData.Username, formData.Password)

	if code == errmsg.SUCCSE {
		setToken(c, formData)
		Redismoon.SetActUser(formData.Username)
	}else {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    formData.Username,
			"id":      formData.ID,
			"message": errmsg.GetErrMsg(code),
			"token":   token,
		})
	}

}


// token生成函数
func setToken(c *gin.Context, user model.User) {

	claims := middleware.MyClaims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 100,
			ExpiresAt: time.Now().Unix() + 7200,
			Issuer:    "GinRss",
			Subject: strconv.Itoa(user.Role),
		},
	}

	//本地
	j := middleware.NewJWT()
	token, err := j.CreateToken(claims)

	//RPC jwt
	//token , err := middleware.GrpcTokenGenerate(claims)


	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  errmsg.ERROR,
			"message": errmsg.GetErrMsg(errmsg.ERROR),
			"token":   token,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    user.Username,
		"id":      user.ID,
		"message": errmsg.GetErrMsg(200),
		"token":   token,
	})
	return
}


