package source

import (
	"bytes"
	"fmt"
	"github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

func BilibiliPic(URL string) (image.Image, error) {
	bvinfo := BvInfo{}
	if _, err := bvinfo.GetBvReplies(URL); err != nil {
		fmt.Println("获取评论失败")
		return nil, err
	}
	picstring, err := PicInit("/Users/wangzehong/Pictures/bilibiliBackground3.png")
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
func appendQr(rgba image.RGBA, picString PicString, URL string) image.Image {
	code, _ := myEncode(URL, qrcode.Medium, codeSize, color.RGBA{255, 255, 255, 255}, color.RGBA{116, 125, 140, 255})
	qrcode, _ := png.Decode(bytes.NewReader(code))
	point := image.Point{(picString.Background.Bounds().Dx()/2 - codeSize/2) * -1, (int(picString.Pt.Y>>6) - codeSize/2 + codeSize/4) * -1}
	draw.Draw(&rgba, rgba.Bounds(), qrcode, point, draw.Src)
	SubImg := rgba.SubImage(image.Rectangle{image.Point{0, 0},
		image.Point{picString.Background.Bounds().Dx(), codeSize + int(picString.Pt.Y>>6)}})
	fmt.Println("point", image.Point{picString.Background.Bounds().Dx(), 2*codeSize + int(picString.Pt.Y>>6)})
	return SubImg
}
