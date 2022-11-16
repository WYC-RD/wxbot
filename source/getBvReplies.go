package source

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 定义评论的结构体
type GetReplies struct {
	Data struct {
		Replies []replies
	}
}
type replies struct {
	Member  Member    `json:"member"`
	Content Content   `json:"content"`
	Replies []replies `json:"replies"`
}
type Member struct {
	Uname string `json:"uname"`
	Sex   string `json:"sex"`
}
type Content struct {
	Message string `json:"message"`
}

func GetBvReplies(URL string) (string, error) {

	Bv, err := GetBvId(URL)
	if err != nil {
		println("GetBvID fail")
		return "", err
	}
	fmt.Println("\nBV:", Bv)
	aid, err := GetAid(Bv)
	if err != nil {
		return "", err
	}
	fmt.Println("\naid:", aid)
	url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/main?oid=%s&plat=1&seek_rpid=&type=1", aid)

	//method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	//defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//defer res.Body.Close()

	var result GetReplies

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	reply := ""
	for i, v := range result.Data.Replies {
		if i > 4 {
			break
		}
		reply += fmt.Sprintf("\n【%s】：%s\n", v.Member.Uname, v.Content.Message)
		fmt.Println(reply)
		for ii, vv := range v.Replies {
			if ii > 2 {
				break
			}
			reply += fmt.Sprintf("\t\tre:%s---【%s】\n", vv.Content.Message, vv.Member.Uname)
		}
	}
	return reply, nil
}

type getAid struct {
	Data data `json:"data"`
}
type data struct {
	Aid int `json:"aid"`
}

func GetAid(Bv string) (string, error) {

	furl := "https://api.bv-av.cn/get-bv-av?id="
	url := fmt.Sprint(furl, Bv)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("aid请求出错")
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	var getData getAid
	if err := json.Unmarshal(body, &getData); err != nil {
		println("unmarshal错误")
		return "", err
	}
	aid := getData.Data.Aid
	println("\nAid:", aid)
	return fmt.Sprintf("%d", aid), nil
}
