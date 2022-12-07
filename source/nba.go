package source

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strings"
)

func NbaScore() (string, error) {

	url := "https://api.sports.163.com/api/nba/v2/schedule/getRecent?product=pc"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	rez, err := simplejson.NewJson([]byte(string(body)))
	if err != nil {
		return "", err
	}
	Score, err := rez.Get("data").Get("1").Array()
	if err != nil {
		return "", err
	}
	//遍历数组
	var slice_sc []string
	for _, sc := range Score {
		if each_map, ok := sc.(map[string]interface{}); ok {
			rs := fmt.Sprintf("%s vs %s\n%s - %s\n\n", each_map["away"], each_map["home"], each_map["awayScore"], each_map["homeScore"])
			slice_sc = append(slice_sc, rs)
		}
	}
	return strings.Join(slice_sc, ""), nil
}
