package source

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

type PicString struct {
	Font       *truetype.Font
	Background *image.RGBA
	Context    *freetype.Context
	Face       font.Face
	DPI        float64
	Pt         fixed.Point26_6
	Padding    int
	SubImg     image.Image
	LastY      int
}

func (x *PicString) BackgroundInit(x0 int, y0 int, backgroundPath string) error {
	tempbg, err := os.Open(backgroundPath)
	defer tempbg.Close()
	if err != nil {
		println("fail to load background")
		return err
	}
	bg, err := png.Decode(tempbg)
	if err != nil {
		println("fail to decode background")
		return err
	}
	x.Background = image.NewRGBA(image.Rect(x0, y0, bg.Bounds().Dx(), bg.Bounds().Dy()))
	draw.Draw(x.Background, x.Background.Bounds(), bg, image.ZP, draw.Src)
	return nil
}

func (x *PicString) ContextInit(DPI float64, bg *image.RGBA) {
	// 设置像素密度
	x.Context.SetDPI(DPI)
	x.DPI = DPI

	// 指定画布对象
	x.Context.SetDst(bg)
	// 指定画布绘制范围
	x.Context.SetClip(bg.Bounds())

}
func (x *PicString) DrawRune(str string, font []byte, fontSize float64, c color.RGBA) {
	//println("DX\n\n", x.Background.Rect.Dx())
	color := image.Uniform{c}
	f, err := truetype.Parse(font)
	if err != nil {
		println("fail to parse ttf ")
	}
	x.Font = f
	// 指定字体
	x.Context.SetFont(f)
	// 指定文字颜色
	x.Context.SetSrc(&color)
	// 指定字体大小
	x.Context.SetFontSize(fontSize)
	// 指定字体宽度
	opts := truetype.Options{}
	opts.Size = fontSize
	opts.DPI = x.DPI
	x.Face = truetype.NewFace(f, &opts)

	//x.Pt = freetype.Pt(padding, padding+int(x.Context.PointToFixed(15)>>6))
	for _, ch := range []rune(str) {
		wordWidth, _ := x.Face.GlyphAdvance(ch)
		if ch == '\t' {
			x.Pt.X += +(2 * wordWidth)
		} else if ch == '\n' {
			x.Pt.X = fixed.Int26_6(x.Padding << 6)
			x.Pt.Y += x.Face.Metrics().Height + x.Face.Metrics().Height>>1
			continue
		} else if x.Font.Index(ch) == 0 {
			continue
		} else if x.Pt.X.Round()+wordWidth.Round() > x.Background.Rect.Dx()-x.Padding {
			x.Pt.X = fixed.Int26_6(x.Padding << 6)
			x.Pt.Y += x.Face.Metrics().Height + x.Face.Metrics().Height>>1
		} else if x.Pt.Y.Round() >= x.Background.Rect.Dy() {
			x.Background.Bounds().Add(image.Point{0, 4 * int(x.Face.Metrics().Height>>6)})
		}

		//fmt.Println("pt.x.round():", pt.X.Round(), "\nwordwidthRound:", wordWidth.Round(),
		//	"\nx.bg.rect.dx:", x.Background.Rect.Dx(), "\npadding:", padding)
		x.Pt, _ = x.Context.DrawString(string(ch), x.Pt)
		//fmt.Println("\nPT.x:", int(x.Pt.X>>6), "Pt.Y:", int(x.Pt.Y>>6))

	}
}
func PicInit(backgroundPath string) (*PicString, error) {
	repliesPic := PicString{}
	repliesPic.Context = freetype.NewContext()
	repliesPic.Context = freetype.NewContext()
	//加载背景图
	if err := repliesPic.BackgroundInit(0, 0, backgroundPath); err != nil {
		fmt.Println("PicInit fail")
		return nil, err
	}
	//加载文字区域画布
	repliesPic.ContextInit(200, repliesPic.Background)
	//设置留白
	repliesPic.Padding = 40
	//设置文字渲染起点像素坐标
	repliesPic.Pt = freetype.Pt(repliesPic.Padding, repliesPic.Padding*6)
	return &repliesPic, nil
}
func appendQr(rgba image.RGBA, picString PicString, URL string, bgColor color.RGBA, qrColor color.RGBA) image.Image {
	code, _ := myEncode(URL, qrcode.Medium, codeSize, bgColor, qrColor)
	qrcode, _ := png.Decode(bytes.NewReader(code))
	point := image.Point{(picString.Background.Bounds().Dx()/2 - codeSize/2) * -1, (int(picString.Pt.Y>>6) - codeSize/2 + codeSize/4) * -1}
	draw.Draw(&rgba, rgba.Bounds(), qrcode, point, draw.Src)
	SubImg := rgba.SubImage(image.Rectangle{image.Point{0, 0},
		image.Point{picString.Background.Bounds().Dx(), codeSize + int(picString.Pt.Y>>6)}})
	fmt.Println("point", image.Point{picString.Background.Bounds().Dx(), 2*codeSize + int(picString.Pt.Y>>6)})
	return SubImg
}

//b站颜色
// color.RGBA{255, 255, 255, 255}, color.RGBA{116, 125, 140, 255}
