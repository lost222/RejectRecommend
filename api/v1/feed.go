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
	//get feedID from query
	feedID, _ := strconv.Atoi(c.Query("feedid"))
	feeddata,ok := model.GetFeedFromId(feedID)
	if !ok{
		code = errmsg.ERROR_Feed_NOT_EXIST
	}else {
		code = errmsg.SUCCSE
	}
	//一大堆redis操作，假设都miss

	//尝试缓存中是否存在
	var feed *gofeed.Feed
	cc := Redismoon.Cache{
		Rssurl: feeddata.Rssurl,
	}
	ok, err := cc.GetFromRedis()
	if ok {
		feed = &cc.Feed
		fmt.Println("cache hit")
	}else {
		//cache miss
		fmt.Println("cache miss")
		fp := gofeed.NewParser()
		feed, err = rss.FetchURL(fp, feeddata.Rssurl)
		errmsg.CheckErr(err)
		fmt.Println(feed)
		cc.Feed = *feed
		//写入cache
		cc.SaveInRedis()
	}

	//最终将获得的XML解码为json传输

	//todo 修改Feed表中LatesTitle项目

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"feed":  *feed,
			"message": errmsg.GetErrMsg(code),
		},
	)

}

