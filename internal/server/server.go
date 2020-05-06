package server

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"time"

	"github.com/artefactop/picgen/internal/image"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	defaultSize        = int(100)
	defaultLabelSize   = float64(65.0)
	defaultImageFormat = string("PNG")
	maxArea            = int(16000000)
	maxWidth           = int(9999)
	maxHeight          = int(9999)
)

var (
	defaultColor      = color.RGBA{1, 173, 216, 255}
	defaultLabelColor = color.RGBA{255, 255, 255, 255}
	errInvalidFormat  = errors.New("invalid format")
)

func addCacheHeaders(w http.ResponseWriter) {
	cacheSince := time.Now().Format(http.TimeFormat)
	cacheUntil := time.Now().AddDate(0, 0, 180).Format(http.TimeFormat)    // 180 days
	w.Header().Set("Cache-Control", "public, immutable, max-age=15552000") // 180 days
	w.Header().Set("Last-Modified", cacheSince)
	w.Header().Set("Expires", cacheUntil)
}

type imageRequest struct {
	Width  int
	Height int
	Color  color.Color
	Format string

	LabelText  string
	LabelColor color.Color
	LabelName  string
	LabelDpi   float64
	LabelSize  float64
	LabelX     int
	LabelY     int
}

// RootHandler ...
func RootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == http.MethodOptions {
		return
	}

	queryValues := req.URL.Query()

	width, height := parseImageSize(queryValues.Get("x"))

	if width*height >= maxArea || width > maxWidth || height > maxHeight {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "kaboom %v", errInvalidFormat)
		return
	}

	backgroundColor := parseColor(queryValues.Get("b"), defaultColor)
	labelText := queryValues.Get("t")
	if labelText == "" {
		labelText = "#"
	}
	labelColor := parseColor(queryValues.Get("f"), defaultLabelColor)

	fontName := "goregular"

	log.Printf("Size:%dx%d, Format:%s, Color:%v, Label.Text:%s, Label.Color:%v, Label.Size:%f Label.Font:%s",
		width, height, defaultImageFormat, backgroundColor, labelText, labelColor, defaultLabelSize, fontName)

	ir := &imageRequest{
		Width:      width,
		Height:     height,
		Color:      backgroundColor,
		Format:     defaultImageFormat,
		LabelText:  labelText,
		LabelColor: labelColor,
		LabelName:  fontName,
		LabelDpi:   72.0,
		LabelSize:  defaultLabelSize,
		LabelX:     width / 2,
		LabelY:     height / 2,
	}

	writeImage(w, ir)
}

// PathHandler ...
func PathHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == http.MethodOptions {
		return
	}

	ir, err := parseRequest(req)
	if err != nil {
		if err == errInvalidFormat {
			w.WriteHeader(http.StatusBadRequest)
		}
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}

	_, err = writeImage(w, ir)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}
}

func writeImage(w http.ResponseWriter, ir *imageRequest) (int, error) {
	img := image.NewImage(ir.Width, ir.Height, ir.Color)

	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return 0, err
	}

	if ir.LabelText != "" {
		label := image.NewLabel(ir.LabelText, ir.LabelColor, ir.LabelDpi, f, ir.LabelSize)

		err = image.DrawLabel(img, *label, ir.LabelX, ir.LabelY)
		if err != nil {
			return 0, err
		}
	}

	addCacheHeaders(w)

	return image.Encode(w, img, ir.Format)
}
