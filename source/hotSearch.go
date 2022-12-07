package source

import (
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetHotSearch(...interface{}) (string, error) {
	//请求微博热搜
	url := "https://weibo.com/ajax/side/hotSearch"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	resquest, err := http.NewRequest("GET", url, Body)
	if err != nil {
		return "", err
	}
	//发送请求
	res, err := client.Do(resquest)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("err!")
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	//用第三发库解析Json为go数据结构
	rez, err := simplejson.NewJson([]byte(string(body)))
	if err != nil {
		return "", err
	}
	//根据返回来的json结构，取到数组
	hotsear, err := rez.Get("data").Get("realtime").Array()
	if err != nil {
		return "", err
	}
	//遍历数组
	var slice_rs []string
	for _, hs := range hotsear {
		//每一组就是一条热搜
		if each_map, ok := hs.(map[string]interface{}); ok {
			rs := fmt.Sprintf("%s:%s\n", each_map["rank"], each_map["word"])
			slice_rs = append(slice_rs, rs)
		}
	}
	return strings.Join(slice_rs, ""), nil
}
