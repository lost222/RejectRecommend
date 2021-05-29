package model

import (
	"ginrss/utils/errmsg"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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

//delete user

func DeleteUser(usrid int) int {
	var usr User
	err := db.Where("id = ?", usrid).Delete(&usr).Error
	if err != nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

func CheckLogin(username string, password string) (User, int) {
	var user User


	db.Where("username = ?", username).First(&user)

	//todo: 加盐哈希
	//var PasswordErr error
	//PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	//if PasswordErr != nil {
	//	return user, errmsg.ERROR_PASSWORD_WRONG
	//}

	if user.Password != password{
		return user, errmsg.ERROR_PASSWORD_WRONG
	}

	if user.ID == 0 {
		return user, errmsg.ERROR_USER_NOT_EXIST
	}

	if user.Role != 1 {
		return user, errmsg.ERROR_USER_NO_RIGHT
	}
	return user, errmsg.SUCCSE
}

func CheckLoginFront(username string, password string) (User, int) {
	var user User
	var PasswordErr error

	db.Where("username = ?", username).First(&user)

	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if user.ID == 0 {
		return user, errmsg.ERROR_USER_NOT_EXIST
	}
	if PasswordErr != nil {
		return user, errmsg.ERROR_PASSWORD_WRONG
	}
	return user, errmsg.SUCCSE
}