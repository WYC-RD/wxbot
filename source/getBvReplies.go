package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"regexp"
)

// 定义评论的结构体
type BvInfo struct {
	Data struct {
		Replies []replies
	}
	Other getAid
	Bv    string
	bchan chan *BvInfo
}
type replies struct {
	Member  Member    `json:"member"`
	Content Content   `json:"content"`
	Replies []replies `json:"replies"`
}
type Member struct {
	Uname string `json:"uname"`
	Sex   string `json:"sex"`
}
type Content struct {
	Message string `json:"message"`
}
type getAid struct {
	Info Info `json:"data"`
}
type Info struct {
	Aid   int    `json:"aid"`
	Title string `json:"title"`
	Pic   string `json:"pic"`
	Image image.Image
	URL   string
}

// 定义使用的字体和颜色
var Blue = color.RGBA{6, 174, 236, 255}
var Pink = color.RGBA{251, 114, 153, 255}
var FzHeiTi, _ = ioutil.ReadFile("./source/material/FangZhengHeiTiJianTi-1.ttf")
var smileHeiTi, _ = ioutil.ReadFile("./source/material/SmileySans-Oblique.ttf")
var DongQing, _ = ioutil.ReadFile("./source/material/冬青黑体.ttf")
var DongQingW4, _ = ioutil.ReadFile("./source/material/冬青黑体简体中文 W4.ttf")
var DongQingW5, _ = ioutil.ReadFile("./source/material/冬青黑体简体中文 W5.ttf")
var codeSize = 400

// 往Bv结构体中写入回复内容
func (bvinfo *BvInfo) GetBvReplies(URL string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	var err error
	bvinfo.Other.Info.URL = URL
	bvinfo.Bv, err = GetBvId(URL)
	if err != nil {
		println("GetBvID fail")
		return err
	}
	//fmt.Println("\nBV:", Bv)
	aid, err := bvinfo.GetAid()
	if err != nil {
		return err
	}
	//fmt.Println("\naid:", getaid.Data.Aid)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/v2/reply/main?oid="+aid+"&plat=1&seek_rpid=&type=1", nil)
	//println(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()
	//var result GetReplies
	if err1 := json.Unmarshal(body, bvinfo); err1 != nil {
		return err
	}
	return err
}

// 获取往Bv结构体中写入AID、封面图、标题
func (bvinfo *BvInfo) GetAid() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.bv-av.cn/get-bv-av?id="+bvinfo.Bv, nil)
	if err != nil {
		fmt.Println("aid请求出错")
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()
	//var getData getAid
	if err := json.Unmarshal(body, &bvinfo.Other); err != nil {
		println("unmarshal错误")
		return "", err
	}
	//aid := getData.Data.Aid
	//println("\nAid:", aid)
	return fmt.Sprintf("%d", bvinfo.Other.Info.Aid), nil
}

func (bvinfo *BvInfo) genBvPic(picString PicString) (image.Image, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	picString.DrawRune(bvinfo.Other.Info.Title, smileHeiTi, 32, Pink)
	picString.DrawRune("\n", DongQing, 20, Pink)
	//获取缩略图
	picIO, err := http.Get(bvinfo.Other.Info.Pic)
	if err != nil {
		return nil, err
	}
	defer picIO.Body.Close()
	rawPic, _, err := image.Decode(picIO.Body)
	if err != nil {
		err = nil
		rawPic, err = jpeg.Decode(picIO.Body)
		if err != nil {
			fmt.Println("图片解码失败")
			return nil, err
		}
	}
	picWidth := picString.Background.Bounds().Dx() - 2*picString.Padding
	size := float64(rawPic.Bounds().Dx()) / float64(picWidth)
	//fmt.Println("picWidth:", picWidth, "\nrawdx:", rawPic.Bounds().Dx(), "\nsize:", size)
	if size != 0 {
		bvinfo.Other.Info.Image = resize.Resize(uint(float64(rawPic.Bounds().Dx())/size), uint(float64(rawPic.Bounds().Dy())/size), rawPic, resize.Lanczos3)
		//粘贴缩略图
		//draw.Draw函数中，目标图的想要向右下移动的话，x、y为负数
		point := image.Point{int(picString.Pt.X>>6) * -1, int(picString.Pt.Y>>6) * -1}
		draw.Draw(picString.Background, picString.Background.Bounds(), bvinfo.Other.Info.Image, point, draw.Src)
		//修改PT的Y坐标，避免文字和缩略图区域重合
		picString.Pt.Y += fixed.Int26_6(bvinfo.Other.Info.Image.Bounds().Dy() << 6)
	}
	picString.DrawRune("\n", DongQing, 20, Pink)
	picString.DrawRune("\n", DongQing, 12, Pink)
	for i, v := range bvinfo.Data.Replies {
		if i > 2 {
			break
		}
		uname := fmt.Sprintf("%s  ：", v.Member.Uname)
		picString.DrawRune(uname, DongQing, 17, Pink)
		message := fmt.Sprintf("%s\n", v.Content.Message)
		picString.DrawRune(message, DongQing, 17, Blue)
		fmt.Println(uname, message)
		for ii, vv := range v.Replies {
			if ii > 2 {
				break
			}
			suname := fmt.Sprintf("------[%s] reply：", vv.Member.Uname)
			picString.DrawRune(suname, DongQing, 17, Pink)
			smessage := fmt.Sprintf("%s\n", vv.Content.Message)
			picString.DrawRune(smessage, DongQing, 17, Blue)
		}
		picString.DrawRune("\n\n", DongQing, 15, Blue)
	}
	picString.SubImg, err = appendQr(*picString.Background, picString, bvinfo.Other.Info.URL, color.RGBA{255, 255, 255, 255}, color.RGBA{116, 125, 140, 255})
	if err != nil {
		return nil, err
	}

	picString.LastY = int(picString.Pt.Y >> 6)
	println("lasty:", picString.LastY)
	return picString.SubImg, nil
}

// 正则获取bv
func GetBv(url string) string {

	regx1 := regexp.MustCompile(`[^https://www.bilibili.com/video/].*?/`)

	bvid1 := regx1.FindAllString(url, -1)
	bvid2 := fmt.Sprint(bvid1)
	//println(bvid2)
	//println(len(bvid2))
	if len(bvid2) > 10 {
		bvid3 := fmt.Sprint(bvid2[1 : len(bvid2)-2])
		fmt.Println(bvid3)
		return bvid3
	}
	return ""
}

// 获取短链重定向之后的链接，方便拿到一些原始信息
func GetBvId(url string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	//if err != nil {
	//	return "", err
	//}
	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	response, _ := client.Do(req)
	//if err != nil {
	//	fmt.Printf("fail to get response on GetBvID\n")
	//	return "", err
	//}
	defer response.Body.Close()
	redirectUrl, _ := response.Location()
	redirectUrl1 := fmt.Sprint(redirectUrl)
	//fmt.Printf(redirectUrl1)
	if bvid := GetBv(redirectUrl1); bvid != "" {
		return bvid, nil
	}

	return GetBvId(redirectUrl1)
}
