package server

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"net/http"
	"strconv"
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
	defaultColor     = color.RGBA{100, 200, 200, 255}
	errInvalidFormat = errors.New("invalid format")
)

func addCacheHeaders(w http.ResponseWriter) {
	cacheSince := time.Now().Format(http.TimeFormat)
	cacheUntil := time.Now().AddDate(0, 0, 180).Format(http.TimeFormat) // 180 days
	w.Header().Set("Cache-Control", "public, max-age=15552000")         // 180 days
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

// PathHandler ...
func PathHandler(w http.ResponseWriter, req *http.Request) {
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

	img := image.NewImage(op.Width, op.Height, op.Color)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}

	f, err := truetype.Parse(goregular.TTF)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "kaboom %v", err)
	}

	label := image.NewLabel(op.LabelText, op.LabelColor, op.LabelDpi, f, op.LabelSize)

	err = image.DrawLabel(img, *label, op.LabelX, op.LabelY)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}

	addCacheHeaders(w)

	buffer := new(bytes.Buffer)
	size, err := image.Encode(buffer, img, op.Format)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(size))
	switch op.Format {
	case "JPEG", "JPG":
		w.Header().Set("Content-Type", "image/jpeg")
	case "PNG":
		fallthrough
	default:
		w.Header().Set("Content-Type", "image/png")
	}

	_, err = w.Write(buffer.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "kaboom %v", err)
		return
	}

}
