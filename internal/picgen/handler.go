package picgen

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/gorilla/mux"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	defaultSize        = int(100)
	defaultLabelSize   = float64(65.0)
	defaultImageFormat = string("png")
	maxArea            = int(16000000)
	maxWidth           = int(9999)
	maxHeight          = int(9999)
)

var (
	defaultColor     = color.RGBA{100, 200, 200, 255}
	errInvalidFormat = errors.New("invalid format")
)

// RootHandler ...
func RootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == http.MethodOptions {
		return
	}

	op, err := parseRequest(req)
	if err != nil {
		if err == errInvalidFormat {
			w.WriteHeader(http.StatusBadRequest)
		}
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}

	img, err := newImage(op)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}

	// TODO move to a middleware
	cacheSince := time.Now().Format(http.TimeFormat)
	cacheUntil := time.Now().AddDate(0, 0, 180).Format(http.TimeFormat) // 180 days
	w.Header().Set("Cache-Control", "public, max-age=15552000")         // 180 days
	w.Header().Set("Last-Modified", cacheSince)
	w.Header().Set("Expires", cacheUntil)

	switch op.Format {
	case "jpeg", "jpg":
		w.Header().Set("Content-Type", "image/jpeg")
		if err := jpeg.Encode(w, img, &jpeg.Options{}); err != nil {
			log.Fatal(err)
		}
	case "png":
		fallthrough
	default:
		w.Header().Set("Content-Type", "image/png")
		if err := png.Encode(w, img); err != nil {
			log.Fatal(err)
		}
	}

}

func parseRequest(r *http.Request) (*options, error) {
	vars := mux.Vars(r)

	sizeStr := vars["size"]
	width, height := parseImageSize(sizeStr)

	if width*height >= maxArea || width > maxWidth || height > maxHeight {
		return &options{}, errInvalidFormat
	}
	colorStr := vars["color"]
	backgroundColor := parseColor(colorStr)

	labelColorStr := vars["labelColor"]
	fontColorAndImgFormatPart := strings.Split(labelColorStr, ".")

	labelColor := parseColor(fontColorAndImgFormatPart[0])

	imageFormat := defaultImageFormat
	if len(fontColorAndImgFormatPart) > 1 {
		imageFormat = fontColorAndImgFormatPart[1]
	}

	queryValues := r.URL.Query()
	label := queryValues.Get("text")
	labelSize := parseLabelSize(queryValues.Get("size"))

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return &options{}, err
	}

	log.Printf("Size:%dx%d, Format:%s, Color:%v, Label:%s, Label.Color:%v, Label.Size:%f Label.Font:%s",
		width, height, imageFormat, backgroundColor, label, labelColor, labelSize, font.Name(4))

	op := &options{
		Width:  width,
		Height: height,
		Color:  backgroundColor,
		Format: imageFormat,
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
	return op, nil
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
	case 3:
		c.R = hexToByte(s[0]) * 17
		c.G = hexToByte(s[1]) * 17
		c.B = hexToByte(s[2]) * 17
	default:
		return c, errInvalidFormat
	}
	return c, err
}
