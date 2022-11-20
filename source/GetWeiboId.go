package source

import (
	"fmt"
	"regexp"
)

// 获取短链重定向之后的链接，方便拿到一些原始信息

func GetWeiboId(url string) string {

	regx1 := regexp.MustCompile(`[^https://m.weibo.cn/status/].*?\?`)
	wbid1 := regx1.FindAllString(url, -1)
	wbid2 := fmt.Sprint(wbid1)
	bvid3 := fmt.Sprint(wbid2[1 : len(wbid2)-2])
	return bvid3

}
func GetLiteWeiboId(url string) string {

	regx1 := regexp.MustCompile(`id=.*?&`)
	wbid1 := regx1.FindAllString(url, -1)
	wbid2 := fmt.Sprint(wbid1)
	println("wb2:", wbid2)
	bvid3 := fmt.Sprint(wbid2[4 : len(wbid2)-2])
	println("wb3", bvid3)
	return bvid3

}

//if wbid := GetWb(url); wbid != "" {
//	return wbid
//}
//req, err := http.NewRequest("GET", url, nil)
//if err != nil {
//	panic(err)
//}
//client := new(http.Client)
//client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
//	return errors.New("Redirect")
//}
//
//response, err := client.Do(req)
//redirectUrl, _ := response.Location()
//redirectUrl1 := fmt.Sprint(redirectUrl)
//println(redirectUrl1)
//if wbid := GetWb(redirectUrl1); wbid != "" {
//	return wbid
//}

//return GetWeiboId(redirectUrl1)
//}
//func GetWb(url string) string {
//	regx1 := regexp.MustCompile(`[^https://m.weibo.cn/status/].*?\?`)
//	wbid1 := regx1.FindAllString(url, -1)
//	wbid2 := fmt.Sprint(wbid1)
//	fmt.Println("wbid2:", wbid2)
//	if len(wbid2) > 5 {
//		bvid3 := fmt.Sprint(wbid2[1 : len(wbid2)-2])
//		return bvid3
//	}
//	return ""
//}
