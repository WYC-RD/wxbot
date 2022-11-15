package source

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func GetWeiboReplies(shortUrl string) string {
	weiboId := GetWeiboId(shortUrl)

	firstUrl := "https://m.weibo.cn/comments/hotflow?mid="

	url := firstUrl + weiboId

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
	replies1, _ := rez.Get("data").Get("data").Array()
	var replies []string

	i := len(replies1)
	if i >= 5 {
		i = 5
	}

	for j := 0; j < i; j++ {
		rpy := replies1[j]

		var holeReply []string

		if each, ok := rpy.(map[string]interface{}); ok {

			if each3, ok2 := each["user"].(map[string]interface{}); ok2 {
				uname := fmt.Sprint(each3["screen_name"])
				holeReply = append(holeReply, uname)
				//holeReply =append(holeReply," : ")
			}
			source := fmt.Sprint("[", each["source"], "]", ":")
			holeReply = append(holeReply, source)
			text := fmt.Sprint(each["text"])
			regx1 := regexp.MustCompile(`<.+>.*</.+>`)
			text = regx1.ReplaceAllString(text, "")
			text = fmt.Sprint(text, "\n\n")
			holeReply = append(holeReply, text)
		}
		replies = append(replies, strings.Join(holeReply, ""))
	}
	//for _, rpy := range replies1{
	//var holeReply []string
	//
	//if each, ok := rpy.(map[string]interface{}); ok {
	//
	//	if each3, ok2 := each["user"].(map[string]interface{}); ok2 {
	//		uname := fmt.Sprint("【", each3["screen_name"], "】")
	//		holeReply = append(holeReply, uname)
	//		//holeReply =append(holeReply," : ")
	//	}
	//	text :=fmt.Sprint("[",each["source"],"]","：",each["text"],"\n")
	//	holeReply =append(holeReply,text)
	//	}
	//	replies = append(replies, strings.Join(holeReply, ""))
	//}
	return strings.Join(replies, "")

}
