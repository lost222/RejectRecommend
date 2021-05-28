package model

import (
	"ginrss/utils/errmsg"
	"github.com/jinzhu/gorm"
)

type MyFeed struct {
	gorm.Model
	//Recid int `gorm:"type:int;not null " json:"recid"`
	Feedname string `gorm:"type:varchar(256);not null " json:"username"`
	Rssurl string `gorm:"type:varchar(256);not null " json:"rssurl"`
	Fav string `gorm:"type:varchar(256);not null " json:"fav"`
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

func GetUserFeeds(username string, pageSize int, pageNum int) []MyFeed{
	var subrec []Record
	err := db.Limit(pageSize).Offset((pageNum-1)*pageSize).Where("username = ?", username).Find(&subrec).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil
	}

	var feeds []MyFeed

	//db.Model(&MyFeed{}).Select("my_feed.feedname, my_feed.rssurl, my_feed.fav").Joins("left join record on my_feed.rssurl = record.rssurl").Where("record.username = ?", username).Scan(feeds)
	// SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id

	//var feeds []MyFeed
	for _, rec := range subrec{
		var feed MyFeed
		err = db.Table("my_feed").Where("rssurl = ?", rec.Rssurl).First(&feed).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil
		}
		feeds = append(feeds, feed)
	}
	return feeds
}

func GetFeedFromId(feedID int) (MyFeed, bool) {
	var feed MyFeed
	db.Where("ID = ?", feedID).First(&feed)
	if feed.ID > 0 {
		return feed,true
	}
	return feed, false

}