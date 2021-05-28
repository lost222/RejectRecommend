package model

import (
	"bufio"
	"fmt"
	"ginrss/rss"
	"ginrss/utils"
	"ginrss/utils/errmsg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/mmcdole/gofeed"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var db *gorm.DB

var err error

func InitDB()  {
	dmm := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	utils.DbUser,
	utils.DbPassWord,
	utils.DbHost,
	utils.DbPort,
	utils.DbName,
	)
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(utils.Db, dmm)

	if err != nil{
		fmt.Println("connect to mysql wrong", err)
	}

	db.SingularTable(true)
	err := db.AutoMigrate(&User{}, &Record{}, &MyFeed{}).Error

	fmt.Println(err)

	db.DB().SetConnMaxLifetime(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(10 * time.Second)

	//initTestData()
}

func readData2url(filePath string) []string{
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return nil
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	line, err := buf.ReadString('\n')
	line = strings.TrimSpace(line)

	rssurl0 := utils.Rsshublink + line

	ans := make([]string, 0)

	for {
		line, err = buf.ReadString('\n')
		line = strings.TrimSpace(line)
		rssurl := rssurl0 + line
		//fmt.Println(rssurl)
		ans = append(ans, rssurl)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return nil
			}
		}
	}

	return ans
}

func initTestData() {
	var rssUrls []string

	dataPaths := []string{
		"./config/bilibili",
		"./config/weibo",
		"./config/zhihu",
	}

	for _, p := range dataPaths{
		ans := readData2url(p)
		rssUrls = append(rssUrls, ans...)
	}

	fp := gofeed.NewParser()


	for i:=0; i<10;i++{
		userstr := strconv.Itoa(i)
		newmember := &User{
			Username: "test"+userstr,
			Usermail: userstr+"@mail.qq.com",
			Password: "pass"+userstr,
			Role: 0,
		}
		CreateUser(newmember)
	}

	M := len(rssUrls)

	//insert feed

	favs := []string{
		"bilibili",
		"weibo",
		"zhihu",
	}
	for f:=0;f<M;f++{
		subindex := f / 10
		feed, err := rss.FetchURL(fp,rssUrls[f])
		if err !=nil{
			fmt.Println(rssUrls[f])
		}
		errmsg.CheckErr(err)
		rssfeedname := feed.Title

		newFeed := &MyFeed{
			Feedname: rssfeedname,
			Rssurl: rssUrls[f],
			Fav: favs[subindex],
		}
		CreateFeed(newFeed)

	}


	//每个user订阅20条
	for i:=0;i<10;i++{
		//sub := rand.Intn(M)
		//插入订阅 useri 订阅了 rssUrls[sub]
		userstr := strconv.Itoa(i)

		subset := map[int]struct{}{}
		var empty struct{}
		for j:=0; j<20;j++{
			sub := rand.Intn(M)
			subset[sub] = empty
		}

		for k, _ := range subset{
			newRec := &Record{
				Username: "test" + userstr,
				Rssurl: rssUrls[k],
			}

			CreateRecord(newRec)
		}

	}

}