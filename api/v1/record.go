package v1

import (
	"encoding/json"
	"fmt"
	"ginrss/model"
	"ginrss/rss"
	"ginrss/utils/errmsg"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"net/http"
	"strconv"
)

func PullRecordSub(c *gin.Context)  {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	userName := c.Query("username")

	if pageSize == 0{
		pageSize = -1
	}

	if pageNum == 0{
		pageNum = 1
	}

	data := model.SearchUserSubRecord(userName, pageSize, pageNum)

	//一大堆redis操作，假设都miss

	//最终将获得的XML解码为json传输
	fp := gofeed.NewParser()
	var feeds [][]byte

	for _, d := range data{
		//fullRssUrl := utils.Rsshublink + d.Rssurl
		feed, err := rss.FetchURL(fp, d.Rssurl)
		//直接序列化解决深拷贝问题
		buffer, _ := json.Marshal(feed)
		feeds = append(feeds, buffer)
		errmsg.CheckErr(err)
		fmt.Println(feed)
	}

	code = errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"feeds":    feeds,
			"message": errmsg.GetErrMsg(code),
		},
	)

}