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
}


func CreateRecord(data *Record) int {
	//data.Password = ScryptPw(data.Password)
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
	//:todo 钩子函数，每次插入都查一查URL 有没有名字，没有就插入名字
}

func SearchUserSubRecord(username string, pageSize int, pageNum int) []Record{
	var ans []Record

	err := db.Limit(pageSize).Offset((pageNum-1)*pageSize).Where("username = ?", username).Find(&ans).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}

	return ans
}