package MessageHandler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/WYC-RD/wxbot/source"
	"github.com/eatmoreapple/openwechat"
	"os"
	"strings"
)

type GroupMsgLog struct {
	MsgID     string
	GroupName string
	Sender    string
	MsgType   int
	Content   string
	MsgApp    string
	IsAt      bool
	URL       string
}
type MsgFile struct {
	MsgID  string
	Header string
	File   []byte
}

func DefaultHandler(message *openwechat.Message) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if message.IsAt() {
		gua, err := os.Open("./source/material/gua.png")
		if err != nil {
			fmt.Println("表情发送失败")
			return
		}
		message.ReplyImage(gua)
	}
	go func() {
		if message.MsgId != "" {
			err := msgLog(message)
			if err != nil {
				println("日志记录失败")
			}
		}
	}()
	if message.IsText() {
		textMessageHandler(message)
	}
	if message.MsgType == 49 {
		AppMessageHandler(message)
	}

}

func textMessageHandler(message *openwechat.Message) {
	switch message.Content {
	case "nba":
		nbaMessageHandler(message)
	case "热搜":
		message.ReplyText(source.GetHotSearch())

	}
}
func AppMessageHandler(message *openwechat.Message) {
	html5 := message.Content
	var appName string
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(html5))
	dom.Find("appname").Each(func(i int, selection *goquery.Selection) {
		appName = selection.Text()
	})
	if message.Url != "" {
		switch appName {
		case "哔哩哔哩":
			go bilibiliHandler(message)
		case "微博":
			go weiboHandler(message, appName)
		case "微博轻享版":
			go weiboHandler(message, appName)
		}

	}
}
