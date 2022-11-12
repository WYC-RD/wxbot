package MessageHandler

import (
	"github.com/eatmoreapple/openwechat"
)

func DefaultHandler(message *openwechat.Message)  {
	if message.IsText() {
		textMessageHandler(message)
	}
}

func textMessageHandler(message *openwechat.Message)  {
	switch message.Content {
	case "nba":
		nbaMessageHandler(message)
	}
}