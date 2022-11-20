package MessageHandler

import (
	"fmt"
	"github.com/WYC-RD/wxbot/source"
	"github.com/eatmoreapple/openwechat"
	"image/png"
	"os"
)

func nbaMessageHandler(message *openwechat.Message) {
	message.ReplyText(source.NbaScore())
}
func bilibiliHandler(message *openwechat.Message, url string) {
	replyPic, err := source.BilibiliPic(url)
	if err != nil {
		fmt.Printf("GetBvReplies fail", err)
	}
	//message.ReplyText(reply)
	picFlie, err := os.Create("/Users/wangzehong/Pictures/wxbot_Bilibi.png")
	defer picFlie.Close()
	if err != nil {
		fmt.Println("creat bilibili replies picture fail")
	}
	png.Encode(picFlie, replyPic)
	pic2, err := os.Open("/Users/wangzehong/Pictures/wxbot_Bilibi.png")
	if err != nil {
		println("加载图片失败")
	}
	defer pic2.Close()
	message.ReplyImage(pic2)
}

func weiboHandler(message *openwechat.Message, url string, appname string) {
	message.ReplyText(source.GetWeiboReplies(url, appname))
}
