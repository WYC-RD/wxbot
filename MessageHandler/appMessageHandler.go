package MessageHandler

import (
	"fmt"
	"github.com/WYC-RD/wxbot/source"
	"github.com/eatmoreapple/openwechat"
)

func nbaMessageHandler(message *openwechat.Message) {
	message.ReplyText(source.NbaScore())
}
func bilibiliHandler(message *openwechat.Message, url string) {
	reply, err := source.GetBvReplies(url)
	if err != nil {
		fmt.Printf("GetBvReplies fail", err)
	}
	message.ReplyText(reply)
}

func weiboHandler(message *openwechat.Message, url string, appname string) {
	message.ReplyText(source.GetWeiboReplies(url, appname))
}
