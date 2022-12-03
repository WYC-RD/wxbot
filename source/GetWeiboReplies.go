package source

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type wbInfo struct {
	URL     string
	App     string
	ID      string
	Pic     image.RGBA
	Text    string
	OK      int  `json:"ok"`
	Replies Data `json:"data"`
	Detail  Datail
}
type Datail struct {
	Data Wbdata `json:"data"`
}
type Wbdata struct {
	LongTextContent string `json:"longTextContent"`
	Reposts_count   int    `json:"reposts_count"`
	Comments_count  int    `json:"comments_count"`
	Attitudes_count int    `json:"attitudes_count"`
}
type Data struct {
	Data []WbReplies `json:"data"`
}

type WbReplies struct {
	Text     string     `json:"text"`
	Source   string     `json:"source"`
	Like     int        `json:"like_count"`
	User     User       `json:"user"`
	Comments []Comments `json:"comments"`
}
type Comments struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	User   User   `json:"user"`
}
type User struct {
	ScreenName string `json:"screen_name"`
	Gender     string `json:"gender"`
}

var Black = color.RGBA{47, 53, 66, 255}
var Gray = color.RGBA{116, 125, 140, 255}
var Brown = color.RGBA{89, 75, 65, 255}
var BoldBrown = color.RGBA{89, 64, 46, 255}
var LiteBrown = color.RGBA{89, 79, 72, 255}

func (wb *wbInfo) GetWeiboReplies() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://m.weibo.cn/comments/hotflow?mid="+wb.ID, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err := json.Unmarshal(body, wb); err != nil {
		if len(wb.Replies.Data) == 0 {
			fmt.Println("Replies Unmarshal fail")
			return err
		}
	}

	//fmt.Println("rawreplies:", wb.Replies.Data)
	return nil
}

func (wb *wbInfo) GetWbDetail() error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	//fmt.Println("ID:", wb.ID)
	clent := &http.Client{}
	req, err := http.NewRequest("GET", "https://m.weibo.cn/statuses/extend?id="+wb.ID, nil)
	if err != nil {
		return err
	}
	res, err := clent.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	//fmt.Println("body:", string(body))
	if err := json.Unmarshal(body, &wb.Detail); err != nil {
		wb.Detail.Data.LongTextContent = "无法获取获取该微博正文"
		return err
	}
	//fmt.Println(wb.Detail.Data)
	//fmt.Println("rawText:", wb.Detail.Data.LongTextContent)
	regxbr, _ := regexp.Compile(`<br /><br />`)
	wb.Detail.Data.LongTextContent = regxbr.ReplaceAllString(wb.Detail.Data.LongTextContent, "\n\n")
	regx1, _ := regexp.Compile(`<.+?>`)
	wb.Detail.Data.LongTextContent = regx1.ReplaceAllString(wb.Detail.Data.LongTextContent, "")
	//println("longtextcontent:", wb.Detail.Data.LongTextContent)
	return nil
}

func (wb *wbInfo) genWbPic(picString PicString) (image.Image, error) {
	if strings.Count(wb.Detail.Data.LongTextContent, "") > 500 {
		picString.DrawRune("\n"+wb.Detail.Data.LongTextContent, DongQing, 17, Brown)
	} else if strings.Count(wb.Detail.Data.LongTextContent, "") > 40 {
		picString.DrawRune("\n"+wb.Detail.Data.LongTextContent, DongQing, 20, Brown)
	} else if strings.Count(wb.Detail.Data.LongTextContent, "") > 15 {
		picString.DrawRune("\n"+wb.Detail.Data.LongTextContent, DongQingW5, 30, Brown)
	} else {
		picString.DrawRune("\n"+wb.Detail.Data.LongTextContent, DongQingW5, 40, Brown)
	}

	picString.DrawRune("\n\n", DongQing, 20, Pink)

	for i, v := range wb.Replies.Data {
		if i > 2 {
			break
		}
		regx1, _ := regexp.Compile(`<.+?>`)
		v.Text = regx1.ReplaceAllString(v.Text, "")
		picString.DrawRune(v.User.ScreenName, DongQingW5, 17, BoldBrown)
		picString.DrawRune("["+v.Source+"]", DongQingW4, 17, BoldBrown)
		picString.DrawRune("：", DongQing, 17, LiteBrown)
		picString.DrawRune(v.Text+"\n", DongQing, 17, LiteBrown)
		for ii, vv := range v.Comments {
			if ii > 2 {
				break
			}
			regx1, _ := regexp.Compile(`<.+?>`)
			vv.Text = regx1.ReplaceAllString(vv.Text, "")
			picString.DrawRune("------"+vv.User.ScreenName, DongQingW5, 15, BoldBrown)
			picString.DrawRune("["+vv.Source+"]", DongQingW4, 15, BoldBrown)
			picString.DrawRune("：", DongQing, 15, LiteBrown)
			picString.DrawRune(vv.Text+"\n", DongQing, 15, LiteBrown)
		}
		picString.DrawRune("\n", DongQing, 15, Pink)
	}
	picString.DrawRune("\n\n", DongQing, 15, color.RGBA{207, 201, 196, 255})
	//生成二维码
	picString.SubImg = appendQr(*picString.Background, picString, wb.URL, color.RGBA{251, 252, 245, 255}, color.RGBA{207, 201, 196, 255})
	picString.LastY = int(picString.Pt.Y >> 6)
	return picString.SubImg, nil
}

// 获取短链重定向之后的链接，方便拿到一些原始信息
