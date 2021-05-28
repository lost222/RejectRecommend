package model

import (
	"ginrss/utils/errmsg"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	//Userid int `gorm:"type:int;not null " json:"userid"`
	Username string `gorm:"type:varchar(20);not null " json:"username"`
	Usermail string `gorm:"type:varchar(256);not null " json:"usermail"`
	Password string `gorm:"type:varchar(256);not null " json:"password"`
	Role int `gorm:"type:int;not null " json:"role"`
}

// CheckUser 查询用户是否存在
func CheckUser(name string) (code int) {
	var user User
	db.Select("id").Where("username = ?", name).First(&user)
	if user.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //1001
	}
	return errmsg.SUCCSE
}

//add user
func CreateUser(data *User) int {
	//data.Password = ScryptPw(data.Password)
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}

