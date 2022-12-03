package source

import (
	"fmt"
	"image"
	"regexp"
)

func Wbhandle(URL string, appname string) (*image.Image, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	wb := wbInfo{}
	wb.URL = URL
	wb.App = appname
	switch wb.App {
	case "微博":
		regx1 := regexp.MustCompile(`[^https://m.weibo.cn/status/].*?\?`)
		wbid1 := regx1.FindAllString(wb.URL, -1)
		wbid2 := fmt.Sprint(wbid1)
		wb.ID = fmt.Sprint(wbid2[1 : len(wbid2)-2])
	case "微博轻享版":
		regx1 := regexp.MustCompile(`id=.*?&`)
		wbid1 := regx1.FindAllString(wb.URL, -1)
		wbid2 := fmt.Sprint(wbid1)
		wb.ID = fmt.Sprint(wbid2[4 : len(wbid2)-2])
	}
	rchan := make(chan int)
	dchan := make(chan int)
	go func() {
		if err := wb.GetWeiboReplies(); err != nil {
			fmt.Println("获取评论失败")
		}
		rchan <- 1
	}()
	go func() {
		if err := wb.GetWbDetail(); err != nil {
			fmt.Println("获取详情失败")
			//return nil, err
		}
		dchan <- 1
	}()

	picstring, err := PicInit("./source/material/weiboBackground.png")
	if err != nil {
		fmt.Println("初始化图片失败")
	}

	<-rchan
	<-dchan
	rgba, err := wb.genWbPic(*picstring)
	if err != nil {
		return nil, err
	}

	return &rgba, nil
}
