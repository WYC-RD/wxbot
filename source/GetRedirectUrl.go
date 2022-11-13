package source

import (
	"errors"
	"fmt"
	"net/http"
)

// 获取短链重定向之后的链接，方便拿到一些原始信息

func GetRedirectUrl(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	response, err := client.Do(req)
	redirectUrl, _ := response.Location()
	redirectUrl1 := fmt.Sprint(redirectUrl)
	if bvid := GetBv(redirectUrl1); bvid != "" {
		return bvid
	}

	return GetRedirectUrl(redirectUrl1)
}
