package source

import (
	"fmt"
	"image"
)

func Wbhandle(URL string, appname string) (image.Image, error) {
	wb := wbInfo{}
	wb.URL = URL
	wb.App = appname
	if err := wb.GetWeiboReplies(); err != nil {
		fmt.Println("获取评论失败")
		//return nil, err
	}
	if err := wb.GetWbDetail(); err != nil {
		fmt.Println("获取详情失败")
		//return nil, err
	}
	picstring, err := PicInit("./source/material/weiboBackground.png")
	if err != nil {
		fmt.Println("初始化图片失败")
		//return nil, err
	}
	rgba, err := wb.genWbPic(*picstring)
	if err != nil {
		return nil, err
	}

	//rgba = image.RGBA
	return rgba, nil
}
