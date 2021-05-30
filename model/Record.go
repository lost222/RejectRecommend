package model

import (
	"ginrss/utils/errmsg"
	"github.com/jinzhu/gorm"
)

type Record struct {
	gorm.Model
	//Recid int `gorm:"type:int;not null " json:"recid"`
	Username string `gorm:"type:varchar(20);not null " json:"username"`
	Rssurl string `gorm:"type:varchar(256);not null " json:"rssurl"`
	Fav string `gorm:"type:varchar(256);not null " json:"fav"`

}

func CheckRecord(rssurl string, username string)  (code int){
	var user Record
	db.Select("id").Where("username = ? AND rssurl = ?", username, rssurl).First(&user)
	if user.ID > 0 {
		return errmsg.ERROR_RECORD_EXIST //1001
	}
	return errmsg.SUCCSE
}


func CreateRecord(data *Record) int {
	//data.Password = ScryptPw(data.Password)
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}


func DeleteRecord(id int)int {
	var re Record
	err := db.Where("id = ?", id).Delete(&re).Error
	if err != nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}


func SearchUserSubRecord(username string, pageSize int, pageNum int) []Record{
	var ans []Record

	err := db.Limit(pageSize).Offset((pageNum-1)*pageSize).Where("username = ?", username).Find(&ans).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}

	return ans
}