package MessageHandler

import (
	"database/sql"
	"github.com/eatmoreapple/openwechat"
	_ "github.com/lib/pq"
	"io/ioutil"
	"strings"
)

func msgLog(message *openwechat.Message) error {
	db, err := sql.Open("postgres", "host=1.12.234.63 port=5432 user=postgres password=AAaa1111# dbname=postgres sslmode=disable")
	if err != nil {
		println("数据库链接失败", err)
		return err
	}
	defer db.Close()
	var Msg GroupMsgLog
	if message.IsSendByGroup() && message.MsgId != "" {
		sendGroup, _ := message.SenderInGroup()
		Msg.Sender = sendGroup.NickName
		group, _ := message.Sender()
		Msg.GroupName = group.NickName
		Msg.IsAt = message.IsAt()
		Msg.URL = message.Url
		Msg.MsgID = message.MsgId
		if message.IsText() {
			Msg.Content = message.Content
		} else if message.MsgType == 49 {
			Msg.Content = "AppID:" + message.AppInfo.AppID + " URL:" + message.Url
			Msg.MsgApp = message.AppInfo.AppID
		} else if message.HasFile() {
			switch message.MsgType {
			case openwechat.MsgTypeImage:
				Msg.Content = "图片消息"
			case openwechat.MsgTypeEmoticon:
				Msg.Content = "[动态表情]"
			case openwechat.MsgTypeVideo:
				Msg.Content = "视频消息"
			case openwechat.MsgTypeVoice:
				Msg.Content = "语言消息"
			}
			var msgFlie MsgFile
			msgFlie.MsgID = message.MsgId
			h, _ := message.GetFile()
			msgFlie.File, _ = ioutil.ReadAll(h.Body)
			for k, v := range h.Header {
				msgFlie.Header += k + ":" + strings.Join(v, " ") + "|"
			}
			stmt, err := db.Prepare("INSERT INTO msgfile (msg_id, header, file ) VALUES($1,$2,$3)")
			if err != nil {
				println("Prepare fail")
				return err
			}
			defer stmt.Close()
			res2, err := stmt.Exec(msgFlie.MsgID, msgFlie.Header, msgFlie.File)
			if err != nil {
				println("exec sql fail")
				println(res2)
				return err
			}
		}
	}

	//插入数据
	stmt, err := db.Prepare("INSERT INTO groupmsglog (msg_id, group_name, sender, msg_type, content, url, is_at) VALUES($1,$2,$3,$4,$5,$6,$7)")
	if err != nil {
		println("Prepare fail")
		return err
	}
	defer stmt.Close()
	res1, err := stmt.Exec(Msg.MsgID, Msg.GroupName, Msg.Sender, Msg.MsgType, Msg.Content, Msg.URL, Msg.IsAt)
	if err != nil {
		println("exec sql fail")
		println(res1)
		return err
	}
	return nil
}
