package server

import (
	"image/color"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/image/colornames"
)

func parseRequest(r *http.Request) (*imageRequest, error) {
	vars := mux.Vars(r)

	sizeStr := vars["size"]
	width, height := parseImageSize(sizeStr)

	if width*height >= maxArea || width > maxWidth || height > maxHeight {
		return &imageRequest{}, errInvalidFormat
	}
	colorStr := vars["color"]
	backgroundColor := parseColor(colorStr)

	labelColorStr := vars["labelColor"]
	fontColorAndImgFormatPart := strings.Split(labelColorStr, ".")

	labelColor := parseColor(fontColorAndImgFormatPart[0])

	imageFormat := defaultImageFormat
	if len(fontColorAndImgFormatPart) > 1 {
		imageFormat = strings.ToUpper(fontColorAndImgFormatPart[1])
	}

	queryValues := r.URL.Query()
	label := queryValues.Get("text")
	labelSize := parseLabelSize(queryValues.Get("size"))

	fontName := "goregular"

	log.Printf("Size:%dx%d, Format:%s, Color:%v, Label:%s, Label.Color:%v, Label.Size:%f Label.Font:%s",
		width, height, imageFormat, backgroundColor, label, labelColor, labelSize, fontName)

	op := &imageRequest{
		Width:      width,
		Height:     height,
		Color:      backgroundColor,
		Format:     imageFormat,
		LabelText:  label,
		LabelColor: labelColor,
		LabelName:  fontName,
		LabelDpi:   72.0,
		LabelSize:  labelSize,
		LabelX:     width / 2,
		LabelY:     height / 2,
	}
	return op, nil
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
	case 3:
		c.R = hexToByte(s[0]) * 17
		c.G = hexToByte(s[1]) * 17
		c.B = hexToByte(s[2]) * 17
	default:
		return c, errInvalidFormat
	}
	return c, err
}
