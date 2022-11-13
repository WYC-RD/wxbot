package main

import (
	"fmt"
	"github.com/WYC-RD/wxbot/MessageHandler"
	"github.com/eatmoreapple/openwechat"
)

func main() {
	//bot := openwechat.DefaultBot()
	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 执行热登录
	bot.HotLogin(reloadStorage)

	// 注册消息处理函数
	bot.MessageHandler = MessageHandler.DefaultHandler
	//	func(msg *openwechat.Message) {
	//	if msg.IsText() && msg.Content == "ping" {
	//		msg.ReplyText("没完没了是吧")
	//	}
	//	if msg.IsText() && msg.Content == "换行测试" {
	//		msg.ReplyText("1 \n 2")
	//	}
	//	if msg.IsText() && msg.Content == "热搜" {
	//		rs := source.GetHotSearch()
	//		msg.ReplyText(rs)
	//	}
	//	if msg.IsText() && msg.Content == "nba" {
	//		sc := source.NbaScore()
	//		msg.ReplyText(sc)
	//	}
	//
	//	fmt.Println(msg.Context, "\ncontent:", msg.Content, msg.MsgType, msg.FromUserName,
	//		"\nismap?:", msg.IsMap(), "\nisText?:", msg.IsText(), "\nisCard?:", msg.IsCard())
	//}
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 登陆
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}
	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	// 获取所有的群组
	groups, err := self.Groups()
	fmt.Println(groups, err)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
