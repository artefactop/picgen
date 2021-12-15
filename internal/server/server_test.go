package server

import (
	"net/http/httptest"
	"testing"
)

func TestWriteImage(t *testing.T) {
	ir := &imageRequest{
		Width:  defaultSize,
		Height: defaultSize,
		Color:  defaultColor,
		Format: defaultImageFormat,
	}

	rr := httptest.NewRecorder()
	_, err := writeImage(rr, ir)
	if err != nil {
		t.Error(err)
	}

	size := len(rr.Body.Bytes())
	expected := 297
	if size != expected {
		t.Errorf("handler returned unexpected body size: got %v want %v", size, expected)
	}
}

func TestWriteImageWithLabel(t *testing.T) {
	ir := &imageRequest{
		Width:      defaultSize,
		Height:     defaultSize,
		Color:      defaultColor,
		Format:     defaultImageFormat,
		LabelText:  "#",
		LabelColor: defaultLabelColor,
		LabelName:  "goregular",
		LabelDpi:   72.0,
		LabelSize:  calculateLabelSize(1, defaultSize, defaultSize),
		LabelX:     defaultSize / 2,
		LabelY:     defaultSize / 2,
	}

	rr := httptest.NewRecorder()
	_, err := writeImage(rr, ir)
	if err != nil {
		t.Error(err)
	}

	size := len(rr.Body.Bytes())
	// FIXME: This image is too big, maybe because the label
	expected := 1077
	if size != expected {
		t.Errorf("handler returned unexpected body size: got %v want %v", size, expected)
	}
}
