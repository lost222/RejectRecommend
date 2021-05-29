package v1

import (
	"fmt"
	"ginrss/model"
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
	userName := c.Query("username")

	if pageSize == 0{
		pageSize = -1
	}

	if pageNum == 0{
		pageNum = -1
	}

	data := model.GetUserFeeds(userName, pageSize, pageNum)

	code = errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"feeds":  data,
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



	//最终将获得的XML解码为json传输
	fp := gofeed.NewParser()
	feed, err := rss.FetchURL(fp, feeddata.Rssurl)
	//直接序列化解决深拷贝问题
	//buffer, _ := json.Marshal(feed)
	//feeds = append(feeds, buffer)
	errmsg.CheckErr(err)
	fmt.Println(feed)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"feed":  *feed,
			"message": errmsg.GetErrMsg(code),
		},
	)

}


