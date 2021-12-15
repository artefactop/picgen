package image

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// Label struct
type Label struct {
	freetype.Context

	Text   string
	Color  color.Color
	Width  int
	Height int
}

// DrawLabel Draw label l on position (x,y) of m Image
func DrawLabel(m draw.Image, l Label, x, y int) error {
	l.SetSrc(image.NewUniform(l.Color))
	l.SetClip(m.Bounds())
	l.SetDst(m)

	// Set the l centered alignment to x,y
	pt := freetype.Pt(x, y)
	if _, err := l.DrawString(l.Text, pt); err != nil {
		return err
	}
	return nil
}

// NewLabel return a new *Label
func NewLabel(text string, color color.Color, dpi float64, f *truetype.Font, fontSize float64) *Label {
	ctx := freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(f)
	ctx.SetFontSize(fontSize)

	opts := truetype.Options{Size: fontSize, DPI: dpi}
	face := truetype.NewFace(f, &opts)
	d := &font.Drawer{
		Face: face,
	}
	w := int(d.MeasureString(text) >> 6)
	h := int(face.Metrics().Height >> 6)

	l := &Label{
		Context: *ctx,
		Text:    text,
		Color:   color,
		Width:   w,
		Height:  h,
	}

	return l
}

// NewImage return a new image.Image of w x h size and c Color
func NewImage(w, h int, c color.Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: c}, image.Point{X: 0, Y: 0}, draw.Src)
	return img
}

// Encode writes the Image m to w in the f format (PNG, JPEG). Any Image may be
// encoded, but images that are not image.NRGBA might be encoded lossily.
func Encode(w io.Writer, m image.Image, f string) (int, error) {
	buffer := new(bytes.Buffer)

	switch f {
	case "JPEG":
		if err := jpeg.Encode(buffer, m, &jpeg.Options{}); err != nil {
			return 0, err
		}
	case "PNG":
		fallthrough
	default:
		if err := png.Encode(buffer, m); err != nil {
			return 0, err
		}
	}

	n, err := w.Write(buffer.Bytes())
	if err != nil {
		return 0, fmt.Errorf("can not encode image %w", err)
	}

	return n, nil
}
