package source

import (
	"encoding/json"
	"fmt"
	"github.com/golang/freetype"
	"image"
	"image/color"
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

var Blue = color.RGBA{6, 174, 236, 255}
var Pink = color.RGBA{251, 114, 153, 255}
var FzHeiTi, _ = ioutil.ReadFile("/Users/wangzehong/Pictures/fonts/FangZhengHeiTiJianTi-1.ttf")
var smileHeiTi, _ = ioutil.ReadFile("/Users/wangzehong/Pictures/fonts/SmileySans-Oblique.ttf")

func GetBvReplies(URL string, picString PicString) (PicString, error) {

	Bv, err := GetBvId(URL)
	if err != nil {
		println("GetBvID fail")
		return picString, err
	}
	fmt.Println("\nBV:", Bv)
	getaid, err := GetAid(Bv)
	if err != nil {
		return picString, err
	}
	fmt.Println("\naid:", getaid.Data.Aid)
	url := fmt.Sprintf("https://api.bilibili.com/x/v2/reply/main?oid=%d&plat=1&seek_rpid=&type=1", getaid.Data.Aid)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	println(url)
	if err != nil {
		fmt.Println(err)
		return picString, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return picString, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return picString, err
	}
	defer res.Body.Close()

	var result GetReplies
	if err1 := json.Unmarshal(body, &result); err1 != nil {
		return picString, err
	}
	picString.DrawRune(getaid.Data.Title+"\n", 40, smileHeiTi, 28, Pink)
	for i, v := range result.Data.Replies {
		if i > 2 {
			break
		}
		uname := fmt.Sprintf("%s  ：", v.Member.Uname)
		picString.DrawRune(uname, 40, FzHeiTi, 17, Pink)
		message := fmt.Sprintf("%s\n", v.Content.Message)
		picString.DrawRune(message, 40, FzHeiTi, 17, Blue)
		fmt.Println(uname, message)
		for ii, vv := range v.Replies {
			if ii > 2 {
				break
			}
			//picString.DrawRune("------",10, FzHeiTi, 40,Pink)
			suname := fmt.Sprintf("------[%s] reply：", vv.Member.Uname)
			picString.DrawRune(suname, 40, FzHeiTi, 17, Pink)
			smessage := fmt.Sprintf("%s\n", vv.Content.Message)
			picString.DrawRune(smessage, 40, FzHeiTi, 17, Blue)
			//reply += fmt.Sprintf("[%s] reply：%s\n", vv.Member.Uname, vv.Content.Message)
		}
		enter := fmt.Sprint("\n\n")
		picString.DrawRune(enter, 40, FzHeiTi, 15, Blue)
	}
	return picString, nil
}

type getAid struct {
	Data data `json:"data"`
}
type data struct {
	Aid   int    `json:"aid"`
	Title string `json:"title"`
	Pic   string `json:"pic"`
	Image image.Image
}

func GetAid(Bv string) (*getAid, error) {

	furl := "https://api.bv-av.cn/get-bv-av?id="
	url := fmt.Sprint(furl, Bv)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("aid请求出错")
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	var getData getAid
	if err1 := json.Unmarshal(body, &getData); err1 != nil {
		println("unmarshal错误")
		return nil, err
	}
	//aid := getData.Data.Aid
	picIO, _ := http.Get(getData.Data.Pic)
	getData.Data.Image, _, _ = image.Decode(picIO.Body)
	//println("\nAid:", aid)
	return &getData, nil
}

//var repliesPic PicString

func GetBvRepliesPic(URL string) (*image.RGBA, error) {
	var err error
	repliesPic := PicString{}
	repliesPic.Context = freetype.NewContext()
	repliesPic.Context = freetype.NewContext()
	repliesPic.BackgroundInit(0, 0, 1080, 1920, "/Users/wangzehong/Pictures/bilibili.png")
	repliesPic.ContextInit(200, repliesPic.Background)
	repliesPic.Pt = freetype.Pt(40, 40+int(repliesPic.Context.PointToFixed(40)>>6))
	repliesPic, err = GetBvReplies(URL, repliesPic)
	if err != nil {
		return nil, err
	}
	//repliesPic.DrawRune(replies, 10, FzHeiTi, 40,Blue)
	return repliesPic.Background, nil
	//return repliesPic.Background
}
