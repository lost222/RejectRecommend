package v1

import (
	"fmt"
	"ginrss/middleware"
	"ginrss/model"
	Redismoon "ginrss/redismoon"
	"ginrss/rss"
	"ginrss/utils/errmsg"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"net/http"
	"strconv"
)

func GetUserFeeds(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	tokenClaim, _ := c.Get("tokenUser")
	tClaim := tokenClaim.(*middleware.MyClaims)
	userName := tClaim.Username
	if tClaim.Username != userName{
		code = errmsg.ERROR_USETOKEN_NOT_MATCH
		c.JSON(
			http.StatusOK, gin.H{
				"status":  code,
				"message": errmsg.GetErrMsg(code),
			},
		)
		return
	}

	if pageSize == 0{
		pageSize = -1
	}

	if pageNum == 0{
		pageNum = -1
	}

	data , count := model.GetUserFeeds(userName, pageSize, pageNum)
	//更新活跃用户
	Redismoon.SetActUser(userName)

	code = errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"feeds":  data,
			"total" : count,
			"message": errmsg.GetErrMsg(code),
		},
	)
}




func GetFeedInfo(c *gin.Context)  {
	//以feedid为参数，获得feed的更新。
	//如果更新在缓存中存在，则直接从缓存中获取。如果不在，则访问网络获取
	//由于网络访问是秒级别，访问redis是毫秒级别，所以如果缓存命中，能大大加快速度。
	feedID, _ := strconv.Atoi(c.Query("feedid"))
	//调用数据库查询feed信息
	feeddata,ok := model.GetFeedFromId(feedID)
	if !ok{
		code = errmsg.ERROR_Feed_NOT_EXIST
	}else {
		code = errmsg.SUCCSE
	}
	//尝试缓存中是否存在
	var feed *gofeed.Feed
	cc := Redismoon.Cache{
		Rssurl: feeddata.Rssurl,
	}
	ok, err := cc.GetFromRedis()
	if ok {
		//cache hit,直接获得
		feed = &cc.Feed
		fmt.Println("cache hit")
	}else {
		//cache miss
		fmt.Println("cache miss")
		fp := gofeed.NewParser()
		//通过网络获得FEED
		feed, err = rss.FetchURL(fp, feeddata.Rssurl)
		errmsg.CheckErr(err)
		//fmt.Println(feed)
		cc.Feed = *feed
		//写入cache
		cc.SaveInRedis()
	}
	//最终将获得的XML解码为json传输
	//修改Feed表中LatesTitle项目
	model.UpdateLastTitle(feedID, feed.Items[0].Title)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"feed":  *feed,
			"message": errmsg.GetErrMsg(code),
		},
	)

}

