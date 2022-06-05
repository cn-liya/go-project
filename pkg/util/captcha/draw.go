package captcha

import (
	"bytes"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand"
	"os"
	"time"
)

const (
	CharsUpperLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsLowerLetter = "abcdefghijklmnopqrstuvwxyz"
	CharsNumbers     = "0123456789"
	adjust           = 12
	fontSize         = adjust * 2
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Drawer struct {
	font   *truetype.Font
	bg     image.Image
	chars  []rune
	length int
}

// NewDrawer font为ttf格式必传，background为jpg格式默认纯灰，chars支持中文默认26个英文大写字母
func NewDrawer(font, background, chars string) *Drawer {
	b, err := os.ReadFile(font)
	if err != nil {
		panic(err)
	}
	f, err := freetype.ParseFont(b)
	if err != nil {
		panic(err)
	}
	d := &Drawer{font: f}
	if background != "" {
		if source, err := os.Open(background); err == nil {
			bg, _ := jpeg.Decode(source)
			source.Close()
			d.bg = bg
		}
	}
	if chars == "" {
		d.chars = []rune(CharsUpperLetter)
	} else {
		d.chars = []rune(chars)
	}
	d.length = len(d.chars)
	return d
}

// Generate 绘制验证码
func (d *Drawer) Generate(n int) (string, []byte) {
	height := fontSize * 2
	width := fontSize*(n+1) - adjust/2
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	if d.bg == nil {
		draw.Draw(canvas, canvas.Bounds(), image.NewUniform(color.Gray{Y: uint8(rand.Intn(256))}),
			image.Pt(0, 0), draw.Src)
	} else {
		x, y := 0, 0
		if bgx := d.bg.Bounds().Dx(); bgx > width {
			x = rand.Intn(bgx - width)
		}
		if bgy := d.bg.Bounds().Dy(); bgy > height {
			y = rand.Intn(bgy - height)
		}
		draw.Draw(canvas, canvas.Bounds(), d.bg, image.Pt(x, y), draw.Src)
	}

	c := freetype.NewContext()
	c.SetDst(canvas)
	c.SetClip(canvas.Bounds())
	c.SetFont(d.font)
	c.SetFontSize(fontSize)

	x := adjust / 2
	code := make([]rune, n)
	for i := range code {
		char := d.chars[rand.Intn(d.length)]
		code[i] = char
		pt := freetype.Pt(x+rand.Intn(adjust), fontSize+rand.Intn(adjust))
		c.SetSrc(randColor())
		c.DrawString(string(char), pt) // nolint
		x += fontSize
	}

	buf := bytes.NewBuffer(nil)
	_ = jpeg.Encode(buf, canvas, nil)
	return string(code), buf.Bytes()
}

func randColor() *image.Uniform {
	return image.NewUniform(color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	})
}
