package v1

import (
	"encoding/json"
	"fmt"
	"ginrss/middleware"
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


func AddRecord(c *gin.Context){
	var data model.Record
	_ = c.ShouldBindJSON(&data)

	//var msg string
	//var validCode int
	//msg, validCode = validator.Validate(&data)
	//if validCode != errmsg.SUCCSE {
	//	c.JSON(
	//		http.StatusOK, gin.H{
	//			"status":  validCode,
	//			"message": msg,
	//		},
	//	)
	//	c.Abort()
	//	return
	//}
	tokenClaim, _ := c.Get("tokenUser")
	tClaim := tokenClaim.(*middleware.MyClaims)
	userName := tClaim.Username
	data.Username = userName

	fmt.Println(data.Username)
	fmt.Println(data.Rssurl)

	code = model.CheckRecord(data.Rssurl, data.Username)
	if code == errmsg.SUCCSE {
		model.CreateRecord(&data)
	}
	if code == errmsg.ERROR_RECORD_EXIST {
		code = errmsg.ERROR_RECORD_EXIST
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}


func DeleteRecord(c *gin.Context){
	tokenClaim, _ := c.Get("tokenUser")
	tClaim := tokenClaim.(*middleware.MyClaims)
	userName := tClaim.Username

	//name := context.PostForm("name")
	var data model.Record
	_ = c.ShouldBindJSON(&data)
	fmt.Println(data.Username)
	fmt.Println(data.Rssurl)


	//查询该用户的该订阅
	recordID, ok := model.SearchRecord(data.Rssurl, userName)
	if !ok {
		//记录不存在，直接返回删除成功
		code = errmsg.SUCCSE
	}else {
		code = model.DeleteRecord(recordID)
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":code,
			"message":errmsg.GetErrMsg(code),
		})
}

func DeleteRecordById(c *gin.Context)  {
	feedid, _ := strconv.Atoi(c.Param("id"))
	tokenClaim, _ := c.Get("tokenUser")
	tClaim := tokenClaim.(*middleware.MyClaims)
	userName := tClaim.Username

	//查询用户订阅
	recordId , ok := model.SearchRecordId(feedid, userName)
	if !ok {
		//记录不存在，直接返回删除成功
		code = errmsg.SUCCSE
	}else {
		code = model.DeleteRecord(recordId)
	}

	c.JSON(
		http.StatusOK, gin.H{
			"status":code,
			"message":errmsg.GetErrMsg(code),
		})
}



func GetFavList(c *gin.Context) {
	favList := model.GetAllFav()
	c.JSON(
		http.StatusOK, gin.H{
			"status":errmsg.SUCCSE,
			"favs" : favList,
			"message":errmsg.GetErrMsg(code),
		})
}

func GetFavFeed(c *gin.Context) {
	tokenClaim, _ := c.Get("tokenUser")
	tClaim := tokenClaim.(*middleware.MyClaims)
	userName := tClaim.Username

	fav := c.Query("fav")

	feedList := model.GetFavFeed(userName, fav)
	c.JSON(
		http.StatusOK, gin.H{
			"status":errmsg.SUCCSE,
			"data" : feedList,
			"message":errmsg.GetErrMsg(code),
		})

}