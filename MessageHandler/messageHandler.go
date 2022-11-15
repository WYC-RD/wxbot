package MessageHandler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/WYC-RD/wxbot/source"
	"github.com/eatmoreapple/openwechat"
	"strings"
)

func DefaultHandler(message *openwechat.Message) {
	if message.IsText() {
		textMessageHandler(message)
	}
	if message.MsgType == 49 {
		fmt.Println(message.Context, "\ncontent:", message.Content, message.MsgType, message.FromUserName,
			"\nismap?:", message.IsMap(), "\nisText?:", message.IsText(), "\nisCard?:", message.IsCard())
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
	var appName, url string
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(html5))
	dom.Find("appname").Each(func(i int, selection *goquery.Selection) {
		appName = selection.Text()

		//if  appName == "哔哩哔哩" {
		//	dom.Find("url").Each(func(i int, selection *goquery.Selection) {
		//		url := selection.Text()
		//		message.ReplyText(source.GetBvReplies(url))
		//	})
		//}
		//if  appName == ""
	})
	dom.Find("url").Each(func(i int, selection *goquery.Selection) {
		url = selection.Text()
		//message.ReplyText(source.GetBvReplies(url))
	})
	switch appName {
	case "哔哩哔哩":
		bilibiliHandler(message, url)
	case "微博":
		weiboHandler(message, url)
	case "微博轻享版":
		weiboHandler(message, url)
	}
}
