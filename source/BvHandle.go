package source

import (
	"fmt"
	"image"
)

func BilibiliPic(URL string) (*image.Image, error) {
	bchan := make(chan int)
	pic := make(chan *PicString)
	bvinfo := BvInfo{}
	go func() {
		err := bvinfo.GetBvReplies(URL)
		if err != nil {
			fmt.Println("获取评论失败")
			return
		}
		bchan <- 1
	}()

	go func() {
		picstring, err := PicInit("./source/material/bilibiliBackground3.png")
		if err != nil {
			fmt.Println("初始化图片失败")
			return
		}
		pic <- picstring
	}()
	picstring := <-pic
	<-bchan
	rgba, err := bvinfo.genBvPic(*picstring)
	if err != nil {
		fmt.Println("生成图片失败")
	}
	return &rgba, nil
}
