package MessageHandler

import (
	"fmt"
	"github.com/WYC-RD/wxbot/source"
	"github.com/eatmoreapple/openwechat"
	"image"
	"image/png"
	"os"
)

func nbaMessageHandler(message *openwechat.Message) {
	message.ReplyText(source.NbaScore())
}
func bilibiliHandler(message *openwechat.Message) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	replyPic, err := source.BilibiliPic(message.Url)
	if err != nil {
		fmt.Printf("GetBvReplies fail", err)
	}
	//message.ReplyText(reply)
	pn := fmt.Sprintf("./wxbot-pic-log/bilibili/%s-%d.png", message.MsgId, message.CreateTime)
	picFlie, err := os.Create(pn)
	defer picFlie.Close()
	if err != nil {
		fmt.Println("creat bilibili replies picture fail")
	}
	png.Encode(picFlie, *replyPic)
	pic2, err := os.Open(pn)
	if err != nil {
		println("加载图片失败")
	}
	defer pic2.Close()
	message.ReplyImage(pic2)
}

func weiboHandler(message *openwechat.Message, appname string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	wchan := make(chan *image.Image)
	go func() {
		replyPic, err := source.Wbhandle(message.Url, appname)
		if err != nil {
			fmt.Printf("GetWbReplies fail", err)
		}
		wchan <- replyPic
	}()

	pn := fmt.Sprintf("./wxbot-pic-log/weibo/%s-%d.png", message.MsgId, message.CreateTime)
	picFlie, err := os.Create(pn)
	defer picFlie.Close()
	if err != nil {
		fmt.Println("creat weibo replies picture fail")
	}
	replyPic := <-wchan
	png.Encode(picFlie, *replyPic)
	pic2, err := os.Open(pn)
	if err != nil {
		println("加载图片失败")
	}
	defer pic2.Close()
	message.ReplyImage(pic2)
}
