package source

import (
	"fmt"
	"image"
)

func BilibiliPic(URL string) (image.Image, error) {
	bvinfo := BvInfo{}
	if _, err := bvinfo.GetBvReplies(URL); err != nil {
		fmt.Println("获取评论失败")
		return nil, err
	}
	picstring, err := PicInit("./source/material/bilibiliBackground3.png")
	if err != nil {
		fmt.Println("初始化图片失败")
		return nil, err
	}
	rgba, err := bvinfo.genBvPic(*picstring)
	if err != nil {
		return nil, err
	}

	//rgba = image.RGBA
	return rgba, nil
}
