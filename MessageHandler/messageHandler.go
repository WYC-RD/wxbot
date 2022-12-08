package MessageHandler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/WYC-RD/wxbot/source"
	"github.com/eatmoreapple/openwechat"
	"log"
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

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		if message.IsSendByGroup() {
			if message.IsAt() {
				if err := source.OpHandle(message); err != nil {
					log.Println("opHandle fail")
					return
				}
			}
		} else {
			if err := source.OpHandle(message); err != nil {
				log.Println("opHandle fail")
				return
			}
		}
	}()

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
		hotsearchHandler(message)
		//message.ReplyText(source.GetHotSearch())

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
