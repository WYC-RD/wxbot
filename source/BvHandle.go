package source

import (
	"fmt"
	"image"
)

func BilibiliPic(URL string) (*image.Image, error) {
	bchan := make(chan *BvInfo)
	go func() {
		bvinfo := BvInfo{}
		bv, err := bvinfo.GetBvReplies(URL)
		if err != nil {
			fmt.Println("获取评论失败")
			//return nil, err
		}
		bchan <- bv
	}()
	pic := make(chan *image.Image)
	go func() {
		picstring, err := PicInit("./source/material/bilibiliBackground3.png")
		if err != nil {
			fmt.Println("初始化图片失败")
			//return nil
		}
		bvinfo := <-bchan
		rgba, err := bvinfo.genBvPic(*picstring)
		if err != nil {
			//return nil, err
		}
		pic <- &rgba
	}()
	rgba := <-pic
	//rgba = image.RGBA
	return rgba, nil
}
