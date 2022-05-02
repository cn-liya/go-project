package captcha

import (
	"bytes"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand/v2"
	"os"
	"unicode/utf8"
)

const (
	unitSize = 6
	adjust   = unitSize * 2
	fontSize = unitSize * 4
	height   = fontSize * 2
)

type Drawer struct {
	font *opentype.Font
}

// NewDrawer font为ttf或otf格式字体文件路径
func NewDrawer(font string) *Drawer {
	b, err := os.ReadFile(font)
	if err != nil {
		panic(err)
	}
	f, err := opentype.Parse(b)
	if err != nil {
		panic(err)
	}
	return &Drawer{font: f}
}

// Draw s字符长度2~8为宜
func (d *Drawer) Draw(s string) []byte {
	width := fontSize*utf8.RuneCountInString(s) + 3*unitSize
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	gray := image.NewUniform(color.Gray{Y: uint8(rand.IntN(256))})
	draw.Draw(canvas, canvas.Bounds(), gray, image.Point{}, draw.Src)

	face, _ := opentype.NewFace(d.font, &opentype.FaceOptions{
		Size: fontSize,
		DPI:  72,
	})
	c := &font.Drawer{Dst: canvas, Face: face}
	seed := rand.Perm(256)
	pos := 0
	x := unitSize
	for _, char := range s {
		c.Src = image.NewUniform(color.RGBA{
			R: uint8(seed[pos&255]),
			G: uint8(seed[(pos+1)&255]),
			B: uint8(seed[(pos+2)&255]),
		})
		c.Dot = fixed.P(x+rand.IntN(adjust), fontSize+rand.IntN(adjust)) // x+范围[0,2]单位, y范围[4,6]单位
		c.DrawString(string(char))
		x += fontSize
		pos += 3
	}

	for x = unitSize; x < width; x += unitSize {
		switch rand.IntN(5) {
		case 0, 4:
			drawCircle(canvas, x+rand.IntN(unitSize), rand.IntN(height), rand.IntN(3)+2, color.RGBA{
				R: uint8(seed[pos&255]),
				G: uint8(seed[(pos+1)&255]),
				B: uint8(seed[(pos+2)&255]),
			})
			pos += 3
		case 1, 3:
			w := rand.IntN(3) + 2
			x0, y0 := x+rand.IntN(unitSize), rand.IntN(height)
			draw.Draw(canvas, image.Rect(x0-w, y0, x0, y0+6-w), image.NewUniform(color.RGBA{
				R: uint8(seed[pos&255]),
				G: uint8(seed[(pos+1)&255]),
				B: uint8(seed[(pos+2)&255]),
			}), image.Point{}, draw.Src)
			pos += 3
		}
	}

	buf := bytes.NewBuffer(nil)
	_ = jpeg.Encode(buf, canvas, nil)
	return buf.Bytes()
}

func drawCircle(dst draw.Image, x0, y0, r int, c color.Color) {
	x, y := r, 0
	t1 := r >> 4
	for x >= y {
		dst.Set(x0+x, y0+y, c)
		dst.Set(x0+y, y0+x, c)
		dst.Set(x0-y, y0+x, c)
		dst.Set(x0-x, y0+y, c)
		dst.Set(x0-x, y0-y, c)
		dst.Set(x0-y, y0-x, c)
		dst.Set(x0+y, y0-x, c)
		dst.Set(x0+x, y0-y, c)
		y++
		t1 += y
		if t2 := t1 - x; t2 >= 0 {
			t1 = t2
			x--
		}
	}
}
