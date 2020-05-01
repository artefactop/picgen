package picgen

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	defaultSize      = int(100)
	defaultLabelSize = float64(65.0)
)

var (
	defaultColor     = color.RGBA{100, 200, 200, 255}
	errInvalidFormat = errors.New("invalid format")
)

// RootHandler ...
func RootHandler(w http.ResponseWriter, req *http.Request) {
	urlPart := strings.Split(req.URL.Path, "/")

	width, height := parseImageSize(urlPart[1])
	backgroundColor := parseColor(urlPart[2])

	fontColorAndImgTypePart := strings.Split(urlPart[3], ".")

	labelColor := parseColor(fontColorAndImgTypePart[0])

	queryValues := req.URL.Query()
	label := queryValues.Get("label")
	labelSize := parseLabelSize(queryValues.Get("size"))

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}

	log.Printf("Size:%dx%d, Color:%v, Label:%s, Label.Color:%v, Label.Size:%f Label.Font:%s",
		width, height, backgroundColor, label, labelColor, labelSize, font.Name(4))

	op := &options{
		Width:  width,
		Height: height,
		Color:  backgroundColor,
		LabelOptions: &labelOptions{
			Text:  label,
			Color: labelColor,
			Font:  font,
			DPI:   72.0,
			Size:  labelSize,
			X:     width / 2,
			Y:     height / 2,
		},
	}

	img, err := newImage(op)
	if err != nil {
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, img); err != nil {
		log.Fatal(err)
	}
}

func encodeImageFormat(w io.Writer, img image.Image, format string) error {
	return png.Encode(w, img)
}

func parseLabelSize(s string) float64 {
	sz, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultLabelSize
	}
	return sz
}

func parseImageSize(s string) (int, int) {
	sizePart := strings.Split(s, "x")
	w, err := strconv.Atoi(sizePart[0])
	if err != nil {
		return defaultSize, defaultSize
	}
	if len(sizePart) > 1 {
		h, err := strconv.Atoi(sizePart[1])
		if err != nil {
			return w, w
		}
		return w, h
	}
	return w, w
}

func parseColor(s string) color.RGBA {
	clr, err := parseNameColor(s)
	if err == nil {
		return clr
	}
	clr, err = parseHexColor(s)
	if err == nil {
		return clr
	}
	return defaultColor
}

func parseNameColor(s string) (color.RGBA, error) {
	if c, ok := colornames.Map[s]; ok == true {
		return c, nil
	}
	return color.RGBA{}, errInvalidFormat
}

func parseHexColor(s string) (color.RGBA, error) {
	c := color.RGBA{}
	c.A = 0xff

	var err error
	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 6:
		c.R = hexToByte(s[0])<<4 + hexToByte(s[1])
		c.G = hexToByte(s[2])<<4 + hexToByte(s[3])
		c.B = hexToByte(s[4])<<4 + hexToByte(s[5])
	case 4:
		c.R = hexToByte(s[0]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[4]) * 17
	default:
		return c, errInvalidFormat
	}
	return c, err
}
