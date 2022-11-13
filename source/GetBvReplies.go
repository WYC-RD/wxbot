package source

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetBvReplies(shortUrl string) string {

	Bv := GetRedirectUrl(shortUrl)
	//Bv := GetBv(rawUrl)
	aid := GetAid(Bv)
	firstUrl := "https://api.bilibili.com/x/v2/reply/main?oid="
	lastUrl := "&plat=1&seek_rpid=&type=1"

	url := fmt.Sprint(firstUrl, aid, lastUrl)

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		//return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		//return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		//return
	}
	rez, _ := simplejson.NewJson([]byte(string(body)))
	//根据返回来的json结构，取到数组
	replies1, _ := rez.Get("data").Get("replies").Array()

	var replies []string

	for _, rpy := range replies1 {
		var holeReply []string

		if each, ok := rpy.(map[string]interface{}); ok {

			if each3, ok2 := each["member"].(map[string]interface{}); ok2 {
				uname := fmt.Sprint("【", each3["uname"], "】", "：")
				holeReply = append(holeReply, uname)
				//holeReply =append(holeReply," : ")
			}
			if each2, ok1 := each["content"].(map[string]interface{}); ok1 {
				message := fmt.Sprint(each2["message"])
				holeReply = append(holeReply, message)
				holeReply = append(holeReply, "\n\n")
			}

			replies = append(replies, strings.Join(holeReply, ""))
		}
	}
	return strings.Join(replies, "")
}

func GetAid(Bv string) string {

	furl := "https://api.bv-av.cn/get-bv-av?id="
	url := fmt.Sprint(furl, Bv)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		//return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		//return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		//return
	}
	rez, _ := simplejson.NewJson([]byte(string(body)))
	data, _ := rez.Get("data").Map()
	aid := fmt.Sprint(data["aid"])
	return aid
}
