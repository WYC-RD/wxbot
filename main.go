package main

import (
	"fmt"
	"github.com/WYC-RD/wxbot/MessageHandler"
	"github.com/WYC-RD/wxbot/source"
	"github.com/eatmoreapple/openwechat"
	_ "github.com/xuri/excelize/v2"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式
	//注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	source.ConsoleQrCode(bot.UUID())
	go func() {
		bot.MessageHandler = MessageHandler.DefaultHandler
	}()
	hotlogin(true, bot)

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
func hotlogin(isHotlogin bool, bot *openwechat.Bot) {
	if isHotlogin {
		// 创建热存储容器对象
		reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
		// 执行热登录
		bot.HotLogin(reloadStorage)
		if err := bot.HotLogin(reloadStorage); err != nil {
			fmt.Println("热登录出错了")
			return
		}
	}
	if err := bot.Login(); err != nil {
		fmt.Println(err)
		return
	}
}
