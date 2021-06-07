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


func SearchRecord(rssurl string, username string)  (uint, bool){
	var re Record
	db.Select("id").Where("username = ? AND rssurl = ?", username, rssurl).First(&re)
	if re.ID > 0 {
		return  re.ID, true
	}

	return 0, false
}

func SearchRecordId(feedid int, username string)  (uint, bool){
	feed, ok := GetFeedFromId(feedid)
	if !ok {
		return 0, false
	}

	var re Record
	db.Select("id").Where("username = ? AND rssurl = ?", username, feed.Rssurl).First(&re)
	if re.ID > 0 {
		return  re.ID, true
	}

	return 0, false
}



func CreateRecord(data *Record) int {
	//data.Password = ScryptPw(data.Password)
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}


func DeleteRecord(id uint)int {
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

func GetAllFav() []string {
	var ans []Record
	err := db.Table("record").Find(&ans).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}

	temp := map[string]struct{}{}
	var result []string
	for _, item := range ans {
		a := item.Fav
		if _, ok := temp[a]; !ok { //如果字典中找不到元素，ok=false，!ok为true，就往切片中append元素。
			temp[a] = struct{}{}
			result = append(result, a)
		}
	}
	return result

}

func GetFavFeed(username string, Favname string) []MyFeed {
	var ans []Record
	err := db.Table("record").Where("username = ? AND fav = ?", username, Favname).Find(&ans).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}

	var result []MyFeed
	for _, a := range ans{
		var f MyFeed
		rssurl := a.Rssurl
		db.Table("my_feed").Where("rssurl = ?", rssurl).First(&f)
		result = append(result, f)
	}
	return result
}