package picgen

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// RootHandler ...
func RootHandler(w http.ResponseWriter, req *http.Request) {
	img, err := buildImage(200, 100, "Hi there!", 16)
	if err != nil {
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, img); err != nil {
		log.Fatal(err)
	}
}

func buildImage(width, height int, label string, fontSize float64) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	clr := color.RGBA{100, 200, 200, 255}

	draw.Draw(img, img.Bounds(), &image.Uniform{clr}, image.ZP, draw.Src)

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}
	x, y := width/2, height/2
	dpi := 72.0
	addLabel(img, font, fontSize, dpi, label, x, y)

	return img, nil
}

func addLabel(img *image.RGBA, drawFont *truetype.Font, fontSize, dpi float64, label string, x, y int) error {
	fontForeGroundColor := image.NewUniform(color.Black)
	ctx := freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(drawFont)
	ctx.SetFontSize(fontSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(fontForeGroundColor)

	opts := truetype.Options{Size: fontSize}
	face := truetype.NewFace(drawFont, &opts)
	d := &font.Drawer{
		Face: face,
	}
	labelWidht := int(d.MeasureString(label) >> 6)
	labelHeight := int(face.Metrics().Height >> 6)

	// Set the label centered alignment to x,y
	pt := freetype.Pt(x-labelWidht/2, y+labelHeight/2)
	if _, err := ctx.DrawString(label, pt); err != nil {
		return err
	}
	return nil

}
