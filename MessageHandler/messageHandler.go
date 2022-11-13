package MessageHandler

import (
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
		AppMessageHandler(message)
	}
}

func textMessageHandler(message *openwechat.Message) {
	switch message.Content {
	case "nba":
		nbaMessageHandler(message)
	}
}
func AppMessageHandler(message *openwechat.Message) {
	html5 := message.Content
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(html5))
	dom.Find("appname").Each(func(i int, selection *goquery.Selection) {
		if appName := selection.Text(); appName == "哔哩哔哩" {
			dom.Find("url").Each(func(i int, selection *goquery.Selection) {
				url := selection.Text()
				message.ReplyText(source.GetBvReplies(url))
			})
		}
	})
}
