package source

import (
	"errors"
	"fmt"
	"net/http"
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
		fmt.Println(bvid3)
		return bvid3
	}
	return ""
}

// 获取短链重定向之后的链接，方便拿到一些原始信息

func GetBvId(url string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	//if err != nil {
	//	println("\nerror on here\n")
	//	return "", err
	//}
	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	response, _ := client.Do(req)
	//if err != nil {
	//	fmt.Printf("fail to get response on GetBvID\n")
	//	return "", err
	//}
	defer response.Body.Close()
	redirectUrl, _ := response.Location()
	redirectUrl1 := fmt.Sprint(redirectUrl)
	//fmt.Printf(redirectUrl1)
	if bvid := GetBv(redirectUrl1); bvid != "" {
		return bvid, nil
	}

	return GetBvId(redirectUrl1)
}
