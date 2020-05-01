package picgen

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type options struct {
	Width, Height int
	Color         color.RGBA
	Format        string
	LabelOptions  *labelOptions
}

type labelOptions struct {
	Text  string
	Color color.RGBA
	Font  *truetype.Font
	DPI   float64
	Size  float64
	X     int
	Y     int
}

func newImage(op *options) (image.Image, error) {
	img := image.NewRGBA(image.Rect(0, 0, op.Width, op.Height))

	draw.Draw(img, img.Bounds(), &image.Uniform{op.Color}, image.ZP, draw.Src)

	if op.LabelOptions != nil {
		err := addLabel(img, op.LabelOptions)
		if err != nil {
			log.Printf("Error adding label to image %v", err)
		}
	}

	return img, nil
}

func addLabel(img *image.RGBA, op *labelOptions) error {
	fontForeGroundColor := image.NewUniform(op.Color)
	ctx := freetype.NewContext()
	ctx.SetDPI(op.DPI)
	ctx.SetFont(op.Font)
	ctx.SetFontSize(op.Size)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetSrc(fontForeGroundColor)

	opts := truetype.Options{Size: op.Size}
	face := truetype.NewFace(op.Font, &opts)
	d := &font.Drawer{
		Face: face,
	}
	labelWidht := int(d.MeasureString(op.Text) >> 6)
	labelHeight := int(face.Metrics().Height >> 6)

	// Set the label centered alignment to x,y
	pt := freetype.Pt(op.X-labelWidht/2, op.Y+labelHeight/2)
	if _, err := ctx.DrawString(op.Text, pt); err != nil {
		return err
	}
	return nil
}
