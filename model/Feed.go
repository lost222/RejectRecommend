package model

import (
	"fmt"
	"ginrss/utils/errmsg"
	"github.com/jinzhu/gorm"
	"strings"
)

type MyFeed struct {
	gorm.Model
	Rssurl string `gorm:"type:varchar(256);not null " json:"rssurl"`
	Feedname string `gorm:"type:varchar(256);not null " json:"feedname"`
	FeedDesc string `gorm:"type:varchar(256);not null " json:"feeddesc"`
	LatesTitle string `gorm:"type:varchar(256);not null " json:"latestitle"`
}

func CheckFeed(url string) (string, bool) {
	var feed MyFeed
	db.Select("id").Where("feedname = ?", url).First(&feed)
	if feed.ID > 0 {
		return feed.Feedname, true
	}
	return "", false
}

func CreateFeed(data *MyFeed) int {
	//data.Password = ScryptPw(data.Password)
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}


func UpdateLastTitle(feedId int, title string) {
	err := db.Model(&MyFeed{}).Where("ID = ?", feedId).Update("lates_title", title).Error
	if err !=nil{
		fmt.Println("UpdateLastTitle", err)
	}
}


func GetUserFeeds(username string, pageSize int, pageNum int) ([]MyFeed, int){
	var subrec []Record

	//先获得数据量
	var count int
	db.Table("record").Where("username = ?", username).Count(&count)
	//分页先查询record
	err := db.Limit(pageSize).Offset((pageNum-1)*pageSize).Where("username = ?", username).Find(&subrec).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil , 0
	}

	var feeds []MyFeed

	//var feeds []MyFeed
	//然后根据record查询Feed
	for _, rec := range subrec{
		var feed MyFeed
		err = db.Table("my_feed").Where("rssurl = ?", rec.Rssurl).First(&feed).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil , 0
		}

		feeds = append(feeds, feed)
	}
	return feeds, count
}


func SearchFeedByname(username string, feedname string) ([]MyFeed, int){
	feeds , _ := GetUserFeeds(username, -1, 0)
	var ans []MyFeed
	for _, feed := range feeds{
		if strings.Contains(feed.Feedname, feedname) {
			ans = append(ans, feed)
		}
	}

	return ans , len(ans)
}

func GetFeedFromId(feedID int) (MyFeed, bool) {
	var feed MyFeed
	db.Where("ID = ?", feedID).First(&feed)
	if feed.ID > 0 {
		return feed,true
	}
	return feed, false
}