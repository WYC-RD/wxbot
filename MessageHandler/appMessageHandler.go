package MessageHandler

import (
	"github.com/WYC-RD/wxbot/source"
	"github.com/eatmoreapple/openwechat"
)

func nbaMessageHandler(message *openwechat.Message) {
	message.ReplyText(source.NbaScore())
}
func bilibiliHandler(message *openwechat.Message, url string) {
	message.ReplyText(source.GetBvReplies(url))
}

func weiboHandler(message *openwechat.Message, url string) {
	message.ReplyText(source.GetWeiboReplies(url))
}
