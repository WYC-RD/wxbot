package source

import (
	"fmt"
	"regexp"
)

func GetBv(url string) string {

	regx1 := regexp.MustCompile(`[^https://www.bilibili.com/video/].*?/`)

	bvid1 := regx1.FindAllString(url, -1)
	bvid2 := fmt.Sprint(bvid1)
	//println(bvid2)
	//println(len(bvid2))
	if len(bvid2) > 10 {
		bvid3 := fmt.Sprint(bvid2[1 : len(bvid2)-2])
		return bvid3
	}
	return ""
}
