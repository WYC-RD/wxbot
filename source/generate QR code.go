package source

import (
	"fmt"
	"github.com/skip2/go-qrcode"
	"image/color"
)

func myEncode(content string, level qrcode.RecoveryLevel, size int, background color.RGBA, foreground color.RGBA) ([]byte, error) {
	var q *qrcode.QRCode
	q, err := qrcode.New(content, level)
	q.BackgroundColor = background
	q.ForegroundColor = foreground
	if err != nil {
		return nil, err
	}
	return q.PNG(size)
}
func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Highest)
	fmt.Println(q.ToString(true))
}
