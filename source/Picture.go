package source

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
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
}

func (x *PicString) BackgroundInit(x0 int, y0 int, x1 int, y1 int, backgroundPath string) {
	x.Background = image.NewRGBA(image.Rect(x0, y0, x1, y1))
	tempbg, err := os.Open(backgroundPath)
	defer tempbg.Close()
	if err != nil {
		println("fail to load background")
		return
	}
	bg, err := png.Decode(tempbg)
	if err != nil {
		println("fail to decode background")
	}
	draw.Draw(x.Background, x.Background.Bounds(), bg, image.ZP, draw.Src)
}

func (x *PicString) ContextInit(DPI float64, bg *image.RGBA) {
	//x.Context = freetype.NewContext()

	//fontcolor :=image.Image
	// 设置像素密度
	x.Context.SetDPI(DPI)
	x.DPI = DPI

	// 指定画布对象
	x.Context.SetDst(bg)
	// 指定画布绘制范围
	x.Context.SetClip(bg.Bounds())

}
func (x *PicString) DrawRune(str string, padding int, font []byte, fontSize float64, c color.RGBA) {
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
			x.Pt.X = fixed.Int26_6(padding << 6)
			x.Pt.Y += x.Face.Metrics().Height + x.Face.Metrics().Height>>1
		} else if x.Font.Index(ch) == 0 {
			continue
		} else if x.Pt.X.Round()+wordWidth.Round() > x.Background.Rect.Dx()-padding {
			x.Pt.X = fixed.Int26_6(padding << 6)
			x.Pt.Y += x.Face.Metrics().Height + x.Face.Metrics().Height>>1
		}
		//fmt.Println("pt.x.round():", pt.X.Round(), "\nwordwidthRound:", wordWidth.Round(),
		//	"\nx.bg.rect.dx:", x.Background.Rect.Dx(), "\npadding:", padding)
		x.Pt, _ = x.Context.DrawString(string(ch), x.Pt)

	}
}
